#! /bin/bash

clone() {
    git clone https://github.com/hootrhino/trailer-demo-app.git ./_temp/trailer-demo-app
    cd ./_temp/trailer-demo-app
    go get
    go build
    cd ../../
}
clone