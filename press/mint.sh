#!/usr/bin/env bash


MAX=$1

okdexcli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CC}
okdexcli tx token mint -s btc -n ${MAX} --from captain --chain-id okchain -y ${CC}
okdexcli tx token mint -s eos -n ${MAX} --from captain --chain-id okchain -y ${CC}

addr=$(okchaincli keys show -a captain)

okdexcli query account ${addr} --chain-id okchain ${CC}
