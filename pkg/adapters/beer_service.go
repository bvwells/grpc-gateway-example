package adapters

import (
	"context"
	"errors"
	"strings"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/proto/beers"

	"github.com/golang/protobuf/ptypes/empty"
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
	// ListBeers lists beers.
	ListBeers(ctx context.Context, params *domain.ListBeersParams) ([]*domain.Beer, error)
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
func (svc *BeerService) CreateBeer(ctx context.Context, params *beers.CreateBeerRequest) (*beers.Beer, error) {
	item, err := svc.interactor.CreateBeer(ctx, &domain.CreateBeerParams{
		Name:    params.Name,
		Type:    fromProtoType(params.Type),
		Brewer:  params.Brewer,
		Country: params.Country,
	})
	if err != nil {
		return nil, toError(err)
	}
	return toProtoBeer(item), nil
}

// GetBeer gets the beer with specified beer identifier.
func (svc *BeerService) GetBeer(ctx context.Context, params *beers.GetBeerRequest) (*beers.Beer, error) {
	item, err := svc.interactor.GetBeer(ctx, &domain.GetBeerParams{ID: params.Id})
	if err != nil {
		return nil, toError(err)
	}
	return toProtoBeer(item), nil
}

// UpdateBeer updates the beer with specified beer identifier.
func (svc *BeerService) UpdateBeer(ctx context.Context, params *beers.UpdateBeerRequest) (*beers.Beer, error) {
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
		case "type":
			beerType := fromProtoType(params.Beer.Type)
			updateParams.Type = &beerType
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
		return nil, toError(err)
	}
	return toProtoBeer(item), nil
}

// DeleteBeer deletes the beer with specified beer identifier.
func (svc *BeerService) DeleteBeer(ctx context.Context, params *beers.DeleteBeerRequest) (*empty.Empty, error) {
	err := svc.interactor.DeleteBeer(ctx, &domain.DeleteBeerParams{ID: params.Id})
	if err != nil {
		return nil, toError(err)
	}
	return &empty.Empty{}, nil
}

// ListBeers lists all beers.
func (svc *BeerService) ListBeers(ctx context.Context, params *beers.ListBeersRequest) (*beers.ListBeersResponse, error) {
	items, err := svc.interactor.ListBeers(ctx, &domain.ListBeersParams{
		Page: int(params.Page),
	})
	if err != nil {
		return nil, toError(err)
	}
	b := &beers.ListBeersResponse{
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

func toProtoType(in domain.BeerType) beers.BeerType {
	switch in {
	case domain.Ale:
		return beers.BeerType_BEER_TYPE_ALE
	case domain.Bitter:
		return beers.BeerType_BEER_TYPE_BITTER
	case domain.Lager:
		return beers.BeerType_BEER_TYPE_LAGER
	case domain.IndiaPaleAle:
		return beers.BeerType_BEER_TYPE_INDIA_PALE_ALE
	case domain.Stout:
		return beers.BeerType_BEER_TYPE_STOUT
	case domain.Pilsner:
		return beers.BeerType_BEER_TYPE_PILSNER
	case domain.Porter:
		return beers.BeerType_BEER_TYPE_PORTER
	case domain.PaleAle:
		return beers.BeerType_BEER_TYPE_PALE_ALE
	case domain.Unspecified:
		return beers.BeerType_BEER_TYPE_UNSPECIFIED
	}
	return beers.BeerType_BEER_TYPE_UNSPECIFIED
}

func fromProtoType(in beers.BeerType) domain.BeerType {
	switch in {
	case beers.BeerType_BEER_TYPE_ALE:
		return domain.Ale
	case beers.BeerType_BEER_TYPE_BITTER:
		return domain.Bitter
	case beers.BeerType_BEER_TYPE_LAGER:
		return domain.Lager
	case beers.BeerType_BEER_TYPE_INDIA_PALE_ALE:
		return domain.IndiaPaleAle
	case beers.BeerType_BEER_TYPE_STOUT:
		return domain.Stout
	case beers.BeerType_BEER_TYPE_PILSNER:
		return domain.Pilsner
	case beers.BeerType_BEER_TYPE_PORTER:
		return domain.Porter
	case beers.BeerType_BEER_TYPE_PALE_ALE:
		return domain.PaleAle
	case beers.BeerType_BEER_TYPE_UNSPECIFIED:
		return domain.Unspecified
	}
	return domain.Unspecified
}

func toError(err error) error {
	if errors.As(err, &domain.ValidationError{}) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return status.Error(codes.Internal, err.Error())
}
