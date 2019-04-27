#!/bin/bash

PROFILE=cloud_okchaind.profile

while getopts "cp:" opt; do
  case $opt in
    c)
      echo "GIT_CLONE"
      GIT_CLONE="true"
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


function push_build {
ssh root@192.168.13.116 << eeooff
    source ~/env.sh
    cd /root/go/src/github.com/ok-chain/okchain
    git stash
    git pull
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

function pull_build {
echo git pull@$1
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}
    git checkout .
    git pull
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/okbins
    ../unzip.sh

    mv ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./service.sh
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

    push_build

    for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
    do
        pull_build ${host} &
    done
}

main