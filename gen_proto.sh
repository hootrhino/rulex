#! /bin/bash
# set Env path
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
# Install protoc
go get -u google.golang.org/grpc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
echo -e "\033[42;33m>>>\033[0m [BEGIN]"
# RulexRpc
echo ">>> Generating RulexRpc Proto..."
protoc -I ./component/rulexrpc --go_out=./component/rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./component/rulexrpc --go-grpc_opt paths=source_relative \
    ./component/rulexrpc/grpc_source.proto
echo ">>> Generate RulexRpc Proto OK"

# Stream
echo ">>> Generating XStream Proto..."
protoc -I ./component/xstream --go_out ./component/xstream --go_opt paths=source_relative \
    --go-grpc_out=./component/xstream --go-grpc_opt paths=source_relative \
    ./component/xstream/xstream.proto
echo ">>> Generate XStream Proto OK."
# Codec
echo ">>> Generating Codec Proto..."
protoc -I ./component/rulexrpc --go_out ./component/rulexrpc --go_opt paths=source_relative \
    --go-grpc_out=./component/rulexrpc --go-grpc_opt paths=source_relative \
    ./component/rulexrpc/xcodec.proto
echo ">>> Generate Codec Proto OK."
# Trailer
echo ">>> Generating Trailer Proto..."
protoc -I ./component/trailer --go_out ./component/trailer --go_opt paths=source_relative \
    --go-grpc_out=./component/trailer --go-grpc_opt paths=source_relative \
    ./component/trailer/trailer.proto
echo ">>> Generate Trailer Proto OK."

echo -e "\033[42;33m>>>\033[0m [FINISHED]"
