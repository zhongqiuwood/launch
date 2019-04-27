#!/bin/bash


for((;;)) do
    ./mint.sh 10000
    ./reward.sh
    ./order.sh -c 10 -x 10 -b 10 -d 50
    ./order.sh -c 10 -x 10 -b 10 -d 50 -s
    ./order.sh -c 10 -x 10 -b 10 -d 50
    ./order.sh -c 10 -x 10 -b 10 -d 50 -s
    ./order.sh -c 10 -x 10 -b 10 -d 50
    ./order.sh -c 10 -x 10 -b 10 -d 50 -s
done



