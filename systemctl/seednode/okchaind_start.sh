#!/bin/bash

OKCHAIN_DAEMON=/usr/local/go/bin/okchaind
OKCHAIN_CLI=/usr/local/go/bin/okchaincli

. /root/go/src/github.com/cosmos/launch/systemctl/seednode/okchaind.profile

LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

if [ ! -d /root/.okchaind ]; then
    
    ${OKCHAIN_CLI} keys add --recover admin   -y -m "keen border system oil inject hotel hood potato shed pumpkin legend actor"
    ${OKCHAIN_CLI} keys add --recover captain -y -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

    ${OKCHAIN_DAEMON} init --chain-id okchain --home /root/.okchaind

    # config okbcli
    ${OKCHAIN_CLI} config chain-id okchain --home /root/.okchaincli
    ${OKCHAIN_CLI} config trust-node true --home /root/.okchaincli
    ${OKCHAIN_CLI} config indent true --home /root/.okchaincli

    ${OKCHAIN_DAEMON} add-genesis-account okchain1m3gmu4zlnv2hmqfu2jwr97r2653w9yshyvde07 \
        2000000okb    --home /root/.okchaind

    ${OKCHAIN_DAEMON} add-genesis-account okchain1kyh26rw89f8a4ym4p49g5z59mcj0xs4j045e39 \
        1000000000okb --home /root/.okchaind

    ${OKCHAIN_DAEMON} gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 \
        --commission-max-rate 0.5 --commission-max-change-rate 0.001 \
        --pubkey $(${OKCHAIN_DAEMON} tendermint show-validator --home /root/.okchaind) \
        --name admin --home /root/.okchaind --home-client /root/.okchaincli

    rm /root/go/src/github.com/cosmos/launch/gentx/data/gentx-*
    cp /root/.okchaind/config/gentx/gentx-*.json /root/go/src/github.com/cosmos/launch/gentx/data

    cd /root/go/src/github.com/cosmos/launch/
    go build
    /root/go/src/github.com/cosmos/launch/launch

    cp /root/go/src/github.com/cosmos/launch/genesis.json /root/.okchaind/config
fi

${OKCHAIN_DAEMON} start --home /root/.okchaind \
    --p2p.seed_mode=true \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656 2>&1 >> /root/okchaind.log &