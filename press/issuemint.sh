#!/usr/bin/env bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

MAX=$1

mint(){
    for token in ${TOKENS[@]}
    do
        okchaincli tx token issue --mintable --symbol ${token} -n 89900000000 --from captain -y
    done

    okchaincli tx token issue --mintable --symbol okb -n 89900000000 --from captain -y

}

main(){
    mint
    for token in ${TOKENS[@]}
    do
        okchaincli query token info ${token}
    done

    okchaincli query token info okb


}

main