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
	"google.golang.org/genproto/protobuf/field_mask"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		Type:    beers.Type_ALE,
		Brewer:  "brewer",
		Country: "country",
	}
	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    params.Name,
		Type:    domain.Ale,
		Brewer:  params.Brewer,
		Country: params.Country,
	}).Return(nil, errors.New(msg))
	_, actual := service.CreateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestCreateBeer_WhenCreateBeerReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.CreateBeerRequest{
		Name:    "a beer",
		Type:    beers.Type_BITTER,
		Brewer:  "brewer",
		Country: "country",
	}
	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    params.Name,
		Type:    domain.Bitter,
		Brewer:  params.Brewer,
		Country: params.Country,
	}).Return(nil, domain.NewValidationError(msg))
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
			Type:    beers.Type_LAGER,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	params := &beers.CreateBeerRequest{
		Name:    "a beer",
		Type:    beers.Type_LAGER,
		Brewer:  "brewer",
		Country: "country",
	}
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    expected.Beer.Name,
		Type:    domain.Lager,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.Lager,
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

	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(nil, errors.New(msg))
	_, actual := service.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeer_WhenGetBeerReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.GetBeerRequest{Id: "ID"}

	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(nil, domain.NewValidationError(msg))
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
			Type:    beers.Type_INDIA_PALE_ALE,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.IndiaPaleAle,
		Brewer:  expected.Beer.Brewer,
		Country: expected.Beer.Country,
	}, nil)
	actual, _ := service.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenFieldMaskNotSpecified_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{
		Beer: &beers.Beer{Id: "id", Name: "name"},
	}
	_, err := service.UpdateBeer(ctx, params)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateBeer_WhenFieldMaskContainsInvalidField_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{
		Beer:       &beers.Beer{Id: "id", Name: "name"},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"invalid_field"}},
	}
	_, err := service.UpdateBeer(ctx, params)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
}

func TestUpdateBeer_WhenUpdateBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{
		Beer:       &beers.Beer{Id: "id", Name: "name"},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"NaMe"}},
	}
	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{ID: params.Beer.Id, Name: &params.Beer.Name}).Return(nil, errors.New(msg))
	_, actual := service.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenUpdateBeerReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{
		Beer:       &beers.Beer{Id: "id", Name: "name"},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"NaMe"}},
	}
	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{ID: params.Beer.Id, Name: &params.Beer.Name}).Return(nil, domain.NewValidationError(msg))
	_, actual := service.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenUpdateBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.UpdateBeerRequest{
		Beer:       &beers.Beer{Id: "id", Name: "name", Type: beers.Type_STOUT, Brewer: "brewer", Country: "Country"},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"NaMe", "type", "Country", "breweR"}},
	}
	expected := &beers.UpdateBeerResponse{
		Beer: &beers.Beer{
			Id:      "id",
			Name:    "a beer",
			Type:    beers.Type_STOUT,
			Brewer:  "brewer",
			Country: "country",
		},
	}
	beerType := domain.Stout
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{
		ID:      params.Beer.Id,
		Name:    &params.Beer.Name,
		Brewer:  &params.Beer.Brewer,
		Type:    &beerType,
		Country: &params.Beer.Country,
	}).Return(&domain.Beer{
		ID:      expected.Beer.Id,
		Name:    expected.Beer.Name,
		Type:    domain.Stout,
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
	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("DeleteBeer", ctx, &domain.DeleteBeerParams{ID: params.Id}).Return(errors.New(msg))
	_, actual := service.DeleteBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestDeleteBeer_WhenDeleteBeerReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	params := &beers.DeleteBeerRequest{Id: "id"}
	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("DeleteBeer", ctx, &domain.DeleteBeerParams{ID: params.Id}).Return(domain.NewValidationError(msg))
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
	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("GetBeers", ctx, &domain.GetBeersParams{Page: 42}).Return(nil, errors.New(msg))
	_, actual := service.GetBeers(ctx, &beers.GetBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}

func TestGetBeers_WhenGetBeersReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("GetBeers", ctx, &domain.GetBeersParams{Page: 42}).Return(nil, domain.NewValidationError(msg))
	_, actual := service.GetBeers(ctx, &beers.GetBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}

func TestGetBeers_WhenGetBeersReturnsBeers_ReturnsBeers(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	expected := &beers.GetBeersResponse{
		Beers: []*beers.Beer{
			{Id: "id1", Type: beers.Type_PILSNER},
			{Id: "id2", Type: beers.Type_PORTER},
			{Id: "id3", Type: beers.Type_PALE_ALE},
			{Id: "id4", Type: beers.Type_UNKNOWN},
		},
	}
	interactor.On("GetBeers", ctx, &domain.GetBeersParams{Page: 42}).Return([]*domain.Beer{
		{ID: expected.Beers[0].Id, Type: domain.Pilsner},
		{ID: expected.Beers[1].Id, Type: domain.Porter},
		{ID: expected.Beers[2].Id, Type: domain.PaleAle},
		{ID: expected.Beers[3].Id, Type: domain.Unknown},
	}, nil)
	actual, _ := service.GetBeers(ctx, &beers.GetBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}
