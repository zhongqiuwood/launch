#!/bin/bash


for((;;)) do
    ./mint.sh 10000
    ./reward.sh
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 0.1 -s
    ./order.sh -p 5 -q 10 -u captain -c 1 -x 1 -b 1 -d 1
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 10
    ./order.sh -p 5 -q 10 -u captain -c 1 -x 1 -b 1 -d 1
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 0.1 -s
    ./order.sh -p 5 -q 10 -u captain -c 1 -x 1 -b 1 -d 1
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 10
    ./order.sh -p 5 -q 10 -u captain -c 1 -x 1 -b 1 -d 1
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 0.1 -s
    ./order.sh -p 5 -q 10 -u captain -c 1 -x 1 -b 1 -d 1
    ./order.sh -c 10 -x 10 -b 10 -d 50 -p 10
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


