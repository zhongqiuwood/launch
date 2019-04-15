#!/bin/bash

. /root/go/src/github.com/cosmos/launch/systemctl/seednode/okchaind.profile

LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

if [ ! -d /root/.okchaind ]; then
    /usr/local/go/bin/okchaincli keys add --recover admin --home /root/.okchaincli <<EOF
12345678
mouse symptom casino left tornado aerobic bargain property fog execute hover also
EOF

    /usr/local/go/bin/okchaind init --chain-id okchain --home /root/.okchaind

    # config okbcli
    /usr/local/go/bin/okchaincli config chain-id okchain --home /root/.okchaincli
    /usr/local/go/bin/okchaincli config trust-node true --home /root/.okchaincli
    /usr/local/go/bin/okchaincli config indent true --home /root/.okchaincli

    /usr/local/go/bin/okchaind add-genesis-account okchain1krmfwu676ac575z8qk7cklurpnjtsmtjke7fzs 2000000okb --home /root/.okchaind

    /usr/local/go/bin/okchaind gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(/usr/local/go/bin/okchaind tendermint show-validator --home /root/.okchaind) --name admin --home /root/.okchaind --home-client /root/.okchaincli 

    rm /root/go/src/github.com/cosmos/launch/gentx/data/gentx-*
    cp /root/.okchaind/config/gentx/gentx-*.json /root/go/src/github.com/cosmos/launch/gentx/data

    cd /root/go/src/github.com/cosmos/launch/
    /root/go/src/github.com/cosmos/launch/launch

    cp /root/go/src/github.com/cosmos/launch/genesis.json /root/.okchaind/config
fi

/usr/local/go/bin/okchaind start --home /root/.okchaind \
    --p2p.seed_mode=true \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656 2>&1 >> /root/okchaind.log &