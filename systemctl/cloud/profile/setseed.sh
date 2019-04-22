#!/bin/bash

. ./okchaind.profile

function main {
    seedid=$(${OKCHAIN_DAEMON} tendermint show-node-id --home ${HOME_DAEMON})
<<<<<<< HEAD:systemctl/script/setseed.sh
    cat>${OKCHAIN_LAUNCH_TOP}/systemctl/script/seednode.profile<<EOF
=======
    cat>${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/profile/seednode.profile<<EOF
>>>>>>> parent of 28f16f9... upd:systemctl/cloud/profile/setseed.sh
SEED_NODE_ID=${seedid}
SEED_NODE_IP=${SEED_NODE_IP}
SEED_NODE_URL=${SEED_NODE_IP}:26656
EOF
}

main