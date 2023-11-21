#! /bin/bash
set -e
APP=rulex
RESPOSITORY="https://github.com/hootrhino"

#
create_pkg() {
    local target=$1
    local version="$(git describe --tags $(git rev-list --tags --max-count=1))"
    local release_dir="_release"
    local pkg_name="${APP}-$target-$version.zip"
    local common_files="./conf/license.* ./LICENSE ./conf/${APP}.ini ./md5.sum"
    local files_to_include="./${APP} $common_files"
    local files_to_include_exe="./${APP}.exe $common_files"

    if [[ "$target" != "windows" ]]; then
        files_to_include="$files_to_include ./script/*.sh"
        mv ./${APP}-$target ./${APP}
        chmod +x ./${APP}
        calculate_and_save_md5 ./${APP}
    else
        files_to_include="$files_to_include_exe ./script/*.bat"
        mv ./${APP}-$target.exe ./${APP}.exe
        calculate_and_save_md5 ./${APP}.exe
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
# 计算md5
calculate_and_save_md5() {
    if [ $# -ne 1 ]; then
        echo "Usage: $0 <file_path>"
        exit 1
    fi
    local file_path="$1"
    local md5_hash
    if [ ! -f "$file_path" ]; then
        echo "File not found: $file_path"
        return 1
    fi
    md5_hash=$(md5sum "$file_path" | awk '{print $1}')
    echo -n "$md5_hash" > md5.sum
}
#
# fetch dashboard
#
#!/bin/bash

fetch_dashboard() {
    local owner="hootrhino"
    local repo="hootrhino-eekit-web"

    # 检查当前目录是否已经存在 www.zip 文件
    if [ -f "www.zip" ]; then
        echo "[!] www.zip already exists. No need to download."
        exit 0
    fi

    # 获取最新 release 的 tag 名称
    local tag=$(curl -s "https://api.github.com/repos/$owner/$repo/releases/latest" | jq -r .tag_name)

    # 获取最新 release 中的 www.zip 下载链接
    local zip_url=$(curl -s "https://api.github.com/repos/$owner/$repo/releases/latest" | jq -r '.assets[] | select(.name == "www.zip") | .browser_download_url')

    if [ -z "$zip_url" ]; then
        echo "[x] Error: www.zip not found in the release assets."
        exit 1
    fi

    # 下载 www.zip 文件
    curl -L -o www.zip "$zip_url"

    echo "[√] Download complete. Tag: $tag"

    # 解压 www.zip 文件到指定目录
    unzip -o www.zip -d /plugin/rulex_api_server/server/www/

    echo "[√] Extraction complete. www.zip contents have been overwritten to /plugin/rulex_api_server/server/www/."
}

#
# gen_changelog
#

gen_changelog() {
    echo -e "[.]Version Change log:"
    log=$(git log --oneline --pretty=format:" \033[0;31m[*]\033[0m%s\n" $(git describe --abbrev=0 --tags).. | cat)
    echo -e $log
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
