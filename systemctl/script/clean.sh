#!/bin/bash

. ./okchaind.profile

sudo rm -rf ${HOME_CLI}
sudo rm -rf ${HOME_DAEMON}
sudo rm -f ${HOME}/okchaind.log