#! /bin/bash
# Rpc
echo "Generate Rpc Proto"
go get -u google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
#
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN

echo "Generate RulexRpc InEnd Proto"
protoc -I ./rulexrpc --go_out=./rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./rulexrpc --go-grpc_opt paths=source_relative \
    ./rulexrpc/grpc_resource.proto
# Stream
echo "Generate Stream Proto"
protoc -I ./xstream --go_out ./xstream --go_opt paths=source_relative \
    --go-grpc_out=./xstream --go-grpc_opt paths=source_relative \
    ./xstream/xstream.proto
