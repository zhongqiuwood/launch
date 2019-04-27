#!/bin/bash

#set -e

RPC_NODE=localhost
RPC_PORT=26657
COIN_PREFIX=coin
COIN_NUM=5
OKDEXCLI_HOME=~/.okchaincli

while getopts "N:n:s:" opt; do
  case $opt in
    n)
      echo "RPC_NODE=$OPTARG"
      RPC_NODE=$OPTARG
      ;;
    N)
      echo "COIN_NUM=$OPTARG"
      COIN_NUM=$OPTARG
      ;;
    s)
      echo "COIN_PREFIX=$OPTARG"
      COIN_PREFIX=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done
RPC_INTERFACE=${RPC_NODE}:${RPC_PORT}


for ((index=0; index<${COIN_NUM}; index++)) do
    okchaincli tx token issue --mintable --symbol ${COIN_PREFIX}${index} \
        -n 990000000 --from captain -y --node ${RPC_INTERFACE}
    okchaincli query token info ${COIN_PREFIX}${index} --node ${RPC_INTERFACE}
done