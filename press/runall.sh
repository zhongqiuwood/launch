#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

./prerun.sh

./runtrade.sh 2>&1 > a.json &

