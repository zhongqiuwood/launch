#!/bin/bash

OKCHAIN_DAEMON=/usr/local/go/bin/okchaind
OKCHAIN_CLI=/usr/local/go/bin/okchaincli
OKCHAIN_TESTNET_FULL_HOSTS=("okchain21" "okchain22" "okchain23" "okchain24" "okchain25")
IFS="!!"
OKCHAIN_TESTNET_FULL_MNEMONIC=("shine left lumber budget elegant margin aunt truly prize snap shy claw"
"tiny sudden coyote idea name thought consider jump occur aerobic approve media"
"hole galaxy armed garlic casino tumble fitness six jungle success tissue jaguar"
"breeze real effort sail deputy spray life real injury universe praise common"
"action verb surge exercise order pause wait special account kid hard devote"
)

. /root/go/src/github.com/cosmos/launch/systemctl/seednode/okchaind.profile

LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

if [ ! -d /root/.okchaind ]; then
    
    ${OKCHAIN_CLI} keys add --recover admin --home /root/.okchaincli  -y -m "keen border system oil inject hotel hood potato shed pumpkin legend actor"
    ${OKCHAIN_CLI} keys add --recover captain --home /root/.okchaincli -y -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

    ${OKCHAIN_DAEMON} init --chain-id okchain --home /root/.okchaind

    ${OKCHAIN_CLI} config chain-id okchain --home /root/.okchaincli
    ${OKCHAIN_CLI} config trust-node true --home /root/.okchaincli
    ${OKCHAIN_CLI} config indent true --home /root/.okchaincli

    ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show admin -a --home /root/.okchaincli) \
        2000000okb --home /root/.okchaind

    ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show captain -a --home /root/.okchaincli) \
        1000000000okb --home /root/.okchaind

    ${OKCHAIN_DAEMON} gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 \
        --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
        --pubkey $(${OKCHAIN_DAEMON} tendermint show-validator --home /root/.okchaind) \
        --name admin --home /root/.okchaind --home-client /root/.okchaincli

    rm /root/go/src/github.com/cosmos/launch/gentx/data/gentx-*
    cp /root/.okchaind/config/gentx/gentx-*.json /root/go/src/github.com/cosmos/launch/gentx/data


    for(( i=0;i<${#OKCHAIN_TESTNET_FULL_HOSTS[@]};i++))
    do
        host=${OKCHAIN_TESTNET_FULL_HOSTS[i]}
        mnemonic=${OKCHAIN_TESTNET_FULL_MNEMONIC[i]}
        home_d=/root/.okchaind/${host}
        home_cli=/root/.okchaincli/${host}

        ${OKCHAIN_CLI} keys add --recover ${host} --home ${home_cli}  -y -m "${mnemonic}"

        ${OKCHAIN_DAEMON} init --chain-id okchain --home ${home_d}
        ${OKCHAIN_CLI} config chain-id okchain --home ${home_cli}
        ${OKCHAIN_CLI} config trust-node true --home ${home_cli}
        ${OKCHAIN_CLI} config indent true --home ${home_cli}

        ${OKCHAIN_DAEMON} add-genesis-account $(${OKCHAIN_CLI} keys show ${host} -a --home ${home_cli}) \
            2000000okb --home ${home_d}
        ${OKCHAIN_DAEMON} gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 \
            --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
            --pubkey $(${OKCHAIN_DAEMON} tendermint show-validator --home ${home_d}) \
            --name ${host} --home ${home_d} --home-client ${home_cli}
        cp ${home_d}/config/gentx/gentx-*.json /root/go/src/github.com/cosmos/launch/gentx/data
    done

    cd /root/go/src/github.com/cosmos/launch/
    /root/go/src/github.com/cosmos/launch/launch

    cp /root/go/src/github.com/cosmos/launch/genesis.json /root/.okchaind/config
    for host in ${OKCHAIN_TESTNET_FULL_HOSTS[@]}
    do
        cp /root/go/src/github.com/cosmos/launch/genesis.json /root/.okchaind/${host}/config
    done
fi

${OKCHAIN_DAEMON} start --home /root/.okchaind \
    --p2p.seed_mode=true \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656 2>&1 >> /root/okchaind.log &