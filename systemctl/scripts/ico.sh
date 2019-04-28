#!/usr/bin/env bash

. ./okchaind.profile
set -e

TOKENS=(xmr6)
OKCHAIN_CLI=okchaincli
proposal() {

    ${OKCHAIN_CLI} tx gov submit-dex-list-proposal \
        --title="list $1/okb" \
        --description="" \
        --type=DexList \
        --deposit="100000okb" \
        --listAsset="$1" \
        --quoteAsset="okb" \
        --initPrice="2.25" \
        --maxPriceDigit=4 \
        --maxSizeDigit=4 \
        --minTradeSize="0.001" \
        --from captain -y \
        --node ${TESTNET_RPC_INTERFACE}
}

issue() {
    ${OKCHAIN_CLI} tx token issue --from captain --symbol ${1} \
        -n 60000000000 --mintable=true -y --node ${TESTNET_RPC_INTERFACE}
}


vote() {
   for ((i=0; i<${#OKCHAIN_TESTNET_FULL_MNEMONIC[@]}; i++))
   do
       mnemonic=${OKCHAIN_TESTNET_FULL_MNEMONIC[i]}
       ${OKCHAIN_CLI} tx gov vote $1 yes --from admin${i} -y --node ${TESTNET_RPC_INTERFACE}
   done
}

recover() {

   ${OKCHAIN_CLI} keys add --recover captain -y -m "${CAPTAIN_MNEMONIC}"

   for ((i=0; i<${#OKCHAIN_TESTNET_FULL_MNEMONIC[@]}; i++))
   do
       mnemonic=${OKCHAIN_TESTNET_FULL_MNEMONIC[i]}
       ${OKCHAIN_CLI} keys add --recover admin${i} -y -m "${mnemonic}"
   done
}


active() {
    ${OKCHAIN_CLI} tx gov dexlist --proposal $1 --from captain -y --node ${TESTNET_RPC_INTERFACE}
}

ico() {
    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        token=${TOKENS[i]}
        issue ${token}
        sleep 2
        proposal ${token}
        ((proposal_id = i + $1))
        sleep 2
        vote ${proposal_id}
    done

    sleep 100

    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        ((proposal_id = i + $1))
        active ${proposal_id}
    done
}

main() {
#    recover
    ico $1
}

main $1