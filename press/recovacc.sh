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

if [ "${REMOVE}" == "Y" ]; then
    rm -rf ${OKDEXCLI_HOME}*
    okchaincli config chain-id okchain
    okchaincli config trust-node true
fi

index=0
cat ${KEY_FILE} | while read line
do
   echo $line
   okchaincli keys add ${USER_NAME}${index} --home ${OKDEXCLI_HOME}${index} --recover -m "$line" -y &
   ((index++))

   if [ $index -eq ${USER_NUM} ]; then
        break
   fi
done

