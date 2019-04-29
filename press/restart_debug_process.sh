#!/bin/bash

ps aux | grep okchaind | grep -v grep | awk '{print $2}' | xargs kill -TERM
rm -rf ~/.okchaind
rm -rf ~/.okchaincli
rm -f /tmp/sqlite3/backend.db

statik_bin_path=`which statik`

if [ ! -z ${statik_bin_path} ]; then
    statik -src=doc/swagger-ui/ -dest=doc/ -f
fi

#make install
go build -gcflags="all=-N -l" ./cmd/okchaincli
go build -gcflags="all=-N -l" ./cmd/okchaind

./okchaind init --chain-id okchain -v --backend
./okchaincli keys add alice
./okchaincli keys add jack

./okchaincli keys add --recover captain -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"


COINS=900100000000000btc,\
900100000000000xmr,\
900100000000000eos,\
90000000000000000okb,\
90000000000000000xxb

./okchaind add-genesis-account $(okchaincli keys show alice -a) ${COINS}
./okchaind add-genesis-account $(okchaincli keys show jack -a) ${COINS}
./okchaind add-genesis-account $(okchaincli keys show captain -a) ${COINS}

# config okbcli
./okchaincli config chain-id okchain
./okchaincli config trust-node true
./okchaincli config output json
./okchaincli config indent true
#
## run
nohup ./okchaind start --log_level *:debug 2>&1 > ./okdexd.log &
#okchaind start --log_level main:info,state:info,x/order:info

#tail -f ./okdexd.log | grep -i -E "(rpc)|(http)"
# nohup ./okchaincli rest-server --chain-id=okchain 2>&1 > ./okchaincli.log &
