#!/bin/bash

source /root/go/src/github.com/cosmos/launch/systemctl/fullnode/okdexd.profile


LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`


scp root@${SEED_NODE_IP}:/root/.okdexd/config/genesis.json ~/.okdexd/config

/root/go/bin/okdexd start --home /root/.okdexd \
    --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_URL} \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656  2>&1 >> /root/okdexd.log &