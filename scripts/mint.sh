#!/usr/bin/env bash


MAX=$1

okdexcli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${C16}
okdexcli tx token mint -s btc -n ${MAX} --from captain --chain-id okchain -y ${C16}
okdexcli tx token mint -s eos -n ${MAX} --from captain --chain-id okchain -y ${C16}

addr=$(okchaincli keys show -a captain ${C16})

okdexcli query account ${addr} --chain-id okchain ${C16}
