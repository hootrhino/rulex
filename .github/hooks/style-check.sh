#!/bin/sh

has_errors=0

# 获取git暂存的所有go代码
# --cached 暂存的
# --name-only 只显示名字
# --diff-filter=ACM 过滤暂存文件，A=Added C=Copied M=Modified, 即筛选出添加/复制/修改的文件
allgofiles=$(git diff --cached --name-only --diff-filter=ACM | grep '.go$')

gofiles=()
godirs=()
for allfile in ${allgofiles[@]}; do 
    # 过滤vendor的
    # 过滤prootobuf自动生产的文件
    if [[ $allfile == "vendor"* || $allfile == *".pb.go" ]];then
        continue
    else
        gofiles+=("$allfile")

        # 文件夹去重
        existdir=0
        dir=`echo "$allfile" |xargs -n1 dirname|sort -u`
        for somedir in ${godirs[@]}; do
            if [[ $dir == $somedir ]]; then 
                existdir=1
                break
            fi
        done

        if [[ $existdir -eq 0 ]]; then 
            godirs+=("$dir")
        fi
    fi
done

[ -z "$gofiles" ] && exit 0

# gofmt 格式化代码
unformatted=$(gofmt -l ${gofiles[@]})
if [ -n "$unformatted" ]; then
	echo >&2 "gofmt FAIL:\n Run following command:"
	for f in ${unformatted[@]}; do
		echo >&2 " gofmt -w $PWD/$f"
	done
	echo "\n"
	has_errors=1
fi

# goimports 自动导包
if goimports >/dev/null 2>&1; then  # 检测是否安装
	unimports=$(goimports -l ${gofiles[@]})
	if [ -n "$unimports" ]; then
		echo >&2 "goimports FAIL:\nRun following command:"
		for f in ${unimports[@]} ; do
			echo >&2 " goimports -w $PWD/$f"
		done
		echo "\n"
		has_errors=1
	fi
else
	echo 'Error: goimports not install. Run: "go get -u golang.org/x/tools/cmd/goimports"' >&2
	exit 1
fi

# golint 代码规范检测
if golint >/dev/null 2>&1; then  # 检测是否安装
	lint_errors=false
	for file in ${gofiles[@]} ; do
		lint_result="$(golint $file)" # run golint
		if test -n "$lint_result" ; then
			echo "golint '$file':\n$lint_result"
			lint_errors=true
			has_errors=1
		fi
	done
	if [ $lint_errors = true ] ; then
		echo "\n"
	fi
else
	echo 'Error: golint not install. Run: "go get -u github.com/golang/lint/golint"' >&2
	exit 1
fi

# go vet 静态错误检查
show_vet_header=true
for dir in ${godirs[@]} ; do
    vet=$(go vet $PWD/$dir 2>&1)
    if [ -n "$vet" -a $show_vet_header = true ] ; then
	echo "govet FAIL:"
	show_vet_header=false
    fi
    if [ -n "$vet" ] ; then
	echo "$vet\n"
	has_errors=1
    fi
done


exit $has_errors