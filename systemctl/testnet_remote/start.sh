#!/bin/bash

PROFILE=cloud_okchaind.profile
TOKENS=(btc eth eos ltc xrp)

while getopts "qrcstap:" opt; do
  case $opt in
    q)
      echo "QUERY"
      QUERY="true"
      ;;
    r)
      echo "RESTART"
      RESTART="true"
      ;;
    c)
      echo "CLEAN"
      CLEAN="true"
      ;;
    p)
      echo "PROFILE=$OPTARG"
      PROFILE=$OPTARG
      ;;
    t)
      echo "TOKEN"
      TOKEN="true"
      ;;
    s)
      echo "STOP"
      STOP="true"
      ;;
    a)
      echo "ACTIVE"
      ACTIVE="true"
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done

. ./${PROFILE}

start_seed_node() {
    echo start_seed_node@$1
${SSH}@$1 << eeooff
    sudo systemctl stop okchaind
    sudo systemctl start okchaind
    sudo systemctl status okchaind
    
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./setseed.sh
    
    exit
eeooff
}
 
start_full_node() {
    echo start_full_node@$1
${SSH}@$1 << eeooff
    sudo systemctl stop okchaind
    sudo systemctl start okchaind
    sudo systemctl status okchaind
    exit
eeooff
}

query_node() {
    echo start_full_node@$1
${SSH}@$1 << eeooff
    sudo systemctl status okchaind
    exit
eeooff
}

function clean {
    echo clean@$1
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./clean.sh
    exit
eeooff
}

function stop {
    echo stop@$1
${SSH}@$1 << eeooff
    sudo systemctl stop okchaind
    sudo systemctl status okchaind
    exit
eeooff
}

function stop_and_cleanup {
    echo stop@$1
${SSH}@$1 << eeooff
    sudo systemctl stop okchaind
    sudo systemctl status okchaind
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./clean.sh
    exit
eeooff
}

exe_stop() {
        for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
        do
            stop ${host}
            if [ ! -z "${CLEAN}" ];then
                clean ${host}
            fi
        done
        exit
}

exe_query() {
        for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
        do
            query_node ${host}
        done
        exit
}

run() {

    echo "========== start seed node =========="
    for host in ${OKCHAIN_TESTNET_SEED_HOST[@]}
    do
        if [ ! -z "${RESTART}" ];then
            stop ${host}
        fi

        if [ ! -z "${CLEAN}" ];then
            clean ${host}
        fi

        start_seed_node ${host}
    done

    echo "========== wating seed node done =========="
    sleep 30

    echo "========== start full node =========="
    for host in ${OKCHAIN_TESTNET_FULL_HOSTS[@]}
    do

        if [ ! -z "${RESTART}" ];then
            stop ${host}
        fi

        if [ ! -z "${CLEAN}" ];then
            clean ${host}
        fi

        start_full_node ${host} &
    done
}


function main {

    if [ ! -z "${TOKEN}" ];then
        ../scripts/ico.sh
        exit
    fi

    if [ ! -z "${STOP}" ];then
        exe_stop
    fi

    if [ ! -z "${QUERY}" ];then
        exe_query
    fi

    run
    exe_query
}

main