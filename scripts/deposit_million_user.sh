#!/bin/bash
set -e
CURDIR=`dirname $0`

# default params
BALANCE=6
HOUSE_NUM=2
ROOM_NUM=2
USER_NUM=16
RPC_NODE=localhost

while getopts "c:b:r:h:n:" opt; do
  case $opt in
    n)
      echo "RPC_NODE=$OPTARG"
      RPC_NODE=$OPTARG
      ;;
    b)
      echo "BALANCE=$OPTARG"
      BALANCE=$OPTARG
      ;;
    h)
      echo "HOUSE_NUM=$OPTARG"
      HOUSE_NUM=$OPTARG
      ;;
    c)
      echo "USER_NUM=$OPTARG"
      USER_NUM=$OPTARG
      ;;
    r)
      echo "ROOM_NUM=$OPTARG"
      ROOM_NUM=$OPTARG
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

for ((house_index=0; house_index<${HOUSE_NUM}; house_index++)) do
    for ((room_index=0;room_index<${ROOM_NUM}; room_index++)) do
        echo "./reward.sh -c ${USER_NUM} -h city/house${house_index}/room${room_index}/user -b ${BALANCE} -n ${RPC_NODE}"
        ./reward.sh -c ${USER_NUM} -h city/house${house_index}/room${room_index}/user -b ${BALANCE} -n ${RPC_NODE}
    done
done



