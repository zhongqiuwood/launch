#!/bin/bash

. ./okchaind.profile

${OKCHAIN_CLI} tx token issue --from captain --symbol $1 -n 60000000000 --mintable=true -y --passwd=12345678 --home ${HOME_CLI} --gas 300000