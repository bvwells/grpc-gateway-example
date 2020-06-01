package usecases_test

import (
	"context"
	"errors"
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"
	"github.com/bvwells/grpc-gateway-example/pkg/usecases"
	"github.com/bvwells/grpc-gateway-example/pkg/usecases/mocks"

	"github.com/stretchr/testify/assert"
)

//go:generate mockery -name=BeerRepository -case=underscore

func TestNewBeerInteractor_ReturnsBeerInteractor(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	assert.NotNil(t, usecases.NewBeerInteractor(repo))
}

func TestCreateBeer_WhenValidateReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	_, err := interactor.CreateBeer(context.Background(), &domain.CreateBeerParams{})
	assert.NotNil(t, err)
}

func TestCreateBeer_WhenCreateBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.CreateBeerParams{Name: "a beer"}
	expected := errors.New("something went wrong")
	repo.On("CreateBeer", ctx, params).Return(nil, expected)
	_, actual := interactor.CreateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestCreateBeer_WhenCreateBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.CreateBeerParams{Name: "a beer"}
	expected := &domain.Beer{ID: "id"}
	repo.On("CreateBeer", ctx, params).Return(expected, nil)
	actual, _ := interactor.CreateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeer_WhenValidateReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	_, err := interactor.GetBeer(context.Background(), &domain.GetBeerParams{})
	assert.NotNil(t, err)
}

func TestGetBeer_WhenGetBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.GetBeerParams{ID: "ID"}
	expected := errors.New("something went wrong")
	repo.On("GetBeer", ctx, params).Return(nil, expected)
	_, actual := interactor.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeer_WhenGetBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.GetBeerParams{ID: "ID"}
	expected := &domain.Beer{ID: "id"}
	repo.On("GetBeer", ctx, params).Return(expected, nil)
	actual, _ := interactor.GetBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenValidateReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	_, err := interactor.UpdateBeer(context.Background(), &domain.UpdateBeerParams{})
	assert.NotNil(t, err)
}

func TestUpdateBeer_WhenUpdateBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	brewer := "brewer"
	params := &domain.UpdateBeerParams{ID: "ID", Brewer: &brewer}
	expected := errors.New("something went wrong")
	repo.On("UpdateBeer", ctx, params).Return(nil, expected)
	_, actual := interactor.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestUpdateBeer_WhenUpdateBeerReturnsBeer_ReturnsBeer(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	brewer := "brewer"
	params := &domain.UpdateBeerParams{ID: "ID", Brewer: &brewer}
	expected := &domain.Beer{ID: "id"}
	repo.On("UpdateBeer", ctx, params).Return(expected, nil)
	actual, _ := interactor.UpdateBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestDeleteBeer_WhenValidateReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	err := interactor.DeleteBeer(context.Background(), &domain.DeleteBeerParams{})
	assert.NotNil(t, err)
}

func TestDeleteBeer_WhenDeleteBeerReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.DeleteBeerParams{ID: "ID"}
	expected := errors.New("something went wrong")
	repo.On("DeleteBeer", ctx, params).Return(expected)
	actual := interactor.DeleteBeer(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestDeleteBeer_WhenDeleteBeerReturnsNilReturnsBeer(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.DeleteBeerParams{ID: "ID"}
	repo.On("DeleteBeer", ctx, params).Return(nil)
	actual := interactor.DeleteBeer(ctx, params)
	assert.Nil(t, actual)
}

func TestGetBeers_WhenValidateReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	_, err := interactor.GetBeers(context.Background(), &domain.GetBeersParams{})
	assert.NotNil(t, err)
}

func TestGetBeers_WhenGetBeersReturnsError_ReturnsError(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.GetBeersParams{Page: 1}
	expected := errors.New("something went wrong")
	repo.On("GetBeers", ctx, params).Return(nil, expected)
	_, actual := interactor.GetBeers(ctx, params)
	assert.Equal(t, expected, actual)
}

func TestGetBeers_WhenGetBeersReturnsBeers_ReturnsBeers(t *testing.T) {
	t.Parallel()
	repo := &mocks.BeerRepository{}
	interactor := usecases.NewBeerInteractor(repo)
	ctx := context.Background()
	params := &domain.GetBeersParams{Page: 1}
	expected := []*domain.Beer{
		{ID: "id1"},
		{ID: "id2"},
	}
	repo.On("GetBeers", ctx, params).Return(expected, nil)
	actual, _ := interactor.GetBeers(ctx, params)
	assert.Equal(t, expected, actual)
}
