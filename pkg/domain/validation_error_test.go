package domain_test

import (
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestNewValidationError_ReturnsValidationError(t *testing.T) {
	t.Parallel()
	assert.NotNil(t, domain.NewValidationError("msg"))
}

func TestError_ReturnsErrorString(t *testing.T) {
	t.Parallel()
	expected := "msg"
	err := domain.NewValidationError(expected)
	assert.Equal(t, expected, err.Error())
}
