#!/bin/sh

# 忽略merge request
MERGE_MSG=`cat $1 | egrep '^Merge branch*'`

if [ "$MERGE_MSG" != "" ]; then
    exit 0
fi

COMMIT_MSG=`cat $1 | egrep "^(feat|fix|docs|chore|dev|ci|test|refactor|optmize)(\(\w+\))?:\s(\S|\w)+"`

if [ "$COMMIT_MSG" = "" ]; then
    echo "Commit Message Irregular，Please check!\n"
    exit 1
fi

if [ ${#COMMIT_MSG} -lt 15 ]; then
    echo "Commit Message Too Short，Please show me more detail!\n"
    exit 1
fi