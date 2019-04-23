#!/bin/bash

. ./okchaind.profile

# rm -rf ${HOME_DAEMON}
# rm -rf ${HOME_CLI}
${OKCHAIN_DAEMON} unsafe-reset-all
rm -f ${HOME}/okchaind.log