#!/usr/bin/env bash

#rm -rf ~/.okdexd
#rm -rf ~/.okdexcli

./killbyname.sh okdexd

okdexd init --chain-id okchain

scp root@192.168.13.116:/root/go/src/github.com/cosmos/launch/genesis.json .

cp genesis.json ~/.okdexd/config

MYIP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

okdexd start --p2p.seeds $1@192.168.13.116:26656 --p2p.addr_book_strict=false --log_level *:info --p2p.laddr tcp://${MYIP}:26656  2>&1 > okdexd.log &

tail -f okdexd.log


