#!/bin/bash


# default params

CONCURRENT_NUM=10
NUM_PER_THREAD=10
DEPTH=10
BATCH_NUM=10
PRODUCT=btc_okb
SIDE=BUY
USER=user
QUANTITY=0.1
PRICE=0.1

while getopts "c:x:b:d:P:p:su:q:" opt; do
  case $opt in
    q)
      echo "QUANTITY=$OPTARG"
      QUANTITY=$OPTARG
      ;;
    p)
      echo "PRICE=$OPTARG"
      PRICE=$OPTARG
      ;;
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
      echo "DEPTH=$OPTARG"
      DEPTH=$OPTARG
      ;;
    P)
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


okchaincli tx order new ${PRODUCT} ${SIDE} ${PRICE} ${QUANTITY} --from ${USER} \
    -y -c ${CONCURRENT_NUM} -x ${NUM_PER_THREAD} -b ${BATCH_NUM} -d ${DEPTH} ${CCC}




exit

okchaincli tx order new btc_okb SELL 0.1 1000000000 --from captain -y ${CCC}
