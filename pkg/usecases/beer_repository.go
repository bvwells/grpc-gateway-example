package usecases

import (
	"context"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
)

// BeerRepository is a repository for beers.
type BeerRepository interface {
	// CreateBeer creates a beer.
	CreateBeer(ctx context.Context, params *domain.CreateBeerParams) (*domain.Beer, error)
	// GetBeer gets a beer.
	GetBeer(ctx context.Context, params *domain.GetBeerParams) (*domain.Beer, error)
	// UpdateBeer updates a beer.
	UpdateBeer(ctx context.Context, params *domain.UpdateBeerParams) (*domain.Beer, error)
	// DeleteBeer deletes a beer.
	DeleteBeer(ctx context.Context, params *domain.DeleteBeerParams) error
	// ListBeers lists beers.
	ListBeers(ctx context.Context, params *domain.ListBeersParams) ([]*domain.Beer, error)
}
