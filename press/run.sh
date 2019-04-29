#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

./killbyname.sh dex.sh
sleep 1
./killbyname.sh reward.sh
sleep 1
./killbyname.sh order.sh
sleep 1
./killbyname.sh okchaincli
sleep 1
rm -rf ./products
rm -rf ~/.okchaincli

sleep 1
for token in ${TOKENS[@]}
do
    ./reward.sh -N -h ./products/${token}_okb/${token}_okb
    ./dex.sh -P ${token}_okb 2>&1 >${token}_okb.json &
done


#for token in ${TOKENS[@]}
#do
#   ./dex.sh -P ${token}_okb 2>&1 >${token}_okb.json &
#   sleep 10
#done
