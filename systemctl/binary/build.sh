#!/bin/bash

function pull_build {
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
    cp /usr/local/go/bin/okchaind .
    cp /root/go/src/github.com/cosmos/launch/launch .
    ./zip.sh
    ./gitpush
eeooff
echo done!
}

function getbinary {
    scp root@192.168.13.116:/usr/local/go/bin/okchaind ~/go/src/github.com/cosmos/launch/systemctl/binary/
    scp root@192.168.13.116:/usr/local/go/bin/okchaincli  ~/go/src/github.com/cosmos/launch/systemctl/binary/
    scp root@192.168.13.116:/root/go/src/github.com/cosmos/launch/launch  ~/go/src/github.com/cosmos/launch/systemctl/binary/

    gzip -f ~/go/src/github.com/cosmos/launch/systemctl/binary/okchaind
    gzip -f ~/go/src/github.com/cosmos/launch/systemctl/binary/okchaincli
    gzip -f ~/go/src/github.com/cosmos/launch/systemctl/binary/launch

#    cp -r ~/go/src/github.com/cosmos/launch/systemctl/binary ~/go/launch/systemctl/
}

function main {
    pull_build
#    getbinary
}

main