package domain_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestBeerValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		beer *domain.Beer
		err  error
	}{
		{
			name: "all good",
			beer: &domain.Beer{ID: "ID"},
			err:  nil,
		},
		{
			name: "missing id field",
			beer: &domain.Beer{},
			err:  errors.New("beer ID is empty"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.beer.Validate())
		})
	}
}

func TestCreateBeerParamsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *domain.CreateBeerParams
		err    error
	}{
		{
			name:   "all good",
			params: &domain.CreateBeerParams{Name: "name"},
			err:    nil,
		},
		{
			name:   "missing name field",
			params: &domain.CreateBeerParams{},
			err:    errors.New("beer name is empty"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.params.Validate())
		})
	}
}

func TestGetBeerParamsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *domain.GetBeerParams
		err    error
	}{
		{
			name:   "all good",
			params: &domain.GetBeerParams{ID: "id"},
			err:    nil,
		},
		{
			name:   "missing id field",
			params: &domain.GetBeerParams{},
			err:    errors.New("beer ID is empty"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.params.Validate())
		})
	}
}

func TestUpdateBeerParamsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *domain.UpdateBeerParams
		err    error
	}{
		{
			name:   "all good",
			params: &domain.UpdateBeerParams{ID: "id"},
			err:    nil,
		},
		{
			name:   "missing id field",
			params: &domain.UpdateBeerParams{},
			err:    errors.New("beer ID is empty"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.params.Validate())
		})
	}
}

func TestDeleteBeerParamsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *domain.DeleteBeerParams
		err    error
	}{
		{
			name:   "all good",
			params: &domain.DeleteBeerParams{ID: "id"},
			err:    nil,
		},
		{
			name:   "missing id field",
			params: &domain.DeleteBeerParams{},
			err:    errors.New("beer ID is empty"),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.params.Validate())
		})
	}
}

func TestGetBeersParamsValidate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		params *domain.GetBeersParams
		err    error
	}{
		{
			name:   "all good",
			params: &domain.GetBeersParams{},
			err:    nil,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			assert.Equal(s, test.err, test.params.Validate())
		})
	}
}
