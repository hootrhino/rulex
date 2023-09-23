#! /bin/bash
# set Env path
export GOROOT=/usr/bin/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
# Install protoc
# go get -u google.golang.org/grpc
# go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
# go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
# RulexRpc
echo ">>> Generate RulexRpc Proto"
protoc -I ./rulexrpc --go_out=./rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./rulexrpc --go-grpc_opt paths=source_relative \
    ./rulexrpc/grpc_source.proto
echo ">>> Generate RulexRpc Proto OK"

# Stream
echo ">>> Generate XStream Proto"
protoc -I ./xstream --go_out ./xstream --go_opt paths=source_relative \
    --go-grpc_out=./xstream --go-grpc_opt paths=source_relative \
    ./xstream/xstream.proto
echo ">>> Generate Rpc Proto OK."
# Codec
echo ">>> Generate Codec Proto."
protoc -I ./rulexrpc --go_out ./rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./rulexrpc --go-grpc_opt paths=source_relative \
    ./rulexrpc/xcodec.proto
echo ">>> Generate Codec Proto OK."
# Trailer
echo ">>> Generate Trailer Proto."
protoc -I ./component/trailer --go_out ./component/trailer --go_opt paths=source_relative \
    --go-grpc_out=./component/trailer --go-grpc_opt paths=source_relative \
    ./component/trailer/trailer.proto
echo ">>> Generate Trailer Proto OK."
