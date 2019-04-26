#!/bin/bash

okchaincli tx token issue --from alice --symbol bnb -n 200000
okchaincli query token info okb
okchaincli query account $(okchaincli keys show alice -a)
okchaincli tx token burn --from alice --symbol okb --amount 100
okchaincli tx token freeze --from alice --symbol okb --amount 0.1
okchaincli tx token unfreeze --from alice --symbol okb --amount 0.1
okchaincli tx send $(okchaincli keys show jack -a) 0.01okb --from alice

okchaincli tx token chown --from $(okchaincli keys show alice -a) --to $(okchaincli keys show jack -a) --symbol okb

okchaincli tx token multi-send --from alice --transfers '[{"to":"okchain192gpmr2dcvjjfk487jztdhskp5lpashq4qqdtt","amount":"1:okb,2:xmr"}]'