#!/bin/bash

# Install dependencies
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/mwitkow/go-proto-validators/protoc-gen-govalidators@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest

# Generate Stubs
protoc -I ./proto \
  --proto_path ./proto/lib \
  --plugin=$(go env GOPATH)/bin/protoc-gen-go \
  --plugin=$(go env GOPATH)/bin/protoc-gen-govalidators \
  --plugin=$(go env GOPATH)/bin/protoc-gen-grpc-gateway \
  --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
  --go_out ./stubs --go_opt paths=source_relative \
  --go-grpc_out ./stubs --go-grpc_opt paths=source_relative \
  --grpc-gateway_out ./stubs --grpc-gateway_opt paths=source_relative \
  ./proto/*/*.proto

# Generate Open API V2
protoc -I ./proto \
  --proto_path ./proto/lib \
  --plugin=$(go env GOPATH)/bin/protoc-gen-grpc-gateway \
  --plugin=$(go env GOPATH)/bin/protoc-gen-openapiv2 \
  --plugin=$(go env GOPATH)/bin/protoc-gen-go-grpc \
  --openapiv2_out ./api \
  --openapiv2_opt logtostderr=true \
  --openapiv2_opt use_go_templates=true \
  ./proto/*/*_service.proto
  