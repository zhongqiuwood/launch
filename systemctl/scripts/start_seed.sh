#!/bin/bash

IFS="!!"

. ${HOME}/okchain/launch/systemctl/scripts/okchaind.profile

if [ ! -d ${HOME_DAEMON} ]; then
    
    ${OKCHAIN_CLI} keys add --recover admin --home ${HOME_CLI}  -y -m "${ADMIN_MNEMONIC}"
    ${OKCHAIN_CLI} keys add --recover captain --home ${HOME_CLI} -y -m "${CAPTAIN_MNEMONIC}"

    ${OKCHAIN_DAEMON} init --chain-id okchain --home ${HOME_DAEMON}

    ${OKCHAIN_CLI} config chain-id okchain --home ${HOME_CLI}
    ${OKCHAIN_CLI} config trust-node true --home ${HOME_CLI}
    ${OKCHAIN_CLI} config indent true --home ${HOME_CLI}

    ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show admin -a --home ${HOME_CLI}) \
        1okb --home ${HOME_DAEMON}

    ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show captain -a --home ${HOME_CLI}) \
        1000000000okb --home ${HOME_DAEMON}

    ${OKCHAIN_DAEMON} gentx --amount 1okb --min-self-delegation 1 --commission-rate 0.1 \
        --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
        --pubkey $(${OKCHAIN_DAEMON} tendermint show-validator --home ${HOME_DAEMON}) \
        --name admin --home ${HOME_DAEMON} --home-client ${HOME_CLI}

    rm ${OKCHAIN_LAUNCH_TOP}/gentx/data/gentx-*
    cp ${HOME_DAEMON}/config/gentx/gentx-*.json ${OKCHAIN_LAUNCH_TOP}/gentx/data


    for (( i=0;i<${#OKCHAIN_TESTNET_FULL_HOSTS[@]};i++))
    do
        host=${HOSTS_PREFIX}${OKCHAIN_TESTNET_FULL_HOSTS[i]}
        mnemonic=${OKCHAIN_TESTNET_FULL_MNEMONIC[i]}
        home_d=${HOME_DAEMON}/${host}
        home_cli=${HOME_CLI}/${host}

        ${OKCHAIN_CLI} keys add --recover captain --home ${home_cli} -y -m "${CAPTAIN_MNEMONIC}"
        ${OKCHAIN_CLI} keys add --recover ${host} --home ${home_cli}  -y -m "${mnemonic}"

        ${OKCHAIN_DAEMON} init --chain-id okchain --home ${home_d}
        ${OKCHAIN_CLI} config chain-id okchain --home ${home_cli}
        ${OKCHAIN_CLI} config trust-node true --home ${home_cli}
        ${OKCHAIN_CLI} config indent true --home ${home_cli}

        ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show ${host} -a --home ${home_cli}) \
            1okb --home ${home_d}
        ${OKCHAIN_DAEMON} gentx --amount 1okb --min-self-delegation 1 --commission-rate 0.1 \
            --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
            --pubkey $(${OKCHAIN_DAEMON} tendermint show-validator --home ${home_d}) \
            --name ${host} --home ${home_d} --home-client ${home_cli}
        cp ${home_d}/config/gentx/gentx-*.json ${OKCHAIN_LAUNCH_TOP}/gentx/data
    done

    cd ${OKCHAIN_LAUNCH_TOP}/
    ${OKCHAIN_LAUNCH_TOP}/launch

    cp ${OKCHAIN_LAUNCH_TOP}/genesis.json ${HOME_DAEMON}/config
    for host in ${OKCHAIN_TESTNET_FULL_HOSTS[@]}
    do
        cp ${OKCHAIN_LAUNCH_TOP}/genesis.json ${HOME_DAEMON}/${HOSTS_PREFIX}${host}/config
    done
fi

if [ ${IP_INNET} = true ];then
    LOCAL_IP=`ifconfig  | grep ${IP_PREFIX} | awk '{print $2}' | cut -d: -f2`
else
    LOCAL_IP=`curl ifconfig.me`
fi

${OKCHAIN_DAEMON} start --home ${HOME_DAEMON} \
    --p2p.seed_mode=true \
    --p2p.addr_book_strict=false \
    --log_level main:info,state:info,x/order:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656