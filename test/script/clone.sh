#! /bin/bash

clone() {
    git clone https://github.com/i4de/grpc_driver_hello_go.git ./_temp/grpc_driver_hello_go
    cd ./_temp/grpc_driver_hello_go
    go get
    go build
    cd ../../
}
clone