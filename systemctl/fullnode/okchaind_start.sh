#!/bin/bash

. /root/go/src/github.com/cosmos/launch/systemctl/fullnode/okchaind.profile

scp root@${SEED_NODE_IP}:/root/go/src/github.com/cosmos/launch/systemctl/seednode/okchaind.profile \
    /root/go/src/github.com/cosmos/launch/systemctl/fullnode/

. /root/go/src/github.com/cosmos/launch/systemctl/fullnode/okchaind.profile

LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

if [ ! -d /root/.okchaind ]; then
    host="okchain"${LOCAL_IP:0-2:2}
    scp -r root@${SEED_NODE_IP}:/root/.okchaind/${host}/ /root/.okchaind/
    scp -r root@${SEED_NODE_IP}:/root/.okchaincli/${host}/ /root/.okchaincli/
fi
/usr/local/go/bin/okchaind start --home /root/.okchaind \
    --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_URL} \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656  2>&1 >> /root/okchaind.log &