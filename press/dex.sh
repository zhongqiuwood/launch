#!/bin/bash


PRODUCT=btc_okb
OKDEXCLI_HOME=~/.okchaincli

while getopts "P:h:" opt; do
  case $opt in
    P)
      echo "PRODUCT=$OPTARG"
      PRODUCT=$OPTARG
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

ENV="-P ${PRODUCT} -h ./products/${PRODUCT}/${PRODUCT}"

round() {

    ./order.sh ${ENV}
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s

    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
    ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s
}


main() {
    ./reward.sh -N -h ./products/${PRODUCT}/${PRODUCT}

    for((;;)) do

        round

        ./reward.sh
    done
}

main


