package adapters

import (
	"context"
	"strings"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/proto/beers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
	if params.UpdateMask == nil {
		return nil, status.Error(codes.InvalidArgument, "no fields specified")
	}

	updateParams := &domain.UpdateBeerParams{
		ID: params.Beer.Id,
	}
	for _, path := range params.UpdateMask.Paths {
		switch field := strings.ToLower(path); field {
		case "name":
			updateParams.Name = &params.Beer.Name
		case "brewer":
			updateParams.Brewer = &params.Beer.Brewer
		case "country":
			updateParams.Country = &params.Beer.Country
		default:
			return nil, status.Errorf(codes.InvalidArgument, "invalid beer field: %s", field)
		}
	}

	item, err := svc.interactor.UpdateBeer(ctx, updateParams)
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
		return beers.Type_ALE
	case domain.Bitter:
		return beers.Type_BITTER
	case domain.Larger:
		return beers.Type_LARGER
	case domain.IndiaPaleAle:
		return beers.Type_INDIA_PALE_ALE
	case domain.Stout:
		return beers.Type_STOUT
	case domain.Pilsner:
		return beers.Type_PILSNER
	case domain.Porter:
		return beers.Type_PORTER
	case domain.PaleAle:
		return beers.Type_PALE_ALE
	case domain.Unknown:
		return beers.Type_UNKNOWN
	}
	return beers.Type_UNKNOWN
}

func fromProtoType(in beers.Type) domain.Type {
	switch in {
	case beers.Type_ALE:
		return domain.Ale
	case beers.Type_BITTER:
		return domain.Bitter
	case beers.Type_LARGER:
		return domain.Larger
	case beers.Type_INDIA_PALE_ALE:
		return domain.IndiaPaleAle
	case beers.Type_STOUT:
		return domain.Stout
	case beers.Type_PILSNER:
		return domain.Pilsner
	case beers.Type_PORTER:
		return domain.Porter
	case beers.Type_PALE_ALE:
		return domain.PaleAle
	case beers.Type_UNKNOWN:
		return domain.Unknown
	}
	return domain.Unknown
}
