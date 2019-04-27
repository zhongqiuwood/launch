#!/usr/bin/env bash

git config --global user.name "zhongqiuwood"
git config --global user.email "zhongqiuwood@gmail.com"

cd ${GOPATH}/src/github.com/ok-chain/okchain
git stash
git checkout master
git pull
make install


cd ${GOPATH}/src/github.com/cosmos/launch

cd ${GOPATH}/src/github.com

mkdir -p cosmos
cd cosmos
git clone https://github.com/okblockchainlab/launch.git
cd launch
