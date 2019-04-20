#!/bin/bash

. ./okchaind.profile

function main {
    seedid=$(${OKCHAIN_DAEMON} tendermint show-node-id --home ${HOME_DAEMON})
    cat>${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/seednode/seednode.profile<<EOF
SEED_NODE_ID=${seedid}
SEED_NODE_IP=${OKCHAIN_TESTNET_SEED_HOST}
SEED_NODE_URL=${OKCHAIN_TESTNET_SEED_HOST}:26656
EOF
}

main