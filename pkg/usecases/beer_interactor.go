package usecases

import (
	"context"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
)

// NewBeerInteractor creates a new beer interactor.
func NewBeerInteractor(repo BeerRepository) *BeerInteractor {
	return &BeerInteractor{repo: repo}
}

// BeerInteractor describes a set of APIs for interacting with beers.
type BeerInteractor struct {
	repo BeerRepository
}

// CreateBeer is an API for getting a beer given its ID.
func (interactor *BeerInteractor) CreateBeer(ctx context.Context, params *domain.CreateBeerParams) (*domain.Beer, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	beer, err := interactor.repo.CreateBeer(ctx, params)
	if err != nil {
		return nil, err
	}
	return beer, nil
}

// GetBeer is an API for getting a beer given its ID.
func (interactor *BeerInteractor) GetBeer(ctx context.Context, params *domain.GetBeerParams) (*domain.Beer, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	beer, err := interactor.repo.GetBeer(ctx, params)
	if err != nil {
		return nil, err
	}
	return beer, nil
}

// UpdateBeer is an API for updating a beer given its ID.
func (interactor *BeerInteractor) UpdateBeer(ctx context.Context, params *domain.UpdateBeerParams) (*domain.Beer, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	beer, err := interactor.repo.UpdateBeer(ctx, params)
	if err != nil {
		return nil, err
	}
	return beer, nil
}

// DeleteBeer is an API for deleting a beer given its ID.
func (interactor *BeerInteractor) DeleteBeer(ctx context.Context, params *domain.DeleteBeerParams) error {
	err := params.Validate()
	if err != nil {
		return err
	}

	err = interactor.repo.DeleteBeer(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

// GetBeers is an API for getting beers.
func (interactor *BeerInteractor) GetBeers(ctx context.Context, params *domain.GetBeersParams) ([]*domain.Beer, error) {
	err := params.Validate()
	if err != nil {
		return nil, err
	}

	beers, err := interactor.repo.GetBeers(ctx, params)
	if err != nil {
		return nil, err
	}
	return beers, nil
}
