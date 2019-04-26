#!/bin/bash

function push_build {
ssh root@192.168.13.116 << eeooff
    source ~/env.sh
    cd /root/go/src/github.com/ok-chain/okchain
    git stash
    git pull
    make install

    cd /root/go/src/github.com/cosmos/launch
    git stash
    git pull
    go build
    cd /root/go/src/github.com/cosmos/launch/systemctl/binary/
    cp /usr/local/go/bin/okchaind .
    cp /usr/local/go/bin/okchaincli .
    cp /root/go/src/github.com/cosmos/launch/launch .
    ./zip.sh
    ./gitpush.sh
eeooff
echo done!
}


function main {
    push_build
}

main