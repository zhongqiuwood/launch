#!/bin/bash

okecho() {
 echo "shell exec: [$@]"
 $@
}

token=$1
index=$2

okecho okchaincli keys show user$index --home ./products/${token}_okb/${token}_okb${index} -a
addr=$(okchaincli keys show user$index --home ./products/${token}_okb/${token}_okb${index} -a)
okecho okchaincli query account $addr --node c22:26657
