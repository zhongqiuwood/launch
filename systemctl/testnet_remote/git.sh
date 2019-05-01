#!/bin/bash

PROFILE=cloud_okchaind.profile

while getopts "bcp:" opt; do
  case $opt in
    c)
      echo "GIT_CLONE"
      GIT_CLONE="true"
      ;;
    b)
      echo "REBUILD_BINARIES"
      REBUILD_BINARIES="true"
      ;;
    p)
      echo "PROFILE=$OPTARG"
      PROFILE=$OPTARG
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done

. ./${PROFILE}



function gitclone {
echo git clone@$1
${SSH}@$1 << eeooff
    sudo rm -rf ${OKCHAIN_LAUNCH_TOP}
    git clone ${LAUNCH_GIT} ${OKCHAIN_LAUNCH_TOP}
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/

    git clone ${OKBINS_GIT}
    cd okbins
    ../unzip.sh

    mv ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./service.sh
eeooff
echo done!
}


function rebuild_and_push {
ssh root@192.168.13.116 << eeooff
    source /root/env.sh
    cd /root/go/src/github.com/ok-chain/okchain
    git stash
    git pull
    git checkout dev
    make install

    cd /root/go/src/github.com/cosmos/launch
    git stash
    git pull
    go build
    cd /root/go/src/github.com/cosmos/launch/systemctl/binary/okbins
    cp /usr/local/go/bin/okchaind .
    cp /usr/local/go/bin/okchaincli .
    cp /root/go/src/github.com/cosmos/launch/launch .
    ../zip.sh
    ../gitpush.sh
eeooff
echo done!
}

function pull_update {
echo git pull@$1
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}
    git checkout .
    git pull
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/okbins
    git checkout .
    git pull
    ../unzip.sh

    mv ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./service.sh
eeooff
echo done!
}

function main {
    if [ ! -z "${GIT_CLONE}" ];then
        for host in ${OKCHAIN_TESTNET_DEPLOYED_HOSTS[@]}
        do
            gitclone ${host}
        done

        exit
    fi

    if [ ! -z "${REBUILD_BINARIES}" ];then
        rebuild_and_push
    fi

    for host in ${OKCHAIN_TESTNET_DEPLOYED_HOSTS[@]}
    do
        pull_update ${host} &
    done
}

main