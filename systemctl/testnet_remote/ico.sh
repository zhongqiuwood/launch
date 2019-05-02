#!/usr/bin/env bash
#set -e

CURDIR=`dirname $0`
. ${CURDIR}/../scripts/okchaind.profile
. ./token.profile

OKCHAIN_CLI=okchaincli
BEGIN_PROPOSALID=1
ADMIN_HOME=~/.okchaincli/admin


while getopts "i:apvVq" opt; do
  case $opt in
    q)
      echo "QUERY_ONLY"
      QUERY_ONLY="Y"
      ;;
    v)
      echo "VOTE_ONLY"
      VOTE_ONLY="Y"
      ;;
    V)
      echo "VOTE_ONE"
      VOTE_ONE="Y"
      ;;
    p)
      echo "PROPOSAL_ONLY"
      PROPOSAL_ONLY="Y"
      ;;
    i)
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
       ${OKCHAIN_CLI} keys add --recover admin${i} -y --home ${ADMIN_HOME}${i} -m "${mnemonic}"
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

proposal_vote_only() {
    # $1: 1st proposal id
    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        token=${TOKENS[i]}
        okecho proposal ${token}
        sleep 1

        ((proposal_id = i + $1))
        okecho vote ${proposal_id}

        echo "------------------------------------"
    done
}


query_only() {
    # $1: 1st proposal id
    for ((i=0; i<${#TOKENS[@]}; i++))
    do
        token=${TOKENS[i]}
        okecho okchaincli query token info $token --node ${TESTNET_RPC_INTERFACE} --chain-id okchain
        echo "------------------------------------"
    done
    okchaincli query token tokenpair --node ${TESTNET_RPC_INTERFACE} --chain-id okchain|grep base_asset_symbol|wc -l

}

main() {

    okchaincli config chain-id okchain
    okchaincli config trust-node true
    okchaincli config output json
    okchaincli config indent true

    if [ ! -z "${QUERY_ONLY}" ]; then
        query_only
        exit
    fi

    if [ ! -z "${ACTIVE_ONLY}" ]; then
        active_all $1
        okchaincli query token tokenpair --node ${TESTNET_RPC_INTERFACE} |grep base_asset_symbol|wc -l
        exit
    fi

    recover

    if [ ! -z "${PROPOSAL_ONLY}" ]; then
        proposal_vote_only $1
        okchaincli query token tokenpair --node ${TESTNET_RPC_INTERFACE} |grep base_asset_symbol|wc -l
        exit
    fi

    ico $1

    okchaincli query token tokenpair --node ${TESTNET_RPC_INTERFACE} |grep base_asset_symbol|wc -l
#    okchaincli query token tokenpair  |grep base_asset_symbol|wc -l
}

main ${BEGIN_PROPOSALID}