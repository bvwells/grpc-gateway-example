package adapters_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/adapters"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/metadata"
)

func TestNewHeaderMatcher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		header          string
		canonicalHeader string
		allowed         bool
	}{
		{
			name:            "header is canonical x-request-id",
			header:          "X-Request-Id",
			canonicalHeader: "X-Request-Id",
			allowed:         true,
		},
		{
			name:            "header is not canonical x-request-id",
			header:          "x-request-id",
			canonicalHeader: "X-Request-Id",
			allowed:         true,
		},
		{
			name:            "header is not allowed",
			header:          "not-allowed",
			canonicalHeader: "",
			allowed:         false,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			canonicalHeader, allowed := adapters.NewHeaderMatcher()(test.header)
			assert.Equal(s, test.canonicalHeader, canonicalHeader)
			assert.Equal(s, test.allowed, allowed)
		})
	}
}

func TestNewAnnotator(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		req  *http.Request
		md   metadata.MD
	}{
		{
			name: "request with required header",
			req: &http.Request{
				Header: http.Header{
					"X-Request-Id": []string{"id"},
				},
			},
			md: metadata.New(map[string]string{"corelationID": "id"}),
		},
		{
			name: "request with required header missing",
			req: &http.Request{
				Header: http.Header{},
			},
			md: metadata.New(map[string]string{"corelationID": ""}),
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("test %s", test.name), func(s *testing.T) {
			md := adapters.NewAnnotator()(context.Background(), test.req)
			assert.Equal(s, test.md, md)
		})
	}
}
