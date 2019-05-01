#!/bin/bash

# default params
USER_NUM=2
USER_NAME=user
OKDEXCLI_HOME=~/.okchaincli
REMOVE=n
KEY_FILE=admin_pkey.json
ADDR_FILE=admin_addr.json

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

if [ -f ${KEY_FILE} ]; then
    rm ${KEY_FILE}
fi

if [ -f ${ADDR_FILE} ]; then
    rm ${ADDR_FILE}
fi

for ((index=0; index < ${USER_NUM}; index++)) do
    okchaincli keys mnemonic >> ${KEY_FILE}
done

index=0
cat ${KEY_FILE} | while read line
do
   echo $line
   okchaincli keys add user${index} --recover -m "$line" -y
   okchaincli keys show user${index} -a >> ${ADDR_FILE}
   ((index++))
   echo "------------------------------"
done