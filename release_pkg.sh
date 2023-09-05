#! /bin/bash
set -e

RESPOSITORY="https://github.com/hootrhino"

#
create_pkg() {
    local target=$1
    local version=$(git describe --tags --always --abbrev=0)
    local release_dir="_release"
    local pkg_name="rulex-$target-$version.zip"
    local common_files="./LICENSE ./conf/rulex.ini"
    local files_to_include="./rulex $common_files"
    local files_to_include_exe="./rulex.exe $common_files"

    if [[ "$target" != "windows" ]]; then
        files_to_include="$files_to_include ./script/*.sh"
        mv ./rulex-$target ./rulex
        chmod +x ./rulex
    else
        files_to_include="$files_to_include_exe"
        mv ./rulex-$target.exe ./rulex
    fi

    echo "Create package: $pkg_name"
    zip -j "$release_dir/$pkg_name" $files_to_include
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

build_windows() {
    make windows
}
build_x64linux() {
    make x64linux
}

build_arm64linux() {
    make arm64
}

build_arm32linux() {
    make arm32
}

build_mips64linux() {
    make mips64
}

build_mips32linux() {
    make mips32
}
#------------------------------------------
cross_compile() {
    ARCHS=("windows" "x64linux" "arm64linux" "arm32linux")
    if [ ! -d "./_release/" ]; then
        mkdir -p ./_release/
    else
        rm -rf ./_release/
        mkdir -p ./_release/
    fi
    for arch in ${ARCHS[@]}; do
        echo -e "\033[34m [★] Compile target =>\033[43;34m ["$arch"]. \033[0m"
        if [[ "${arch}" == "windows" ]]; then
            build_windows $arch
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
    local www_zip="./www.zip"
    local http_server_dir="./plugin/http_server/www"

    # 检查文件是否存在
    if [ -f "$www_zip" ]; then
        echo -e "\033[44;32m [√] File www.zip already downloaded \033[0m"
        unzip -q "$www_zip" -d "$http_server_dir"
    else
        local VERSION=$(git describe --tags --always --abbrev=0)
        local URL="${RESPOSITORY}/hootrhino-eekit-web/releases/download/${VERSION}/www.zip"
        echo -e "\033[41;37m [*] Fetch www.zip from: ${URL}\033[0m"
        # 发送HEAD请求来检查URL是否存在
        response=$(curl -s --head -w %{http_code} "$URL" -o /dev/null)

        if [ "$response" = "200" ]; then
            echo -e "\033[40;32m [*] Unzip www.zip to:${http_server_dir} \033[0m"
            wget -q --show-progress "$URL"
            unzip -q "$www_zip" -d "$http_server_dir"
        else
            echo -e "\033[41;30m [x] Error with http code 404, check if ${URL} exists \033[0m"
            exit 1
        fi
    fi
}

#
# gen_changelog
#
gen_changelog() {
    PreviewVersion=$(git describe --tags --abbrev=0 $(git rev-list --tags --skip=1 --max-count=1))
    CurrentVersion=$(git describe --tags --abbrev=0)
    echo "----------------------------------------------------------------"
    echo -e "\033[34m [★] Change log between\033[44;34m ["${PreviewVersion}"] <--> [${CurrentVersion}]. \033[0m"
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
    DEPS=("bash" "git" "jq" "gcc" "make" "x86_64-w64-mingw32-gcc")
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
