#!/usr/bin/env bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

MAX=$1

mint(){
    for token in ${TOKENS[@]}
    do
        okchaincli tx token mint -s ${token} -n ${MAX} --from captain --chain-id okchain -y ${CCC}
    done
    okchaincli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CCC}
    okchaincli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CCC}
}

main(){
    mint
    addr=$(okchaincli keys show -a captain)
    okchaincli query account ${addr} --chain-id okchain ${CCC}
}

main