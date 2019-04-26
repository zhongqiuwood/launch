#!/bin/bash

okchaincli tx order new xxb_okb BUY 9.8 10.0 --generate-only --from alice > unsignedTx.json

okchaincli tx sign unsignedTx.json --from alice > signedTx.json