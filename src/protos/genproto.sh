#!/bin/sh

# See README.md on how to set up the protobuf compiler and Go plugin

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative *.proto