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
		Type:    beers.BeerType_BEER_TYPE_ALE,
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
		Type:    beers.BeerType_BEER_TYPE_BITTER,
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
	expected := &beers.Beer{
		Id:      "id",
		Name:    "a beer",
		Type:    beers.BeerType_BEER_TYPE_LAGER,
		Brewer:  "brewer",
		Country: "country",
	}
	params := &beers.CreateBeerRequest{
		Name:    "a beer",
		Type:    beers.BeerType_BEER_TYPE_LAGER,
		Brewer:  "brewer",
		Country: "country",
	}
	interactor.On("CreateBeer", ctx, &domain.CreateBeerParams{
		Name:    expected.Name,
		Type:    domain.Lager,
		Brewer:  expected.Brewer,
		Country: expected.Country,
	}).Return(&domain.Beer{
		ID:      expected.Id,
		Name:    expected.Name,
		Type:    domain.Lager,
		Brewer:  expected.Brewer,
		Country: expected.Country,
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
	expected := &beers.Beer{
		Id:      "id",
		Name:    "a beer",
		Type:    beers.BeerType_BEER_TYPE_INDIA_PALE_ALE,
		Brewer:  "brewer",
		Country: "country",
	}
	interactor.On("GetBeer", ctx, &domain.GetBeerParams{ID: params.Id}).Return(&domain.Beer{
		ID:      expected.Id,
		Name:    expected.Name,
		Type:    domain.IndiaPaleAle,
		Brewer:  expected.Brewer,
		Country: expected.Country,
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
		Beer:       &beers.Beer{Id: "id", Name: "name", Type: beers.BeerType_BEER_TYPE_STOUT, Brewer: "brewer", Country: "Country"},
		UpdateMask: &field_mask.FieldMask{Paths: []string{"NaMe", "type", "Country", "breweR"}},
	}
	expected := &beers.Beer{
		Id:      "id",
		Name:    "a beer",
		Type:    beers.BeerType_BEER_TYPE_STOUT,
		Brewer:  "brewer",
		Country: "country",
	}
	beerType := domain.Stout
	interactor.On("UpdateBeer", ctx, &domain.UpdateBeerParams{
		ID:      params.Beer.Id,
		Name:    &params.Beer.Name,
		Brewer:  &params.Beer.Brewer,
		Type:    &beerType,
		Country: &params.Beer.Country,
	}).Return(&domain.Beer{
		ID:      expected.Id,
		Name:    expected.Name,
		Type:    domain.Stout,
		Brewer:  expected.Brewer,
		Country: expected.Country,
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

func TestListBeers_WhenListBeersReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	const msg = "something went wrong"
	expected := status.Error(codes.Internal, msg)
	interactor.On("ListBeers", ctx, &domain.ListBeersParams{Page: 42}).Return(nil, errors.New(msg))
	_, actual := service.ListBeers(ctx, &beers.ListBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}

func TestListBeers_WhenListBeersReturnsValidationError_ReturnsInvalidArgumentError(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	const msg = "something went wrong"
	expected := status.Error(codes.InvalidArgument, msg)
	interactor.On("ListBeers", ctx, &domain.ListBeersParams{Page: 42}).Return(nil, domain.NewValidationError(msg))
	_, actual := service.ListBeers(ctx, &beers.ListBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}

func TestListBeers_WhenListBeersReturnsBeers_ReturnsBeers(t *testing.T) {
	t.Parallel()
	interactor := &mocks.BeerInteractor{}
	service := adapters.NewBeerService(interactor)
	ctx := context.Background()
	expected := &beers.ListBeersResponse{
		Beers: []*beers.Beer{
			{Id: "id1", Type: beers.BeerType_BEER_TYPE_PILSNER},
			{Id: "id2", Type: beers.BeerType_BEER_TYPE_PORTER},
			{Id: "id3", Type: beers.BeerType_BEER_TYPE_PALE_ALE},
			{Id: "id4", Type: beers.BeerType_BEER_TYPE_UNSPECIFIED},
		},
	}
	interactor.On("ListBeers", ctx, &domain.ListBeersParams{Page: 42}).Return([]*domain.Beer{
		{ID: expected.Beers[0].Id, Type: domain.Pilsner},
		{ID: expected.Beers[1].Id, Type: domain.Porter},
		{ID: expected.Beers[2].Id, Type: domain.PaleAle},
		{ID: expected.Beers[3].Id, Type: domain.Unspecified},
	}, nil)
	actual, _ := service.ListBeers(ctx, &beers.ListBeersRequest{Page: 42})
	assert.Equal(t, expected, actual)
}
