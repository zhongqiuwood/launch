#!/bin/bash

. ./okchaind.profile

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

if [ ${LOCAL_IP} = "${SEED_NODE_IP}" ];then
<<<<<<< HEAD:systemctl/script/service.sh
    sudo cp ${OKCHAIN_LAUNCH_TOP}/systemctl/script/okchaind_seed.service /etc/systemd/system/okchaind.service
else
    sudo cp ${OKCHAIN_LAUNCH_TOP}/systemctl/script/okchaind_full.service /etc/systemd/system/okchaind.service
=======
    sudo cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/profile/okchaind_seed.service /etc/systemd/system/okchaind.service
else
    sudo cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/profile/okchaind_full.service /etc/systemd/system/okchaind.service
>>>>>>> parent of 28f16f9... upd:systemctl/cloud/profile/service.sh
fi
sudo systemctl daemon-reload
