package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/textproto"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// NewProtoErrorHandler returns a new runtime.ProtoErrorHandlerFunc and
// illustrates how custom responses can be returned from the grpc gateway.
func NewProtoErrorHandler(logger *logrus.Logger) runtime.ProtoErrorHandlerFunc {
	return func(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
		s, ok := status.FromError(err)
		if !ok {
			s = status.New(codes.Unknown, err.Error())
		}

		st := runtime.HTTPStatusFromCode(s.Code())

		type Error struct {
			Code    int32  `json:"code"`
			Message string `json:"message"`
		}
		httpError := &Error{
			Code:    int32(st),
			Message: s.Message(),
		}
		w.Header().Set("Content-Type", "application/json")
		buf, err := json.Marshal(httpError)
		if err != nil {
			logger.Infof("failed to marshall error response: %v", err)
			return
		}
		w.WriteHeader(st)
		if _, err := w.Write(buf); err != nil {
			logger.Infof("failed to write response: %v", err)
			return
		}
	}
}

// NewHeaderMatcher returns a new runtime.HeaderMatcherFunc and illustrates
// how to match request headers and add them to the grpc metadata.
func NewHeaderMatcher() runtime.HeaderMatcherFunc {
	return func(key string) (string, bool) {
		key = textproto.CanonicalMIMEHeaderKey(key)
		if key == "X-Request-Id" {
			return key, true
		}

		return "", false
	}
}

// NewAnnotator returns a new grpc metadata annotator and illustrates
// how custom data can be added to the grpc metadata request context.
func NewAnnotator() func(context.Context, *http.Request) metadata.MD {
	return func(ctx context.Context, req *http.Request) metadata.MD {
		xRequestID := req.Header.Get("X-Request-Id")
		return metadata.New(map[string]string{"corelationID": xRequestID})
	}
}
