package adapters

import (
	"context"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/proto/beers"
)

// BeerInteractor defines a set of APIs for interacting with beers.
type BeerInteractor interface {
	// CreateBeer creates a beers.
	CreateBeer(ctx context.Context, params *domain.CreateBeerParams) (*domain.Beer, error)
	// GetBeer gets a beers.
	GetBeer(ctx context.Context, params *domain.GetBeerParams) (*domain.Beer, error)
	// UpdateBeer uodates a beers.
	UpdateBeer(ctx context.Context, params *domain.UpdateBeerParams) (*domain.Beer, error)
	// DeleteBeer deletes a beers.
	DeleteBeer(ctx context.Context, params *domain.DeleteBeerParams) error
	// GetBeers gets beers.
	GetBeers(ctx context.Context, params *domain.GetBeersParams) ([]*domain.Beer, error)
}

// NewBeerService creates a new beer service.
func NewBeerService(interactor BeerInteractor) *BeerService {
	return &BeerService{interactor: interactor}
}

// BeerService implements the BeerService service gRPC API.
type BeerService struct {
	interactor BeerInteractor
}

// CreateBeer create a beer with specified beer parameters.
func (svc *BeerService) CreateBeer(ctx context.Context, params *beers.CreateBeerRequest) (*beers.CreateBeerResponse, error) {
	item, err := svc.interactor.CreateBeer(ctx, &domain.CreateBeerParams{
		Name:    params.Name,
		Type:    fromProtoType(params.Type),
		Brewer:  params.Brewer,
		Country: params.Country,
	})
	if err != nil {
		return nil, err
	}
	return &beers.CreateBeerResponse{
		Beer: toProtoBeer(item),
	}, nil
}

// GetBeer gets the beer with specified beer identifier.
func (svc *BeerService) GetBeer(ctx context.Context, params *beers.GetBeerRequest) (*beers.GetBeerResponse, error) {
	item, err := svc.interactor.GetBeer(ctx, &domain.GetBeerParams{ID: params.Id})
	if err != nil {
		return nil, err
	}
	return &beers.GetBeerResponse{
		Beer: toProtoBeer(item),
	}, nil
}

// UpdateBeer updates the beer with specified beer identifier.
func (svc *BeerService) UpdateBeer(ctx context.Context, params *beers.UpdateBeerRequest) (*beers.UpdateBeerResponse, error) {
	item, err := svc.interactor.UpdateBeer(ctx, &domain.UpdateBeerParams{ID: params.Id})
	if err != nil {
		return nil, err
	}
	return &beers.UpdateBeerResponse{
		Beer: toProtoBeer(item),
	}, nil
}

// DeleteBeer deletes the beer with specified beer identifier.
func (svc *BeerService) DeleteBeer(ctx context.Context, params *beers.DeleteBeerRequest) (*beers.DeleteBeerResponse, error) {
	err := svc.interactor.DeleteBeer(ctx, &domain.DeleteBeerParams{ID: params.Id})
	if err != nil {
		return nil, err
	}
	return &beers.DeleteBeerResponse{}, nil
}

// GetBeers gets all beers.
func (svc *BeerService) GetBeers(ctx context.Context, _ *beers.GetBeersRequest) (*beers.GetBeersResponse, error) {
	items, err := svc.interactor.GetBeers(ctx, &domain.GetBeersParams{})
	if err != nil {
		return nil, err
	}
	b := &beers.GetBeersResponse{
		Beers: make([]*beers.Beer, 0, len(items)),
	}
	for _, item := range items {
		b.Beers = append(b.Beers, toProtoBeer(item))
	}
	return b, nil
}

func toProtoBeer(in *domain.Beer) *beers.Beer {
	return &beers.Beer{
		Id:      in.ID,
		Name:    in.Name,
		Type:    toProtoType(in.Type),
		Brewer:  in.Brewer,
		Country: in.Country,
	}
}

func toProtoType(in domain.Type) beers.Type {
	switch in {
	case domain.Ale:
		return beers.Type_Ale
	case domain.Bitter:
		return beers.Type_Bitter
	case domain.Larger:
		return beers.Type_Larger
	case domain.IndiaPaleAle:
		return beers.Type_IndiaPaleAle
	case domain.Stout:
		return beers.Type_Stout
	case domain.Pilsner:
		return beers.Type_Pilsner
	case domain.Porter:
		return beers.Type_Porter
	case domain.PaleAle:
		return beers.Type_PaleAle
	case domain.Unknown:
		return beers.Type_Unknown
	}
	return beers.Type_Unknown
}

func fromProtoType(in beers.Type) domain.Type {
	switch in {
	case beers.Type_Ale:
		return domain.Ale
	case beers.Type_Bitter:
		return domain.Bitter
	case beers.Type_Larger:
		return domain.Larger
	case beers.Type_IndiaPaleAle:
		return domain.IndiaPaleAle
	case beers.Type_Stout:
		return domain.Stout
	case beers.Type_Pilsner:
		return domain.Pilsner
	case beers.Type_Porter:
		return domain.Porter
	case beers.Type_PaleAle:
		return domain.PaleAle
	case beers.Type_Unknown:
		return domain.Unknown
	}
	return domain.Unknown
}
