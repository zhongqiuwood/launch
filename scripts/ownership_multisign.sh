#!/bin/bash

# add from sig
okchaincli tx token chown --from $(okchaincli keys show alice -a) --to $(okchaincli keys show jack -a) --symbol okb > unsignedTx.json

# add to sign
okchaincli tx token multisigns unsignedTx.json --from jack > unsignedTx2.json

# owner sign
okchaincli tx sign --from alice unsignedTx2.json > signedTx.json

# broadcast tx
okchaincli tx broadcast signedTx.json
