#!/bin/bash


function main {
    tar -zcvf okchainbins.tar.gz okchaind okchaincli launch
}

main