#!/bin/bash

. ./okchaind.profile

TOKENS=(btc eth eos ltc xrp)

while getopts "csta" opt; do
  case $opt in
    c)
      echo "CLEAN"
      CLEAN="true"
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


function startseed {
    echo startseed@$1
ssh root@$1 << eeooff
    systemctl stop okchaind
    systemctl start okchaind
    systemctl status okchaind
    
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/testnet_remote
    ./setseed.sh
    
    exit
eeooff
}

function startfull {
    echo startfull@$1
ssh root@$1 << eeooff
    systemctl stop okchaind
    systemctl start okchaind
    systemctl status okchaind
    
    exit
eeooff
}

function clean {
    echo clean@$1
ssh root@$1 << eeooff
    rm -rf ${HOME_CLI}
    rm -rf ${HOME_DAEMON}
    rm -f /tmp/okchain/okchaind.log
    exit
eeooff
}

function stop {
    echo stop@$1
ssh root@$1 << eeooff
    systemctl stop okchaind
    systemctl status okchaind
    exit
eeooff
}

function vote {
        echo vote@$1 proposal=$2
ssh root@$1 << eeooff
    cd ${OKCHAIN_LAUNCH_TOP}/systemctl/testnet_remote
    ./vote.sh $2
    exit
eeooff
}

function issue {
        echo issue@$1 token=$2
ssh root@$1 << eeooff
    ${OKCHAIN_CLI} tx token issue --from captain --symbol $2 -n 10000000000 --mintable=true -y --passwd=12345678 --home ${HOME_CLI}
    exit
eeooff
}

function proposal {
        echo proposal@$1 token=$2
ssh root@$1 << eeooff
    ${OKCHAIN_CLI} tx gov submit-dex-list-proposal \
    --title="list $2/okb" \
    --description="list $2/okb" \
    --type=DexList \
    --deposit="100000okb" \
    --listAsset="$2" \
    --quoteAsset="okb" \
    --initPrice="2.25" \
    --maxPriceDigit=4 \
    --maxSizeDigit=4 \
    --minTradeSize="0.001" \
    --from captain -y \
    --home ${HOME_CLI}
    exit
eeooff
}


function active {
        echo proposal@$1 proposal=$2
ssh root@$1 << eeooff
    ${OKCHAIN_CLI} tx gov dexlist --proposal=$2 --from captain -y --home ${HOME_CLI}
    exit
eeooff
}

function main {
    if [ ! -z "${ACTIVE}" ];then
        for (( i=1;i<=${#TOKENS[@]};i++))
        do
            for host in ${OKCHAIN_TESTNET_SEED_HOST[@]}
            do
                active ${host} ${i}
            done
        done
        exit
    fi

    if [ ! -z "${TOKEN}" ];then
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
                let id=${i}+1
                vote ${host} ${id}
            done
        done
        exit
    fi

    if [ ! -z "${STOP}" ];then
        for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
        do
            stop ${host}
            if [ ! -z "${CLEAN}" ];then
                clean ${host}
            fi
        done
        exit
    fi

    echo "========== start seed node =========="
    for host in ${OKCHAIN_TESTNET_SEED_HOST[@]}
    do
        if [ ! -z "${CLEAN}" ];then
            stop ${host}
            clean ${host}
        fi

        startseed ${host}
    done
    
    echo "========== wating seed node done =========="
    sleep 30

    echo "========== start full node =========="
    for host in ${OKCHAIN_TESTNET_FULL_HOSTS[@]}
    do
        if [ ! -z "${CLEAN}" ];then
            clean ${host}
        fi

        startfull ${host}
    done
}

main