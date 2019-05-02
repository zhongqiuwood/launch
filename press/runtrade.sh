#!/bin/bash

CURDIR=`dirname $0`

. ${CURDIR}/../systemctl/testnet_remote/token.profile

#TOKENS=(dash)

okecho() {
    echo "shell exec: [$@]"
    $@
}

trade() {

    admin_index=$1
    token=$2
    # 1000000 per user
    okecho ${CURDIR}/rewardby_admin.sh -N -i ${admin_index} -n c22 -c 16 -b 1000000 -h ${CURDIR}/products/${token}_okb/${token}_okb
    ${CURDIR}/dex.sh -P ${token}_okb 2>&1 >products/${token}_okb.json &
}


main() {

    ${CURDIR}/stop.sh

    sleep 1

    index=0
    for token in ${TOKENS[@]}
    do
        trade $index ${token}
        ((index++))
    done
}

main
