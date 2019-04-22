#!/bin/bash

. ./okchaind.profile

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

if [ ${LOCAL_IP} = "${SEED_NODE_IP}" ];then
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/profile/okchaind_seed.service /etc/systemd/system/okchaind.service
else
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/profile/okchaind_full.service /etc/systemd/system/okchaind.service
fi
systemctl daemon-reload
