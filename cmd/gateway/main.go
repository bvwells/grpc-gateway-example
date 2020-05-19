package main

import (
	"context"
	"net"
	"net/http"

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
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
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
	mux := runtime.NewServeMux()
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
