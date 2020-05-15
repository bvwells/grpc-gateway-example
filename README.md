# gRPC Gateway Example

[![Build Status](https://travis-ci.org/bvwells/grpc-gateway-example.svg?branch=master)](https://travis-ci.org/bvwells/grpc-gateway-example)
[![codecov](https://codecov.io/gh/bvwells/grpc-gateway-example/branch/master/graph/badge.svg)](https://codecov.io/gh/bvwells/grpc-gateway-example)
[![Go Report Card](https://goreportcard.com/badge/github.com/bvwells/grpc-gateway-example)](https://goreportcard.com/report/github.com/bvwells/grpc-gateway-example)

This repo contains an example usage of grpc gateway (https://github.com/grpc-ecosystem/grpc-gateway).

## Developing

Install protoc (https://github.com/protocolbuffers/protobuf)

```
brew install protobuf
```

Run 

```
go mod tidy
```

```
go install \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger \
    github.com/golang/protobuf/protoc-gen-go
```

Add $GOBIN to path

```
export PATH=$PATH:~/go/bin
```

## Generate gRPC server/client stub

```
protoc -I. --go_out=plugins=grpc,paths=source_relative:./ api.proto
```

## Generate reverse-proxy using protoc-gen-grpc-gateway

```
protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:./ api.proto 
```

## Generate swagger definitions using protoc-gen-swagger

```
protoc -I. --swagger_out=disable_default_errors=true,logtostderr=true:../api/openapi-spec api.proto
```


## TODOs

- What to do with request headers?
- Patch request.
- How to remove 200 response from delete resource request.
