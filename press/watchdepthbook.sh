#!/usr/bin/env bash

if [ $# -eq 0 ]; then
    echo "product required!"
    exit
fi

for((;;)) do
    okchaincli query order depthbook $1 99999999 ${CCC} |jq
    echo "---------------------------------"
    sleep 3
done
