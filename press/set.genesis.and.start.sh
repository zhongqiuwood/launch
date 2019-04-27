#!/usr/bin/env bash


cp genesis.json ~/.okdexd/config

okdexd start --p2p.seed_mode=true --p2p.addr_book_strict=false --log_level *:info --p2p.laddr tcp://192.168.13.116:26656
