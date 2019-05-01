#!/bin/bash
set -e
CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile
. ./local.profile
# default params
USER_NUM=10
USER_NAME=user
OKDEXCLI_HOME=~/.okchaincli
BALANCE=8000000000
RPC_NODE=c22
RPC_PORT=26657
#ADMIN_NAME=captain
ADMIN_HOME=~/.okchaincli
while getopts "c:h:u:a:b:n:NA:H:" opt; do
  case $opt in
    H)
      echo "ADMIN_HOME=$OPTARG"
      ADMIN_HOME=$OPTARG
      ;;
    A)
      echo "ADMIN_NAME=$OPTARG"
      ADMIN_NAME=$OPTARG
      ;;
    N)
      echo "CREATE_NEW_USER"
      CREATE_NEW_USER="Y"
      ;;
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

init() {

    for token in ${TOKENS[@]}
    do
        COINS=${COINS}${AMOUNT}${token}","
        REWARD=${REWARD}${BALANCE}${token}","
    done

#    COINS=${COINS}${AMOUNT}okb
#    REWARD=${REWARD}${BALANCE}okb
    COINS=${COINS}10okb
    REWARD=${REWARD}10okb

    echo "[${COINS}]"
    echo "[${REWARD}]"
}

okecho() {
 echo "shell exec: [$@]"
 $@
}

transfer2admin(){

    okchaincli keys add --recover captain -y \
        -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

#    ./mint.sh ${AMOUNT}

    okchaincli keys add ${ADMIN_NAME} -y --home ${ADMIN_HOME}

    admin=$(okchaincli keys show ${ADMIN_NAME} -a --home ${ADMIN_HOME})

    okecho okchaincli tx send ${admin} ${COINS} \
        --from captain \
        -y --chain-id okchain --node ${RPC_INTERFACE}
}

main() {

    init

    if [ ! -z "${ADMIN_NAME}" ]; then
        transfer2admin
        exit
    fi

    if [ ! -z "${CREATE_NEW_USER}" ]; then
        okecho ${CURDIR}/genacc.sh -u ${USER_NAME} -c ${USER_NUM} -r -h ${OKDEXCLI_HOME}
    fi



    sleep 5

    header=$(okchaincli keys show ${USER_NAME}0 -a --home ${OKDEXCLI_HOME}0)
    okecho okchaincli tx send ${header} ${COINS} \
        --from ${ADMIN_NAME} \
        --home ${ADMIN_HOME} \
        -y --chain-id okchain --node ${RPC_INTERFACE}

    sleep 5
    okecho okchaincli query account ${header} --chain-id okchain --node ${RPC_INTERFACE}
    okecho okchaincli tx send dummy ${REWARD} --from ${USER_NAME} -r ${USER_NUM} -y --chain-id okchain \
        --home ${OKDEXCLI_HOME} --node ${RPC_INTERFACE}
}


main


exit


for ((index=0;index<${USER_NUM};index++)) do
    res=$(okchaincli query account $(okchaincli keys show ${USER_NAME}${index} -a --home ${OKDEXCLI_HOME}${index}) \
        --node ${RPC_INTERFACE}|grep Coins)
    echo "${USER_NAME}${index} ${res}"
done

res=$(okchaincli query account $(okchaincli keys show captain -a) --node ${RPC_INTERFACE} |grep Coins)
echo "captain ${res}"



