#!/bin/bash

./killbyname.sh dex.sh
sleep 1
./killbyname.sh reward.sh
sleep 1
./killbyname.sh order.sh
sleep 1
./killbyname.sh okchaincli
sleep 1
rm -rf ./products
rm -rf ~/.okchaincli
