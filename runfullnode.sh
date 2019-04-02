#!/usr/bin/env bash

./killbyname.sh okdexd

okdexd init --chain-id okchain -o -v

cp genesis.json ~/.okdexd/config

okdexd start --p2p.seeds 8ac17d229f7ea0f00ea11a376e56d61d199d5476@192.168.13.116:26656 --p2p.addr_book_strict=false
