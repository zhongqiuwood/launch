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

function gitpull {
echo git pull@$1
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}
    git checkout .
    git pull
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/
    ./unzip.sh
    mv ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/launch ${OKCHAIN_LAUNCH_TOP}/
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./service.sh
eeooff
echo done!
}

function gitclone {
echo git clone@$1
${SSH}@$1 << eeooff
    sudo rm -rf ${OKCHAIN_LAUNCH_TOP}
    git clone ${LAUNCH_GIT} ${OKCHAIN_LAUNCH_TOP}
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/binary/
    ./unzip.sh
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

    for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
    do
        gitpull ${host} &
    done
}

main