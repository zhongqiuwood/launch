#!/usr/bin/env bash

./killbyname.sh okdexd

rm -rf ~/.okdexd
rm -rf ~/.okdexcli

./recover.admin.sh

okdexd init --chain-id okchain

# config okbcli
okdexcli config chain-id okchain
okdexcli config trust-node true
okdexcli config indent true

okdexd add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 2000000okb
okdexd gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okdexd tendermint show-validator) --name admin
rm gentx/data/gentx-*
cp ~/.okdexd/config/gentx/gentx-*.json gentx/data


go run main.go

cp ~/.okdexd/config/genesis.json genesis.tmp.json

okdexd tendermint show-node-id