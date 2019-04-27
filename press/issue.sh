#!/bin/bash



echo "usage: ./issue  symbol amount from"

okchaincli tx token issue --mintable --symbol $1 -n $2 --from $3 -y

sleep 2
okchaincli query token info $1


okchain14q9yyysp5uphhtnlztxf3pl4huk4fk63crsf7d