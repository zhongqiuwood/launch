#!/bin/bash
rm -rf ~/.okdexd
rm -rf ~/.okdexcli

go install ./cmd/okdexd
go install ./cmd/okdexcli



okdexcli keys add boos<<EOF
asdfghjkl
EOF

okdexd init --chain-id okchain

okdexd add-genesis-account $(okdexcli keys show boos -a) 1000000000okb

okdexd gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okdexd tendermint show-validator) --name boos<<EOF
asdfghjkl
EOF


# okdexd init --chain-id testchain -v
# okdexcli keys add alice <<EOF
# asdfghjkl
# EOF
# okdexcli keys add jack <<EOF
# asdfghjkl
# EOF
# okdexd add-genesis-account $(okdexcli keys show alice -a) 1000mycoin,1000alicecoin
# okdexd add-genesis-account $(okdexcli keys show jack -a) 1000mycoin,1000jackcoin

# # config okbcli
# okdexcli config chain-id testchain
# okdexcli config trust-node true
# okdexcli config output json
# okdexcli config indent true

# # run
# okdexd start
