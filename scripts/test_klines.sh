#!/bin/bash

# Place Sell order
okchaincli tx order new xxb_okb SELL 10.0 1.0 --from jack --yes
okchaincli tx order new xxb_okb SELL 10.1 1.0 --from jack --yes

# Place Buy Order
okchaincli tx order new xxb_okb BUY 9.9 2.0 --from alice --yes
okchaincli tx order new xxb_okb BUY 9.8 10.0 --from alice --yes
okchaincli tx order new xxb_okb BUY 10.0 1.1 --from alice --yes

sleep 60

# place SELL order
okchaincli tx order new xxb_okb SELL 10.0 1.0 --from jack --yes
okchaincli tx order new xxb_okb SELL 10.2 1.0 --from jack --yes

# place BUY order
okchaincli tx order new xxb_okb BUY 9.9 2.0 --from alice --yes
okchaincli tx order new xxb_okb BUY 9.8 10.0 --from alice --yes
okchaincli tx order new xxb_okb BUY 10.1 1.1 --from alice --yes


# Query Client Test
sleep 121
okchaincli backend tickers
okchaincli backend klines -g 60 -p xxb_okb -s 100
okchaincli backend klines -g 180 -p xxb_okb -s 100

# Restful API test
nohup okchaincli rest-server --chain-id=okchain 2>&1 > ./okdexcli.log &
sleep 0.5
curl http://localhost:1317/tickers; echo
curl "http://localhost:1317/tickers?count=100&sort=true"; echo
curl "http://localhost:1317/tickers?count=100&sort=false"; echo
curl "http://localhost:1317/tickers/xxb_okb?count=10"; echo

curl http://localhost:1317/candles; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=60"; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=60&size=1"; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=60&size=10"; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=60&size=1000"; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=60&size=1001"; echo
curl "http://localhost:1317/candles/xxb_okb?granularity=180"; echo

curl http://localhost:1317/deals; echo
curl http://localhost:1317/fees; echo

sleep 0.5
ps aux | grep okchaincli | grep -v grep | awk '{print $2}' | xargs kill -TERM
