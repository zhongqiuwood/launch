#!/bin/bash
source /root/go/src/github.com/cosmos/launch/systemctl/seednode/okdexd.profile

/root/go/bin/okdexd start --p2p.seed_mode=true --p2p.addr_book_strict=false --log_level *:info --p2p.laddr tcp://${LOCAL_IP}:26656 2>&1 > okdexd.log &