#!/bin/bash

. ./okchaind.profile

function main {
    seedid=$(${OKCHAIN_DAEMON} tendermint show-node-id --home ${HOME_DAEMON})
    cat>${OKCHAIN_LAUNCH_TOP}/systemctl/script/seednode.profile<<EOF
SEED_NODE_ID=${seedid}
SEED_NODE_IP=${SEED_NODE_IP}
SEED_NODE_URL=${SEED_NODE_IP}:26656
EOF
}

main