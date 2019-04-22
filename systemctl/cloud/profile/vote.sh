#!/bin/bash

. ./okchaind.profile

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

if [ ${LOCAL_IP} = "${SEED_NODE_IP}" ];then
    ${OKCHAIN_CLI} tx gov vote $1 yes --from admin --home ${HOME_CLI} -y
else
    ${OKCHAIN_CLI} tx gov vote $1 yes --from ${HOSTS_PREFIX}${LOCAL_IP} --home ${HOME_CLI} -y
fi
