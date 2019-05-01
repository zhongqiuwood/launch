#!/bin/bash


PRODUCT=btc_okb
OKDEXCLI_HOME=~/.okchaincli

while getopts "P:h:n" opt; do
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

okecho() {
    echo "shell exec: [$@]"
    $@
}

round() {

    okecho ./order.sh ${ENV} -n c22
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c21
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c22
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c23
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c24
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c25

    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c21
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c22
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c23
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c24
    okecho ./order.sh ${ENV} -c 10 -x 10 -b 10 -d 1 -p 0.1 -s -n c25
}


main() {


    for((;;)) do

        round

#        ./reward.sh
    done
}

main


