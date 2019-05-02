#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

#TOKENS=(dash)

okecho() {
    echo "------------------------------------------------"
    echo "shell exec: [$@]"
#    $@
}



main() {

    index=0
    for token in ${TOKENS[@]}
    do

         okecho okchaincli query order depthbook ${token}_okb --node c22:26657
         okchaincli query order depthbook ${token}_okb --node c22:26657 |jq
        ((index++))
    done
}

main
