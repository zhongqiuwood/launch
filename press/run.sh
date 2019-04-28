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





#for((;;)) do
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100
#
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100 -s
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100 -s
#    ./order.sh -q 1 -p 0.1 -c 20 -x 10 -b 10 -d 100 -s
#
#    ./reward.sh
#    ./mint.sh 1000
#done


