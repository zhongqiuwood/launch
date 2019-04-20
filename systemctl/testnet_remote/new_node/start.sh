#!/bin/bash

. ./okchaind.profile

if [ ! -d ${HOME_DAEMON} ]; then
    ${OKCHAIN_DAEMON} init --chain-id okchain --home ${HOME_DAEMON}
else
    ${OKCHAIN_DAEMON} unsafe-reset-all
fi

scp root@${SEED_NODE_IP}:${SEED_NODE_GENESIS} ${HOME_DAEMON}/config

${OKCHAIN_DAEMON} start --home ${HOME_DAEMON} \
    --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_URL} \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656  2>&1 >> okchaind.log &