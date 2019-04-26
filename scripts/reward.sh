#!/bin/bash
set -e
CURDIR=`dirname $0`

# default params
USER_NUM=16
USER_NAME=user
OKDEXCLI_HOME=~/.okchaincli
AMOUNT=10000
BALANCE=6
RPC_NODE=localhost
RPC_PORT=26657

while getopts "c:h:u:a:b:n:" opt; do
  case $opt in
    n)
      echo "RPC_NODE=$OPTARG"
      RPC_NODE=$OPTARG
      ;;
    u)
      echo "USER_NAME=$OPTARG"
      USER_NAME=$OPTARG
      ;;
    a)
      echo "AMOUNT=$OPTARG"
      AMOUNT=$OPTARG
      ;;
    b)
      echo "BALANCE=$OPTARG"
      BALANCE=$OPTARG
      ;;
    c)
      echo "USER_NUM=$OPTARG"
      USER_NUM=$OPTARG
      ;;
    h)
      echo "OKDEXCLI_HOME=$OPTARG"
      OKDEXCLI_HOME=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done
((AMOUNT = USER_NUM * BALANCE))
RPC_INTERFACE=${RPC_NODE}:${RPC_PORT}



okecho() {
 echo "shell exec: [$@]"
 $@
}

okecho ${CURDIR}/genacc.sh -u ${USER_NAME} -c ${USER_NUM} -r -h ${OKDEXCLI_HOME}

okchaincli keys add --recover captain -y \
    -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

sleep 5

#okdexcli tx token mint -s btc -n ${AMOUNT} --from captain --chain-id okchain -y
#okdexcli tx token mint -s eos -n ${AMOUNT} --from captain --chain-id okchain -y
#okdexcli tx token mint -s okb -n ${AMOUNT} --from captain --chain-id okchain -y


COINS=${AMOUNT}okb

header=$(okchaincli keys show ${USER_NAME}0 -a --home ${OKDEXCLI_HOME}0)
okecho okchaincli tx send ${header} ${COINS} --from captain -y --chain-id okchain --node ${RPC_INTERFACE}

REWARD=${BALANCE}okb

sleep 5
okecho okchaincli query account ${header} --chain-id okchain --node ${RPC_INTERFACE}

okecho okchaincli tx send dummy ${REWARD} --from ${USER_NAME} -r ${USER_NUM} -y --chain-id okchain -p \
    --home ${OKDEXCLI_HOME} --node ${RPC_INTERFACE}

for ((index=0;index<${USER_NUM};index++)) do
    res=$(okchaincli query account $(okchaincli keys show ${USER_NAME}${index} -a --home ${OKDEXCLI_HOME}${index}) \
        --node ${RPC_INTERFACE}|grep Coins)
    echo "${USER_NAME}${index} ${res}"
done

res=$(okchaincli query account $(okchaincli keys show captain -a) --node ${RPC_INTERFACE} |grep Coins)
echo "captain ${res}"







