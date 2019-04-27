#!/bin/bash
# jack issue token xxb
#okchaincli tx token issue --from jack --symbol xxb -n 100 --yes

# place order
okchaincli tx order new xxb_okb SELL 10.0 1.0 --from jack --yes

okchaincli tx order new xxb_okb SELL 10.1 1.0 --from jack --yes
#
okchaincli tx order new xxb_okb BUY 9.9 2.0 --from alice --yes

okchaincli tx order new xxb_okb BUY 9.8 10.0 --from alice --yes

# query depth book
okchaincli query order depthbook xxb_okb

okchaincli tx order new xxb_okb BUY 10.0 1.1 --from alice --yes
#
okchaincli query order depthbook xxb_okb
