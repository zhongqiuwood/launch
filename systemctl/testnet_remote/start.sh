#!/bin/bash

PROFILE=cloud_okchaind.profile
TOKENS=(btc eth eos ltc xrp)

while getopts "rcstap:" opt; do
  case $opt in
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

function vote {
        echo vote@$1 proposal=$2
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./vote.sh $2
    exit
eeooff
}

function issue {
        echo issue@$1 token=$2
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./issue.sh $2
    exit
eeooff
}

function proposal {
        echo proposal@$1 token=$2
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./proposal.sh $2
    exit
eeooff
}


function active {
        echo proposal@$1 proposal=$2
${SSH}@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/scripts
    ./active.sh $2
    exit
eeooff
}

exe_active() {
        for (( i=1;i<=${#TOKENS[@]};i++))
        do
            for host in ${OKCHAIN_TESTNET_SEED_HOST[@]}
            do
                active ${host} ${i}
            done
        done
        exit
}

exe_ico() {
        for (( i=0;i<${#TOKENS[@]};i++))
        do
            token=${TOKENS[i]}

            for host in ${OKCHAIN_TESTNET_SEED_HOST[@]}
            do
                issue ${host} ${token}
                sleep 2
                proposal ${host} ${token}
            done
            sleep 2
            for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
            do
                ((id = ${i} + 1))
                vote ${host} ${id}
            done
        done
        exit
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
    if [ ! -z "${ACTIVE}" ];then
        exe_active
    fi

    if [ ! -z "${TOKEN}" ];then
       exe_ico
    fi

    if [ ! -z "${STOP}" ];then
        exe_stop
    fi

    run
}

main