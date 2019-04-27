#!/usr/bin/env bash


for((;;)) do
    okchaincli query order depthbook $1 $CC |jq
    sleep 3
done
