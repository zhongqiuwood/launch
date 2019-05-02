#!/usr/bin/env bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

AMOUNT=$1
RPC_INTERFACE=$2

okecho() {
 echo "shell exec: [$@]"
 $@
}

mint(){
    for token in ${TOKENS[@]}
    do
        okecho okchaincli tx token mint -s ${token} -n ${AMOUNT} \
            --from captain --chain-id okchain -y --node ${RPC_INTERFACE}
    done
    okecho okchaincli tx token mint -s okb -n ${AMOUNT} \
        --from captain --chain-id okchain -y --node ${RPC_INTERFACE}
}

main(){
    mint
    addr=$(okchaincli keys show -a captain)
    okchaincli query account ${addr} --chain-id okchain --node ${RPC_INTERFACE}
}

main