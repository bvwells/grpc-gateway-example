package domain_test

import (
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestBeerTypeString(t *testing.T) {
	tests := []struct {
		ty  domain.BeerType
		val string
	}{
		{
			ty:  domain.Ale,
			val: "Ale",
		},
		{
			ty:  domain.Bitter,
			val: "Bitter",
		},
		{
			ty:  domain.Lager,
			val: "Lager",
		},
		{
			ty:  domain.IndiaPaleAle,
			val: "IndiaPaleAle",
		},
		{
			ty:  domain.Stout,
			val: "Stout",
		},
		{
			ty:  domain.Pilsner,
			val: "Pilsner",
		},
		{
			ty:  domain.Porter,
			val: "Porter",
		},
		{
			ty:  domain.PaleAle,
			val: "PaleAle",
		},
		{
			ty:  domain.BeerType(-100),
			val: "Unspecified",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.val, test.ty.String())
	}
}
