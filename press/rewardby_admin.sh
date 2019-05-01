#!/bin/bash
set -e
CURDIR=`dirname $0`
. ${CURDIR}/env.profile

. ${CURDIR}/../systemctl/testnet_remote/token.profile

# default params
USER_NUM=32
USER_NAME=user
OKDEXCLI_HOME=user_home/user
BALANCE=1000000
RPC_NODE=localhost
RPC_PORT=26657
ADMIN_INDEX=0

while getopts "c:h:u:a:b:n:Ni:" opt; do
  case $opt in
    N)
      echo "CREATE_NEW_USER"
      CREATE_NEW_USER="Y"
      ;;
    i)
      echo "ADMIN_INDEX=$OPTARG"
      ADMIN_INDEX=$OPTARG
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
ADMIN_HOME=admin${ADMIN_INDEX}_home


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

recover_admin() {
    if [ ! -f ${KEY_FILE} ]; then
        "${KEY_FILE} does not exist!"
        exit
    fi

    index=0
    cat ${KEY_FILE} | while read line
    do
        if [ $index -eq ${ADMIN_INDEX} ]; then
            okchaincli keys add admin${ADMIN_INDEX} --recover -m "${line}" -y --home ${ADMIN_HOME}
            okecho okchaincli keys show admin${ADMIN_INDEX} -a --home ${ADMIN_HOME}
            echo "$line"
            break
        fi
        ((index++))
    done
}

main() {

    init

    recover_admin

    if [ ! -z "${CREATE_NEW_USER}" ]; then
        okecho ${CURDIR}/genacc.sh -u ${USER_NAME} -c ${USER_NUM} -r -h ${OKDEXCLI_HOME}
    fi

    sleep 3

    header=$(okchaincli keys show ${USER_NAME}0 -a --home ${OKDEXCLI_HOME}0)
    okecho okchaincli tx send ${header} ${COINS} --from admin${ADMIN_INDEX} -y \
        --chain-id okchain --node ${RPC_INTERFACE} \
        --home ${ADMIN_HOME}

    sleep 5
    okecho okchaincli query account ${header} --chain-id okchain --node ${RPC_INTERFACE}
    okecho okchaincli tx send dummy ${REWARD} --from ${USER_NAME} -r ${USER_NUM} -y --chain-id okchain \
        --home ${OKDEXCLI_HOME} --node ${RPC_INTERFACE}
}


main



