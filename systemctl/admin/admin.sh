#!/bin/bash

PROFILE=cluster.profile

while getopts "qrcstap:" opt; do
  case $opt in
    q)
      echo "QUERY"
      QUERY="true"
      ;;
    \?)
      echo "Invalid option: -$OPTARG"
      ;;
  esac
done

. ./${PROFILE}



update_bash_profile() {
   echo "$1 update_bash_profile"
${SSH}@$1 << eeooff
    echo "source /home/ubuntu/.env.sh" >>  /home/ubuntu/.bash_profile
eeooff
}


#update_bash_profile() {
#   echo "$1 update_bash_profile"
#${SSH}@$1 << eeooff
#/bin/cat>/home/ubuntu/.bash_profile<<EOF
#export PATH="ENV_PREFIXPATH:/opt/mssql-tools/bin:"
#source /home/ubuntu/.env.sh
#EOF
#/bin/sed -i "s/ENV_PREFIX/$/g" /home/ubuntu/.bash_profile
#
#eeooff
#}

function copy_env {
    scp ${SCP_TAG} env.sh ubuntu@${1}:/home/ubuntu/.env.sh
}


visit_all() {
    for host in ${OKCHAIN_TESTNET_ALL_HOSTS[@]}
    do
        echo "$1 ${host}"
        $1 ${host}
    done
}



function main {
   visit_all copy_env
}

main

