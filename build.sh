#!/usr/bin/env bash


cd ${GOPATH}/src/github.com/okchain/okdex
git stash
git checkout master
git pull
make install


cd ${GOPATH}/src/github.com
mkdir -p cosmos
cd cosmos
git clone https://github.com/okblockchainlab/launch.git
cd launch
./runfullnode.sh