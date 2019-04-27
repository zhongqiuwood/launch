#!/usr/bin/env bash


MAX=$1

okdexcli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CCC}
okdexcli tx token mint -s btc -n ${MAX} --from captain --chain-id okchain -y ${CCC}
okdexcli tx token mint -s eos -n ${MAX} --from captain --chain-id okchain -y ${CCC}

addr=$(okchaincli keys show -a captain)

okdexcli query account ${addr} --chain-id okchain ${CCC}
