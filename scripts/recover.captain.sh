#!/usr/bin/env bash


okdexcli keys add --recover admin   -y -m "keen border system oil inject hotel hood potato shed pumpkin legend actor"
okdexcli keys add --recover captain -y -m "puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"

okdexcli tx token mint -s okb -n 1 --from captain --chain-id okchain -y $CC
okdexcli tx token mint -s btc -n 1 --from captain --chain-id okchain -y $CC
