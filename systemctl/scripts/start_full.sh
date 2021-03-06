#!/bin/bash

. ${HOME}/okchain/launch/systemctl/scripts/okchaind.profile

scp ${SCP_TAG}@${SEED_NODE_IP}:${OKCHAIN_LAUNCH_TOP}/systemctl/scripts/seednode.profile \
    ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts/

. ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts/seednode.profile

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

if [ ! -d ${HOME_DAEMON} ]; then
    host=${HOSTS_PREFIX}${LOCAL_IP}
    scp -r ${SCP_TAG}@${SEED_NODE_IP}:${HOME_DAEMON}/${host}/ ${HOME_DAEMON}/
    scp -r ${SCP_TAG}@${SEED_NODE_IP}:${HOME_CLI}/${host}/ ${HOME_CLI}/
fi

${OKCHAIN_DAEMON} start --home ${HOME_DAEMON} \
    --p2p.seeds ${SEED_NODE_ID}@${SEED_NODE_URL} \
    --p2p.addr_book_strict=false \
    --log_level main:info,state:info,x/order:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656