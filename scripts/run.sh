#!/bin/bash


ps aux | grep okchaind | grep -v grep | awk '{print $2}' | xargs kill -TERM

rm -rf ~/.okchaind*
rm -rf ~/.okchaincli*

rm -f /tmp/sqlite3/backend.db

statik_bin_path=`which statik`

if [ ! -z ${statik_bin_path} ]; then
    statik -src=doc/swagger-ui/ -dest=doc/ -f
fi

make install


okchaind init --chain-id okchain -v
okchaincli keys add alice

okchaincli keys add jack

# 1000*

okchaind add-genesis-account $(okchaincli keys show alice -a) \
    9001000000000000000000000btc,9001000000000000000000000xmr,9001000000000000000000000eos,900000000000000000000000000okb,900000000000000000000000000xxb
okchaind add-genesis-account $(okchaincli keys show jack -a) \
    9001000000000000000000000btc,9001000000000000000000000xmr,9001000000000000000000000eos,900000000000000000000000000okb,900000000000000000000000000xxb


HOME_DIR=~/.okchaincli

for ((n=0;n<16;n++)) do
   okchaincli keys add user${n} --home ${HOME_DIR}${n}
   okchaind add-genesis-account $(okchaincli keys show user${n} -a --home ${HOME_DIR}${n}) \
    1000000000000000000000btc,1000000000000000000000xmr,1000000000000000000000eos,100000000000000000000000okb,900000000000000000000000000xxb
done


# config okbcli
okchaincli config chain-id okchain
okchaincli config trust-node true
okchaincli config indent true

## run
nohup okchaind start 2>&1 > ./okdexd.log &

sleep 1
tail -f ./okdexd.log
