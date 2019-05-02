#!/bin/bash

okecho() {
 echo "shell exec: [$@]"
 $@
}

token=$1
index=$2

okecho okchaincli query account $(okchaincli keys show user$index --home ./products/$token_okb/$token_okb$index -a) --node c22:26657

