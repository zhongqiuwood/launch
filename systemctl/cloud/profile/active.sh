#!/bin/bash

. ./okchaind.profile

${OKCHAIN_CLI} tx gov dexlist --proposal $1 --from captain -y --home ${HOME_CLI}