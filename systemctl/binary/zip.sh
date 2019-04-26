#!/bin/bash


function main {
    echo "tar okchain binaries..."
    tar -zcvf okchainbins.tar.gz okchaind okchaincli launch
}

main