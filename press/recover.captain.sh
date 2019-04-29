#!/usr/bin/env bash


okchaincli keys add --recover admin   -y -m "keen border system oil inject hotel hood potato shed pumpkin legend actor"
okchaincli keys add --recover captain -y -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

okchaincli tx token mint -s okb -n 1 --from captain --chain-id okchain -y ${CCC}
okchaincli tx token mint -s btc -n 1 --from captain --chain-id okchain -y ${CCC}
