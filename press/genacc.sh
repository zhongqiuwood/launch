#!/bin/bash
set -e
CURDIR=`dirname $0`
. ${CURDIR}/env.profile


# default params
USER_NUM=2
USER_NAME=user
OKDEXCLI_HOME=~/.okchaincli
REMOVE=n

while getopts "c:h:u:r" opt; do
  case $opt in
    r)
      echo "REMOVE=Y"
      REMOVE=Y
      ;;
    u)
      echo "USER_NAME=$OPTARG"
      USER_NAME=$OPTARG
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


CUR_DIR=`pwd`

function newkey() {
    echo "okchaincli keys add ${1}${2} --home ${OKDEXCLI_HOME}${2}"
    okchaincli keys add ${1}${2} --home ${OKDEXCLI_HOME}${2} -y
}

function recover() {
    echo "okchaincli keys add ${1}${2} --home ${OKDEXCLI_HOME}${2}"
    okchaincli keys add ${1}${2} --home ${OKDEXCLI_HOME}${2} -y
}

if [ "${REMOVE}" == "Y" ]; then
    rm -rf ${OKDEXCLI_HOME}*
    okchaincli config chain-id okchain
    okchaincli config trust-node true
fi

flow_control=0
for ((index=0;index<${USER_NUM};index++,flow_control++)) do
    if [ ${flow_control} -eq 128 ]; then
        echo "flow_control: ${flow_control}, sleep 5s..."
        sleep 5
        flow_control=0
    fi
    newkey ${USER_NAME} ${index} &
done




