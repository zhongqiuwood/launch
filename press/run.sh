#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

for token in ${TOKENS[@]}
do
    ./reward.sh -N -h ./products/${token}_okb/${token}_okb
done


for token in ${TOKENS[@]}
do
   ./dex.sh -P ${token}_okb 2>&1 >${token}_okb.json &
   sleep 10
done