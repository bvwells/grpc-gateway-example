package grpc_gateway_example

//go:generate protoc -I. --go_out=plugins=grpc,paths=source_relative:./ api.proto
//go:generate protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:./ api.proto
//go:generate protoc -I. --swagger_out=disable_default_errors=true,logtostderr=true:../api/openapi-spec api.proto
