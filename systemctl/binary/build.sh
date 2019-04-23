#!/bin/bash

function pull_build {
ssh root@192.168.13.116 << eeooff
    source ~/env.sh
    cd /root/go/src/github.com/ok-chain/okchain
    git pull
    make install

    cd /root/go/src/github.com/cosmos/launch
    git pull
    go build
eeooff
echo done!
}

function getbinary {
    scp root@192.168.13.116:/usr/local/go/bin/okchaind ~/go/src/github.com/cosmos/launch/systemctl/binary/
    scp root@192.168.13.116:/usr/local/go/bin/okchaincli  ~/go/src/github.com/cosmos/launch/systemctl/binary/
    scp root@192.168.13.116:/root/go/src/github.com/cosmos/launch/launch  ~/go/src/github.com/cosmos/launch/systemctl/binary/

    cp -r ~/go/src/github.com/cosmos/launch/systemctl/binary ~/go/launch/systemctl/
}

function main {
    pull_build
    getbinary
}

main