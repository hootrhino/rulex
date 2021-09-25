#! /bin/bash
# Rpc
echo "Generate Rpc Proto"
go get google.golang.org/grpc
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway
go get github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2
go get google.golang.org/protobuf/cmd/protoc-gen-go
go get google.golang.org/grpc/cmd/protoc-gen-go-grpc
#
export GOROOT=/usr/local/go
export GOPATH=$HOME/go
export GOBIN=$GOPATH/bin
export PATH=$PATH:$GOROOT:$GOPATH:$GOBIN
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ./rulexrpc/grpc_resource.proto
# Stream
echo "Generate Stream Proto"
protoc -I ./xstream \
    --go_out ./xstream --go_opt paths=source_relative \
    --go-grpc_out ./xstream --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./xstream --grpc-gateway_opt paths=source_relative \
    --openapiv2_out ./xstream \
    --openapiv2_opt logtostderr=true \
    ./xstream/xstream.proto
# Stream
echo "Generate Cloud Proto"
protoc -I ./cloud \
    --go_out ./cloud --go_opt paths=source_relative \
    --go-grpc_out ./cloud --go-grpc_opt paths=source_relative \
    --grpc-gateway_out ./cloud --grpc-gateway_opt paths=source_relative \
    --openapiv2_out ./cloud \
    --openapiv2_opt logtostderr=true \
    ./cloud/atomic_cloud.proto
