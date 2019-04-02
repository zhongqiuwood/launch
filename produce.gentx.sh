#!/usr/bin/env bash
rm -rf ~/.okdexd
okdexd init --chain-id okchain -o
okdexd add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 100okb
okdexd gentx --amount 1okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okdexd tendermint show-validator) --name admin
rm gentx/data/gentx-*
cp ~/.okdexd/config/gentx/gentx-*.json gentx/data