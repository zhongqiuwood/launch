#!/usr/bin/env bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

MAX=$1

okecho() {
 echo "shell exec: [$@]"
 $@
}

mint(){
    for token in ${TOKENS[@]}
    do
        okecho okchaincli tx token mint -s ${token} -n ${MAX} --from captain --chain-id okchain -y ${CC1}
    done
    okecho okchaincli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CC1}
#    okchaincli tx token mint -s okb -n ${MAX} --from captain --chain-id okchain -y ${CC1}
}

main(){
    mint
    addr=$(okchaincli keys show -a captain)
    okchaincli query account ${addr} --chain-id okchain ${CC1}
}

main