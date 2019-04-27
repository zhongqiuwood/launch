#!/usr/bin/env bash


for((;;)) do
    okchaincli query order depthbook $1 99999999 $CC |jq
    echo "---------------------------------"
    sleep 3
done
