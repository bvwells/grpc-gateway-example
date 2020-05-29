package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/textproto"

	"github.com/bvwells/grpc-gateway-example/pkg/adapters"
	"github.com/bvwells/grpc-gateway-example/pkg/infrastructure"
	"github.com/bvwells/grpc-gateway-example/pkg/usecases"
	"github.com/bvwells/grpc-gateway-example/proto/beers"
	gw "github.com/bvwells/grpc-gateway-example/proto/beers"

	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/joonix/log"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	address = "127.0.0.1:50000"
)

func newLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(log.NewFormatter())
	return logger
}

func newBeerService() (*adapters.BeerService, error) {
	settings := &infrastructure.PostgresSettings{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "docker",
		DBName:   "beers",
	}

	generateID := func() string {
		return uuid.New().String()
	}
	repo, err := infrastructure.NewPostgresBeerRepository(settings, generateID)
	if err != nil {
		return nil, err
	}

	interactor := usecases.NewBeerInteractor(repo)
	return adapters.NewBeerService(interactor), nil
}

func main() {
	logger := newLogger()
	service, err := newBeerService()
	if err != nil {
		logger.Fatalf("error creating beer service: %v", err)
	}

	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("error listening on port '%s': %v", address, err)
	}

	s := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logger)),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logger)),
			grpc_recovery.StreamServerInterceptor(),
		),
	)
	grpc_health_v1.RegisterHealthServer(s, health.NewServer())
	beers.RegisterBeerServiceServer(s, service)

	logger.Infof("starting gRPC service at '%s'", address)
	go func() {
		logger.Fatalf("error serving gRPC server: %v", s.Serve(lis))
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		address,
		grpc.WithInsecure(),
		grpc.WithBlock(),
	)
	if err != nil {
		logger.Fatalf("error dialing gRPC server: %v", err)
	}

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux(
		runtime.WithProtoErrorHandler(NewProtoErrorHandler(logger)),
		runtime.WithIncomingHeaderMatcher(NewHeaderMatcher()),
		runtime.WithMetadata(NewAnnotator()),
	)
	err = gw.RegisterBeerServiceHandler(context.Background(), mux, conn)
	if err != nil {
		logger.Fatalf("error registering beer service handler: %v", err)
	}

	logger.Info("starting http service at ':8080'")

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	err = http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Fatalf("error serving beer service: %v", err)
	}
}

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
		}
		w.WriteHeader(st)
		if _, err := w.Write(buf); err != nil {
			logger.Infof("failed to write response: %v", err)
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
