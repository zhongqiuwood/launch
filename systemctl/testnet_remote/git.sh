#!/bin/bash

. ./okchaind.profile

while getopts "c" opt; do
  case $opt in
    c)
      echo "GIT_CLONE"
      GIT_CLONE="true"
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done


function gitpull {
echo git pull@$1
ssh root@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}
    git pull
eeooff
echo done!
}

function gitclone {
echo git clone@$1
ssh root@$1 << eeooff
    rm -rf ${OKCHAIN_LAUNCH_TOP}
    git clone git@gitlab.okcoin-inc.com:dex/launch.git ${OKCHAIN_LAUNCH_TOP}
eeooff
echo done!
}

function main {
    if [ ! -z "${GIT_CLONE}" ];then
        for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
        do
            gitclone ${host}
        done

        exit
    fi

    for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
    do
        gitpull ${host}
    done
}

main