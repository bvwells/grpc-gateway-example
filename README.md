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
    github.com/golang/protobuf/protoc-gen-go \
    golang.org/x/tools/cmd/stringer \
    github.com/vektra/mockery/cmd/mockery \
    github.com/golangci/golangci-lint/cmd/golangci-lint
```

Add $GOBIN to path

```
export PATH=$PATH:~/go/bin
```

## Code generation 

From root of repo run:

```
go generate ./...
```

This will re-generate required code. Go generate will generate grpc server/client
stubs, the grpc gateway reverse proxy and open api definition for the protobuf
definition. These commands can be run with the following commands:

### Generate gRPC server/client stub

```
protoc -I. --go_out=plugins=grpc,paths=source_relative:./ api.proto
```

### Generate reverse-proxy using protoc-gen-grpc-gateway

```
protoc -I. --grpc-gateway_out=logtostderr=true,paths=source_relative:./ api.proto 
```

### Generate swagger definitions using protoc-gen-swagger

```
protoc -I. --swagger_out=disable_default_errors=true,logtostderr=true:../api/openapi-spec api.proto
```

## Run PostgreSQL database

To install PostgreSQL (https://www.postgresql.org/) run:

```
brew install postgresql
```

Run the docker image postgres (https://hub.docker.com/_/postgres):

```
docker run --rm --name beers -e POSTGRES_PASSWORD=docker -p 5432:5432 -v $HOME/Git/github.com/bvwells/grpc-gateway-example/pkg/infrastructure/postres:/var/lib/postgresql/data postgres
```

NOTE: The environment variable POSTGRES_PASSWORD should be set to a secret when running in a real environment.

Connect to running instance:

```
psql -h localhost -U postgres -d postgres
```

To create the beers database run:
```
CREATE DATABASE BEERS;
```

Connect to the database:
```
\c beers
```

To create the beers table run:
```
CREATE TABLE BEERS (
  id VARCHAR(36) PRIMARY KEY,
  name TEXT,
  type INT,
  brewer TEXT,
  country TEXT
);
```

Some useful psql commands:

List databases:
```
\l 
```

Connect to a database:
```
\c table_name username

```
List tables:
```
\dt 
```

Describe a table:
```
\d table_name
```

Help:
```
\? table_name
```

## TODOs

- What to do with request headers?
- Implement paging for get all beers.
