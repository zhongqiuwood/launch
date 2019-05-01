#!/usr/bin/env bash
set -e

CURDIR=`dirname $0`
. ${CURDIR}/../scripts/okchaind.profile
. ./token.profile

OKCHAIN_CLI=okchaincli
BEGIN_PROPOSALID=1
ADMIN_HOME=~/.okchaincli/admin


while getopts "ap:" opt; do
  case $opt in
    p)
      echo "BEGIN_PROPOSALID=$OPTARG"
      BEGIN_PROPOSALID=$OPTARG
      ;;
    a)
      echo "ACTIVE_ONLY"
      ACTIVE_ONLY="true"
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done


okecho() {
    echo "shell exec: [$@]"
    $@
}

echoonly() {
    echo "shell exec: [$@]"
}

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
        --chain-id okchain \
        --node ${TESTNET_RPC_INTERFACE}
}

issue() {
    ${OKCHAIN_CLI} tx token issue --from captain --symbol ${1} --chain-id okchain \
        -n 89999999999 --mintable=true -y --node ${TESTNET_RPC_INTERFACE}
}

vote() {
   for ((idx=0; idx<${#OKCHAIN_TESTNET_ALL_ADMIN_MNEMONIC[@]}; idx++))
   do
       okecho ${OKCHAIN_CLI} tx gov vote $1 yes --from admin${idx} -y \
        --node ${TESTNET_RPC_INTERFACE} --home ${ADMIN_HOME}${idx} \
        --chain-id okchain &
   done
}

recover() {

   ${OKCHAIN_CLI} keys add --recover captain -y -m "${CAPTAIN_MNEMONIC}"
   for ((i=0; i<${#OKCHAIN_TESTNET_ALL_ADMIN_MNEMONIC[@]}; i++))
   do
       mnemonic=${OKCHAIN_TESTNET_ALL_ADMIN_MNEMONIC[i]}
       ${OKCHAIN_CLI} keys add --recover admin${i} -y --home ${ADMIN_HOME}${i} -m "${mnemonic}" &
   done
}

active() {
    okecho ${OKCHAIN_CLI} tx gov dexlist --proposal $1 --chain-id okchain --from captain -y --node ${TESTNET_RPC_INTERFACE}
}

active_all() {
    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        ((proposal_id = i + $1))
        okecho active ${proposal_id}
    done
}


ico() {
    # $1: 1st proposal id
    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        token=${TOKENS[i]}
        okecho issue ${token}
        sleep 2
        okecho proposal ${token}
        ((proposal_id = i + $1))
        sleep 2
        okecho vote ${proposal_id}

        echo "token_index[$i]: vote $token done"
        echo "------------------------------------"

    done

    echo "sleeping ..."
    sleep 60

    active_all $1
}

main() {

    if [ ! -z "${ACTIVE_ONLY}" ]; then
        active_all $1
        exit
    fi

    recover
    ico $1
}

main ${BEGIN_PROPOSALID}