#!/usr/bin/env bash
if [ $# -lt 1 ]; then
    echo "missing params..."
    exit
fi

for i in $@; do
    HOME_PATH=${GOPATH}/src/github.com/ok-chain/okchain/dev/okdex_testnet/cache/node${i}/okchaincli
    ADDR=$(okdexcli keys show node${i} --bech val -a --home ${HOME_PATH})
    echo "okdexcli tx staking delegate ${ADDR} 1okb --from node${i} --home ${HOME_PATH} --chain-id okchain --node tcp://localhost:10${i}57"
    okdexcli tx staking delegate ${ADDR} 1okb --from node${i} --home ${HOME_PATH} --chain-id okchain --node tcp://localhost:10${i}57<<EOF
y
EOF
done
