#!/bin/bash


function main {
    gzip -dkf okchaind.gz
    gzip -dkf okchaincli.gz
    gzip -dkf launch.gz
}

main