#!/bin/bash

okchaincli tx order new xxb_okb BUY 0.1 0.1 --chain-id okchain --trust-node -y --from user -c $1 -x $2
