#!/usr/bin/env bash

rm -rf ~/.okdexd

./killbyname.sh okdexd

okdexd init --chain-id okchain -o -v

scp root@192.168.13.116:/root/go/src/github.com/cosmos/launch/genesis.json .

cp genesis.json ~/.okdexd/config

okdexd start --p2p.seeds $1@192.168.13.116:26656 --p2p.addr_book_strict=false
