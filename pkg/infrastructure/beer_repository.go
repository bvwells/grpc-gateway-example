package infrastructure

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"

	_ "github.com/lib/pq" // imported for side effect.
)

const numberRows = 50

// PostgresSettings describes all the settings required for setting up
// a connection to a postgres database.
type PostgresSettings struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// String returns the string representation for a postgres database.
func (s *PostgresSettings) String() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		s.Host, s.Port, s.User, s.Password, s.DBName)
}

// postgresBeer is the postgres representation of a beer.
type postgresBeer struct {
	ID      string
	Name    string
	Type    int
	Brewer  string
	Country string
}

// GenerateID generates a unique identifier.
type GenerateID func() string

// NewPostgresBeerRepository creates a new postgres beer repository.
func NewPostgresBeerRepository(settings *PostgresSettings,
	generateID GenerateID) (*PostgresBeerRepository, error) {
	db, err := sql.Open("postgres", settings.String())
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return &PostgresBeerRepository{
		db:         db,
		generateID: generateID,
	}, nil
}

// PostgresBeerRepository is a postgres beer repository.
type PostgresBeerRepository struct {
	db         *sql.DB
	generateID func() string
}

// Close closes the postgres database.
func (repo *PostgresBeerRepository) Close() error {
	if repo.db != nil {
		return repo.db.Close()
	}
	return nil
}

// CreateBeer creates a beer in the postgres database.
func (repo *PostgresBeerRepository) CreateBeer(ctx context.Context, params *domain.CreateBeerParams) (*domain.Beer, error) {
	id := repo.generateID()
	sqlStatement := `
	INSERT INTO BEERS (id, name, type, brewer, country)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id`
	err := repo.db.QueryRow(sqlStatement, id, params.Name, params.Type, params.Brewer, params.Country).Scan(&id)
	if err != nil {
		return nil, err
	}
	return repo.GetBeer(ctx, &domain.GetBeerParams{ID: id})
}

// GetBeer gets a beer from the postgres database.
func (repo *PostgresBeerRepository) GetBeer(ctx context.Context, params *domain.GetBeerParams) (*domain.Beer, error) {
	var beer postgresBeer
	row := repo.db.QueryRow("SELECT * FROM BEERS WHERE id=$1;", params.ID)
	err := row.Scan(&beer.ID, &beer.Name, &beer.Type, &beer.Brewer, &beer.Country)
	switch err {
	case sql.ErrNoRows:
		return nil, errors.New("not found")
	case nil:
		return toDomainBeer(&beer), nil
	default:
		return nil, errors.New("something unexpected happened")
	}
}

// UpdateBeer updates a beer in the postgres database.
func (repo *PostgresBeerRepository) UpdateBeer(ctx context.Context, params *domain.UpdateBeerParams) (*domain.Beer, error) {
	if params.Name != nil {
		sqlStatement := `UPDATE BEERS
						 SET name = $2
						 WHERE id = $1;`
		_, err := repo.db.Exec(sqlStatement, params.ID, params.Name)
		if err != nil {
			return nil, err
		}
	}
	if params.Brewer != nil {
		sqlStatement := `UPDATE BEERS
						 SET brewer = $2
						 WHERE id = $1;`
		_, err := repo.db.Exec(sqlStatement, params.ID, params.Brewer)
		if err != nil {
			return nil, err
		}
	}
	if params.Country != nil {
		sqlStatement := `UPDATE BEERS
						 SET country = $2
						 WHERE id = $1;`
		_, err := repo.db.Exec(sqlStatement, params.ID, params.Country)
		if err != nil {
			return nil, err
		}
	}
	return repo.GetBeer(ctx, &domain.GetBeerParams{ID: params.ID})
}

// DeleteBeer deletes a beer from the postgres database.
func (repo *PostgresBeerRepository) DeleteBeer(ctx context.Context, params *domain.DeleteBeerParams) error {
	sqlStatement := `
	DELETE FROM BEERS
	WHERE id = $1;`
	_, err := repo.db.Exec(sqlStatement, params.ID)
	if err != nil {
		return err
	}
	return nil
}

// GetBeers gets all beers from the postgres database.
// TODO - return cursor to new batch of beers.
func (repo *PostgresBeerRepository) GetBeers(ctx context.Context, params *domain.GetBeersParams) ([]*domain.Beer, error) {
	rows, err := repo.db.Query("SELECT * FROM BEERS LIMIT $1", numberRows)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var beers []*domain.Beer
	for rows.Next() {
		var beer postgresBeer
		err := rows.Scan(&beer.ID, &beer.Name, &beer.Type, &beer.Brewer, &beer.Country)

		if err != nil {
			return nil, err
		}
		beers = append(beers, toDomainBeer(&beer))
	}

	// Check for any errors encountered.
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return beers, nil
}

func toDomainBeer(in *postgresBeer) *domain.Beer {
	return &domain.Beer{
		ID:      in.ID,
		Name:    in.Name,
		Type:    domain.Type(in.Type),
		Brewer:  in.Brewer,
		Country: in.Country,
	}
}
