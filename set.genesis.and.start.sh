#!/usr/bin/env bash


cp genesis.json ~/.okdexd/config

okdexd start --p2p.seed_mode=true --p2p.addr_book_strict=false --log_level *:debug
