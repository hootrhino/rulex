#! /bin/bash
set -e

RESPOSITORY="https://github.com/i4de"

#
create_pkg() {
    VERSION=$(git describe --tags --always --abbrev=0)
    echo "Create package: ${rulex-$1-${VERSION}}"
    if [ "$1" == "x64windows" ]; then
        zip -r _release/rulex-$1-${VERSION}.zip \
            ./rulex-$1.exe \
            ./conf/rulex.ini
        rm -rf ./rulex-$1.exe
    else
        zip -r _release/rulex-$1-${VERSION}.zip \
            ./rulex-$1 \
            ./script/crulex.sh \
            ./conf/rulex.ini
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
        go build -ldflags "-s -w -linkmode external -extldflags -static" -o rulex-$1 main.go
}

#------------------------------------------
cross_compile() {
    ARCHS=("x64windows" "x64linux" "arm64linux" "arm32linux")
    if [ ! -d "./_release/" ]; then
        mkdir -p ./_release/
    else
        rm -rf ./_release/
        mkdir -p ./_release/
    fi
    for arch in ${ARCHS[@]}; do
        echo -e "\033[34m [★] Compile target =>\033[43;34m ["$arch"]. \033[0m"
        if [[ "${arch}" == "x64windows" ]]; then
            # sudo apt install gcc-mingw-w64-x86-64 -y
            build_x64windows $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"
        fi
        if [[ "${arch}" == "x86linux" ]]; then
            build_x86linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
        if [[ "${arch}" == "x64linux" ]]; then
            build_x64linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
        if [[ "${arch}" == "arm64linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm64linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
        if [[ "${arch}" == "arm32linux" ]]; then
            # sudo apt install gcc-arm-linux-gnueabi -y
            build_arm32linux $arch
            make_zip $arch
            echo -e "\033[33m [√] Compile target => ["$arch"] Ok. \033[0m"

        fi
    done
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
# 检查是否安装了这些软件
#
check_cmd() {
    DEPS=("git" "jq" "gcc" "make")
    for dep in ${DEPS[@]}; do
        echo -e "\033[34m [*] Check dependcy command: $dep. \033[0m"
        if ! [ -x "$(command -v $dep)" ]; then
            echo -e "\033[31m |x| Error: $dep is not installed. \033[0m"
            exit 1
        else
            echo -e "\033[32m [√] $dep has been installed. \033[0m"
        fi
    done

}
#
#-----------------------------------
#
check_cmd
init_env
cp -r $(ls | egrep -v '^_build$') ./_build/
cd ./_build/
# fetch_dashboard
cross_compile
gen_changelog
