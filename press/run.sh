#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile


${CURDIR}/stop.sh

sleep 1
for token in ${TOKENS[@]}
do
    ${CURDIR}/reward.sh -N -h ${CURDIR}/products/${token}_okb/${token}_okb
    ${CURDIR}/dex.sh -P ${token}_okb 2>&1 >${token}_okb.json &
done

