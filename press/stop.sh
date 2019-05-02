#!/bin/bash
CURDIR=`dirname $0`

${CURDIR}/killbyname.sh dex.sh
sleep 1
${CURDIR}/killbyname.sh reward.sh
sleep 1
${CURDIR}/killbyname.sh order.sh
sleep 1
${CURDIR}/killbyname.sh okchaincli
sleep 1
rm -rf ${CURDIR}/products
rm -rf ${HOME}/.okchaincli
