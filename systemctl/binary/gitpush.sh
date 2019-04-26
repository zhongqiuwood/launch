#!/bin/bash


function main {
    git add okchainbins.tar.gz
    git commit -m "update okchainbins"
    git push
}

main