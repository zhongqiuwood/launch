#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

# 16 concurrent/admins, each 1000000 coins
./reward.sh -R -n c22 -b 10000000 -c 16 -u admin -h admin_home/admin
