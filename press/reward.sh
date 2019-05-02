#!/bin/bash
#set -e
CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

# default params
USER_NUM=32
USER_NAME=user
OKDEXCLI_HOME=~/.okchaincli
BALANCE=1000000
RPC_NODE=c22
RPC_PORT=26657


while getopts "c:h:u:a:b:n:N:Rm" opt; do
  case $opt in
    m)
      echo "MINT_COINS=Y"
      MINT_COINS="Y"
      ;;
    N)
      echo "CREATE_NEW_USER"
      CREATE_NEW_USER="Y"
      ;;
    R)
      echo "RECOVER_USER"
      RECOVER_USER="Y"
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
echo "RPC_INTERFACE: $RPC_INTERFACE"


init() {

    for token in ${TOKENS[@]}
    do
        COINS=${COINS}${AMOUNT}${token}","
        REWARD=${REWARD}${BALANCE}${token}","
    done

    COINS=${COINS}${AMOUNT}okb
    REWARD=${REWARD}${BALANCE}okb

    echo "[${COINS}]"
    echo "[${REWARD}]"
}

okecho() {
 echo "shell exec: [$@]"
 $@
}

main() {

    init

    if [ ! -z "${CREATE_NEW_USER}" ]; then
        okecho ${CURDIR}/genacc.sh -u ${USER_NAME} -c ${USER_NUM} -r -h ${OKDEXCLI_HOME}
        sleep 5
    fi

    if [ ! -z "${RECOVER_USER}" ]; then
        okecho ${CURDIR}/recovacc.sh -u ${USER_NAME} -c ${USER_NUM} -r -h ${OKDEXCLI_HOME}
        sleep 5
    fi

    okchaincli keys add --recover captain -y \
        -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

#    sleep 5

    header=$(okchaincli keys show ${USER_NAME}0 -a --home ${OKDEXCLI_HOME}0)
    okecho okchaincli tx send ${header} ${COINS} --from captain -y --chain-id okchain --node ${RPC_INTERFACE}

    sleep 5
    okecho okchaincli query account ${header} --chain-id okchain --node ${RPC_INTERFACE}
    okecho okchaincli tx send dummy ${REWARD} --from ${USER_NAME} -r ${USER_NUM} -y --chain-id okchain \
        --home ${OKDEXCLI_HOME} --node ${RPC_INTERFACE}

    if [ ! -z "${MINT_COINS}" ]; then
        okecho ./mint.sh ${AMOUNT} ${RPC_INTERFACE}
    fi
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



