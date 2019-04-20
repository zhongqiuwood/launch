#!/bin/bash

. ./okchaind.profile

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

host=${HOSTS_PREFIX}${LOCAL_IP}

if [ ${host} = "${HOSTS_PREFIX}${OKCHAIN_TESTNET_SEED_HOST}" ];then
    ${OKCHAIN_CLI} tx gov vote $1 yes --from admin -y --home ${HOME_CLI}
else
    ${OKCHAIN_CLI} tx gov vote $1 yes --from ${host} -y --home ${HOME_CLI}
fi
