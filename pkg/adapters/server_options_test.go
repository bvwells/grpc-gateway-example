package adapters_test

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bvwells/grpc-gateway-example/pkg/adapters"

	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestNewProtoErrorHandler_SetsResponseContentType(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()

	handler := adapters.NewProtoErrorHandler(logger)

	w := httptest.NewRecorder()
	err := status.Error(codes.Internal, "something went wrong")
	handler(nil, nil, nil, w, nil, err)

	assert.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
}

func TestNewProtoErrorHandler_IfCalledWithStatusError_ReturnMappedHTTPError(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()

	handler := adapters.NewProtoErrorHandler(logger)

	w := httptest.NewRecorder()
	err := status.Error(codes.NotFound, "something went wrong")
	handler(nil, nil, nil, w, nil, err)

	assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
}

func TestNewProtoErrorHandler_IfCalledWithError_ReturnInternalServerError(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()

	handler := adapters.NewProtoErrorHandler(logger)

	w := httptest.NewRecorder()
	err := errors.New("something went wrong")
	handler(nil, nil, nil, w, nil, err)

	assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
}

func TestNewProtoErrorHandler_ReturnsErrorBody(t *testing.T) {
	t.Parallel()

	logger, _ := test.NewNullLogger()

	handler := adapters.NewProtoErrorHandler(logger)

	w := httptest.NewRecorder()
	err := status.Error(codes.NotFound, "something went wrong")
	handler(nil, nil, nil, w, nil, err)

	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, "{\"code\":404,\"message\":\"something went wrong\"}", string(body))
}

func TestNewIncomingHeaderMatcher(t *testing.T) {
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
			canonicalHeader, allowed := adapters.NewIncomingHeaderMatcher()(test.header)
			assert.Equal(s, test.canonicalHeader, canonicalHeader)
			assert.Equal(s, test.allowed, allowed)
		})
	}
}

func TestNewOutgoingHeaderMatcher(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name            string
		header          string
		canonicalHeader string
		allowed         bool
	}{
		{
			name:            "header is canonical Content-Type",
			header:          "Content-Type",
			canonicalHeader: "Content-Type",
			allowed:         true,
		},
		{
			name:            "header is not canonical Content-Type",
			header:          "content-type",
			canonicalHeader: "Content-Type",
			allowed:         true,
		},
		{
			name:            "header is canonical Content-Length",
			header:          "Content-Length",
			canonicalHeader: "Content-Length",
			allowed:         true,
		},
		{
			name:            "header is not canonical Content-Length",
			header:          "content-length",
			canonicalHeader: "Content-Length",
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
			canonicalHeader, allowed := adapters.NewOutgoingHeaderMatcher()(test.header)
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
