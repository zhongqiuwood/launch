#!/bin/bash

source /root/go/src/github.com/cosmos/launch/systemctl/seednode/okdexd.profile

if [ ! -d /root/.okdexd ]; then

    /root/go/bin/okdexd init --chain-id okchain --home /root/.okdexd

    scp root@192.168.13.116:/root/go/src/github.com/cosmos/launch/genesis.json /root/.okdexd/config/
fi

MYIP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

# --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_ADD}
/root/go/bin/okdexd start --home /root/.okdexd --p2p.addr_book_strict=false --log_level *:info --p2p.laddr tcp://${MYIP}:26656  2>&1 > /root/okdexd.log &
