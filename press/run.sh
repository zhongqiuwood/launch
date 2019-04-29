#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

for token in ${TOKENS[@]}
do
   ./dex.sh -P ${token}_okb &
done