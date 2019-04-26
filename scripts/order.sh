#!/bin/bash


# default params

CONCURRENT_NUM=10
NUM_PER_THREAD=10
DEPTH=10
BATCH_NUM=10
PRODUCT=btc_okb
SIDE=BUY
USER=user


while getopts "c:x:b:d:p:su:" opt; do
  case $opt in
    c)
      echo "CONCURRENT_NUM=$OPTARG"
      CONCURRENT_NUM=$OPTARG
      ;;
    x)
      echo "NUM_PER_THREAD=$OPTARG"
      NUM_PER_THREAD=$OPTARG
      ;;
    b)
      echo "BATCH_NUM=$OPTARG"
      BATCH_NUM=$OPTARG
      ;;
    d)
      echo "BALANCE=$OPTARG"
      BALANCE=$OPTARG
      ;;
    p)
      echo "PRODUCT=$OPTARG"
      PRODUCT=$OPTARG
      ;;
    s)
      echo "SIDE=SELL"
      SIDE=SELL
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done


okchaincli tx order new ${PRODUCT} ${SIDE} 0.1 0.1 --from ${USER} \
    -y -c ${CONCURRENT_NUM} -x ${NUM_PER_THREAD} -b ${BATCH_NUM} -d ${DEPTH} ${CC}

