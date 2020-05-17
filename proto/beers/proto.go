package beers

//go:generate protoc -I. -I../../third_party/protobuf --go_out=plugins=grpc,paths=source_relative:./ api.proto
//go:generate protoc -I. -I../../third_party/protobuf --grpc-gateway_out=logtostderr=true,paths=source_relative:./ api.proto
//go:generate protoc -I. -I../../third_party/protobuf --swagger_out=disable_default_errors=true,logtostderr=true:../../api/openapi-spec api.proto
