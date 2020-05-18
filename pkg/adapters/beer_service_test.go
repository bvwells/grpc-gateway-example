package adapters_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/adapters"
	"github.com/bvwells/grpc-gateway-example/pkg/adapters/mocks"
	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/proto/beers"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery -name=BeerInteractor -case=underscore

func TestNewBeerService_ReturnsBeerInteractor(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	assert.NotNil(t, adapters.NewBeerService(interactor))
}

func TestCreateBeer_WhenCreateBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.CreateBeerRequest{
		Name:    "a beer",
		Type:    beers.Type_Ale,
		Brewer:  "brewer",
		Country: "country",
	}
	expected := errors.New("something went wrong")
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    params.Name,
		Type:    domain.Ale,
		Brewer:  params.Brewer,
		Country: params.Country,
	}).Return(nil, expected)
	_, actual := service.CreateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestCreateBeer_WhenCreateBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	expected := &beers.CreateBeerResponse{
		Beer: &beers.Beer{
			Id:      "id",
			Name:    "a beer",
			Type:    beers.Type_Ale,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	params := &beers.CreateBeerRequest{
		Name:    "a beer",
		Type:    beers.Type_Ale,
		Brewer:  "brewer",
		Country: "country",
	}
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    expected.Beer.Name,
		Type:    domain.Ale,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.Ale,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}, nil)
	actual, _ := service.CreateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeer_WhenGetBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.GetBeerRequest{Id: "ID"}
	expected := errors.New("something went wrong")
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(nil, expected)
	_, actual := service.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeer_WhenGetBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.GetBeerRequest{Id: "ID"}
	expected := &beers.GetBeerResponse{
		Beer: &beers.Beer{
			Id:      "id",
			Name:    "a beer",
			Type:    beers.Type_Ale,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.Ale,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}, nil)
	actual, _ := service.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenUpdateBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{Id: "id"}
	expected := errors.New("something went wrong")
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{ID: params.Id}).Return(nil, expected)
	_, actual := service.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenUpdateBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{Id: "id"}
	expected := &beers.UpdateBeerResponse{
		Beer: &beers.Beer{
			Id:      "id",
			Name:    "a beer",
			Type:    beers.Type_Ale,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{ID: params.Id}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.Ale,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}, nil)
	actual, _ := service.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestDeleteBeer_WhenDeleteBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.DeleteBeerRequest{Id: "id"}
	expected := errors.New("something went wrong")
	interactor.On("DeleteBeer", ctx, &domain.DeleteBeerParams{ID: params.Id}).Return(expected)
	_, actual := service.DeleteBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestDeleteBeer_WhenDeleteBeerReturnsNil_ReturnsNilError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.DeleteBeerRequest{Id: "id"}
	interactor.On("DeleteBeer", ctx, &domain.DeleteBeerParams{ID: params.Id}).Return(nil)
	_, actual := service.DeleteBeer(ctx, params)
	assert.Nil(t, actual)
}

func TestGetBeers_WhenGetBeersReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	expected := errors.New("something went wrong")
	interactor.On("GetBeers", ctx, &domain.GetBeersParams{}).Return(nil, expected)
	_, actual := service.GetBeers(ctx, &beers.GetBeersRequest{})
	assert.Equal(t, expected, actual)
}

func TestGetBeers_WhenGetBeersReturnsBeers_ReturnsBeers(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	expected := &beers.GetBeersResponse{
		Beers: []*beers.Beer{
			{Id: "id1"},
			{Id: "id2"},
		},
	}
	interactor.On("GetBeers", ctx, &domain.GetBeersParams{}).Return([]*domain.Beer{
		{ID: expected.Beers[0].Id},
		{ID: expected.Beers[1].Id},
	}, nil)
	actual, _ := service.GetBeers(ctx, &beers.GetBeersRequest{})
	assert.Equal(t, expected, actual)
}
