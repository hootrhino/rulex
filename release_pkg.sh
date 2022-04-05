#! /bin/bash
set -e

RESPOSITORY="https://github.com/i4de"

#
create_pkg() {
    VERSION=$(git describe --tags --always --abbrev=0)
    echo "Create package: ${rulex-$1-${VERSION}}"
    if [ "$1" == "x64windows" ]; then
        cd rulexc
        go get
        CGO_ENABLED=1 GOOS=windows GOARCH=amd64 go build -v -o rulex-cli.exe main.go
        cd ../
        cp ./rulexc/rulex-cli.exe ./
        zip -r _release/rulex-$1-${VERSION}.zip \
            ./rulex-$1.exe \
            ./rulex-cli.exe \
            ./conf/rulex.ini
        rm -rf ./rulex-cli.exe
        rm -rf ./rulex-$1.exe
    else
        cd rulexc
        go get
        go build -v -o rulex-cli main.go
        cd ../
        cp ./rulexc/rulex-cli ./
        zip -r _release/rulex-$1-${VERSION}.zip \
            ./rulex-$1 \
            ./rulex-cli \
            ./script/crulex.sh \
            ./conf/rulex.ini
        rm -rf ./rulex-cli
        rm -rf ./rulex-$1
    fi

}

#
make_zip() {
    if [ -n $1 ]; then
        create_pkg $1
    else
        echo "Should have release target."
        exit 1
    fi

}

#
build_x64windows() {
    CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc \
        go build -ldflags "-s -w" -o rulex-$1.exe main.go
}
build_x86linux() {
    CGO_ENABLED=1 GOOS=linux GO386=softfloat \
        go build -ldflags "-s -w" -o rulex-$1 main.go
}
build_x64linux() {
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
        go build -ldflags "-s -w" -o rulex-$1 main.go
}

build_arm64linux() {
    CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
        go build -ldflags "-s -w" -o rulex-$1 main.go
}
build_arm32linux() {
    CGO_ENABLED=1 GOARM=7 GOOS=linux GOARCH=arm CC=arm-linux-gnueabi-gcc \
        go build -ldflags "-s -w" -o rulex-$1 -ldflags "-linkmode external -extldflags -static" main.go
}
build_arm64android() {
    CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc \
        go build -ldflags "-s -w" -o rulex-$1 main.go
}
build_x64android() {
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64 \
        go build -ldflags "-s -w" -o rulex-$1 main.go
}

#------------------------------------------
cross_compile() {
    ARCHS=("x64windows" "x86linux" "x64linux" "arm64linux" "arm32linux" "arm64android" "x64android")
    if [ ! -d "./_release/" ]; then
        mkdir -p ./_release/
    else
        rm -rf ./_release/
        mkdir -p ./_release/
    fi
    for arch in ${ARCHS[@]}; do
        echo "Compile target => ["$arch"]"
        if [[ "${arch}" == "x64windows" ]]; then
            # sudo apt install gcc-mingw-w64-x86-64 -y
            build_x64windows $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."
        fi
        if [[ "${arch}" == "x86linux" ]]; then
            build_x86linux $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."

        fi
        if [[ "${arch}" == "x64linux" ]]; then
            build_x64linux $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."

        fi
        if [[ "${arch}" == "arm64linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm64linux $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."

        fi
        if [[ "${arch}" == "arm32linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm32linux $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."

        fi
        if [[ "${arch}" == "arm64android" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm64android $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."

        fi
        if [[ "${arch}" == "x64android" ]]; then
            build_x64android $arch
            make_zip $arch
            echo "Compile target => ["$arch"] Ok."
        fi
    done
}
#
#
#
build_rulexcli() {
    git clone ${RESPOSITORY}/rulex-cli.git rulexc
}
#
# fetch dashboard
#
fetch_dashboard() {
    VERSION=$(git describe --tags --always --abbrev=0)
    wget -q --show-progress ${RESPOSITORY}/rulex-dashboard/releases/download/${VERSION}/${VERSION}.zip
    unzip -q ${VERSION}.zip
    cp -r ./dist/* ./plugin/http_server/www
}
#
# gen_changelog
#
gen_changelog() {
    PreviewVersion=$(git describe --tags --abbrev=0 $(git rev-list --tags --skip=1 --max-count=1))
    CurrentVersion=$(git describe --tags --abbrev=0)
    echo "----------------------------------------------------------------"
    echo "|  Change log between [${PreviewVersion}] <--> [${CurrentVersion}]"
    echo "----------------------------------------------------------------"
    ChangeLog=$(git log ${PreviewVersion}..${CurrentVersion} --oneline --decorate --color)
    printf "${ChangeLog}\n"

}

#
#-----------------------------------
#
init_env() {
    if [ ! -d "./_build/" ]; then
        mkdir -p ./_build/
    else
        rm -rf ./_build/
        mkdir -p ./_build/
    fi

}
#
#-----------------------------------
#
init_env
cp -r $(ls | egrep -v '^_build$') ./_build/
cd ./_build/
build_rulexcli
fetch_dashboard
cross_compile
gen_changelog
