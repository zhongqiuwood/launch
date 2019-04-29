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
OKDEXCLI_HOME=~/.okchaincli
NODE=c22
while getopts "c:x:b:d:P:p:su:q:h:n:" opt; do
  case $opt in
    n)
      echo "NODE=$OPTARG"
      NODE=$OPTARG
      ;;
    h)
      echo "OKDEXCLI_HOME=$OPTARG"
      OKDEXCLI_HOME=$OPTARG
      ;;
    q)
      echo "QUANTITY=$OPTARG"
      QUANTITY=$OPTARG
      ;;
    u)
      echo "USER=$OPTARG"
      USER=$OPTARG
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

okecho() {
    echo "shell exec: [$@]"
    $@
}

okecho okchaincli tx order new ${PRODUCT} ${SIDE} ${PRICE} ${QUANTITY} \
    --from ${USER} \
    --home ${OKDEXCLI_HOME} \
    --chain-id okchain \
    -y -c ${CONCURRENT_NUM} \
    -x ${NUM_PER_THREAD} -b ${BATCH_NUM} -d ${DEPTH} --node ${NODE}:26657

