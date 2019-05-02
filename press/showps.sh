#!/usr/bin/env bash


ps() {
    echo "===============$1================"
    h=$1 ./callcloud.sh ps -ef|grep bash
}


ps c13
ps c16
ps c21
ps c22
ps c23
ps c24
ps c25
