#! /bin/bash
# Rpc
echo ">>> Generate GRpc Proto"
go get -u google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
#
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
# RulexRpc
echo ">>> Generate RulexRpc Proto"
protoc -I ./rulexrpc --go_out=./rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./rulexrpc --go-grpc_opt paths=source_relative \
    ./rulexrpc/grpc_resource.proto
# Stream
echo ">>> Generate XStream Proto"
protoc -I ./xstream --go_out ./xstream --go_opt paths=source_relative \
    --go-grpc_out=./xstream --go-grpc_opt paths=source_relative \
    ./xstream/xstream.proto

echo ">>> Generate Rpc Proto OK."