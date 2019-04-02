#!/usr/bin/env bash

okdexcli keys add $1<<EOF
12345678
EOF

okdexcli keys add --recover org <<EOF
12345678
puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer
EOF

okdexcli tx send $(okdexcli keys show $1 -a) 500okb --from=org --yes --chain-id okchain

sleep 1

okdexcli tx staking create-validator --amount 1000000okb --pubkey $(okdexd tendermint show-validator) --chain-id okchain --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --min-self-delegation 1 --from $1 --moniker $1<<EOF
y
EOF
