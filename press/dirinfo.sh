#!/usr/bin/env bash

DIR_NAME=~/.okchaind/data/application.db
SAMPLING=30

if [ $# -gt 0 ]; then
    DIR_NAME=$1
fi

if [ $# -gt 1 ]; then
    SAMPLING=$2
fi

echo "Start to monitor ${DIR_NAME}..."

index=0
for((;;))
do
    sdate=`date +"%Y-%m-%d %H:%M:%S"`
    smem=`du -ms ${DIR_NAME}`
    ((index++))
    echo $index "date=["$sdate"], dir_info["$smem"]"
    sleep ${SAMPLING}
done
