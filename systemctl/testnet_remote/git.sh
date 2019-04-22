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
ssh -i "~/okchain-dex-test.pem" ubuntu@$1 << eeooff
    sudo su
    cd ${OKCHAIN_LAUNCH_TOP}
    git pull
<<<<<<< HEAD
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/script
    ./service.sh
=======
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/binary/launch ${OKCHAIN_LAUNCH_TOP}/
>>>>>>> parent of 28f16f9... upd
eeooff
echo done!
}

function gitclone {
echo git clone@$1
<<<<<<< HEAD
${SSH}@$1 << eeooff
    sudo rm -rf ${OKCHAIN_LAUNCH_TOP}
    git clone ${LAUNCH_GIT} ${OKCHAIN_LAUNCH_TOP}
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/script
    ./service.sh
=======
ssh -i "~/okchain-dex-test.pem" ubuntu@$1 << eeooff
    sudo su
    rm -rf ${OKCHAIN_LAUNCH_TOP}
    git clone https://github.com/okblockchainlab/launch.git ${OKCHAIN_LAUNCH_TOP}
    cp ${OKCHAIN_LAUNCH_TOP}/systemctl/cloud/binary/launch ${OKCHAIN_LAUNCH_TOP}/
>>>>>>> parent of 28f16f9... upd
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