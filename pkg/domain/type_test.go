package domain_test

import (
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/domain"

	"github.com/stretchr/testify/assert"
)

func TestTypeString(t *testing.T) {
	tests := []struct {
		ty  domain.Type
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
			ty:  domain.Type(-100),
			val: "Unknown",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.val, test.ty.String())
	}
}
