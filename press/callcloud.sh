#!/bin/bash


SSH="ssh -i ~/okchain-dex-test.pem ubuntu"

call() {
${SSH}@${h} << eeooff
    source /home/ubuntu/.env.sh
    $@
    exit
eeooff
}

call $@