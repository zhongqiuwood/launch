#!/bin/bash

okchaincli keys add a1
okchaincli keys add a2
okchaincli keys add a3

okchaincli keys add a0 --multisig a1,a2,a3 --multisig-threshold 2
okchaincli keys show a0 -m


okchaincli tx send $(okchaincli keys show a0 -a) 10okb --from jack -y

sleep 2
#查询 a0 余额
okchaincli query account $(okchaincli keys show a0 -a)

okchaincli tx send $(okchaincli keys show jack -a) 6okb --from a0 -y --generate-only > unsignedtx.json

okchaincli tx sign --multisig=$(okchaincli keys show a0 -a) --from=a1 --output-document=signed.by.a1.json unsignedtx.json
okchaincli tx sign --multisig=$(okchaincli keys show a0 -a) --from=a3 --output-document=signed.by.a3.json unsignedtx.json

okchaincli tx multisign unsignedtx.json a0 signed.by.a1.json signed.by.a3.json > signedtx.json

okchaincli tx sign --validate-signatures signedtx.json

#after broadcast, validate signatures will fail
exit
okchaincli tx broadcast signedtx.json
okchaincli query account $(okchaincli keys show jack -a)
okchaincli query account $(okchaincli keys show a0 -a)

