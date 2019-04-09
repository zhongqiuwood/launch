#!/bin/bash

. /root/go/src/github.com/cosmos/launch/systemctl/seednode/okdexd.profile

LOCAL_IP=`ifconfig  | grep 192.168 | awk '{print $2}' | cut -d: -f2`

if [ ! -d /root/.okdexd ]; then
    /usr/local/go/bin/okdexcli keys add --recover admin --home /root/.okdexcli <<EOF
12345678
keen border system oil inject hotel hood potato shed pumpkin legend actor
EOF

    /usr/local/go/bin/okdexd init --chain-id okchain --home /root/.okdexd

    # config okbcli
    /usr/local/go/bin/okdexcli config chain-id okchain --home /root/.okdexcli
    /usr/local/go/bin/okdexcli config trust-node true --home /root/.okdexcli
    /usr/local/go/bin/okdexcli config indent true --home /root/.okdexcli

    /usr/local/go/bin/okdexd add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 2000000okb --home /root/.okdexd

    /usr/local/go/bin/okdexd gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(/usr/local/go/bin/okdexd tendermint show-validator --home /root/.okdexd) --name admin --home /root/.okdexd --home-client /root/.okdexcli 

    rm /root/go/src/github.com/cosmos/launch/gentx/data/gentx-*
    cp /root/.okdexd/config/gentx/gentx-*.json /root/go/src/github.com/cosmos/launch/gentx/data

    cd /root/go/src/github.com/cosmos/launch/
    /root/go/src/github.com/cosmos/launch/launch

    cp /root/go/src/github.com/cosmos/launch/genesis.json /root/.okdexd/config
fi

/usr/local/go/bin/okdexd start --home /root/.okdexd \
    --p2p.seed_mode=true \
    --p2p.addr_book_strict=false \
    --log_level *:info \
    --p2p.laddr tcp://${LOCAL_IP}:26656 2>&1 >> /root/okdexd.log &