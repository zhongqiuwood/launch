#!/bin/bash

source /root/go/src/github.com/cosmos/launch/systemctl/fullnode/okdexd.profile

/root/go/bin/okdexd start --home /root/.okdexd --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_ADD} --p2p.addr_book_strict=false --log_level *:info --p2p.laddr tcp://${LOCAL_IP}:26656  2>&1 > /root/okdexd.log &