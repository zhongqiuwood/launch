#!/bin/bash

. ./okchaind.profile

${OKCHAIN_CLI} tx gov submit-dex-list-proposal \
    --title="list $1/okb" \
    --description="list $1/okb" \
    --type=DexList \
    --deposit="100000okb" \
    --listAsset="$1" \
    --quoteAsset="okb" \
    --initPrice="2.25" \
    --maxPriceDigit=4 \
    --maxSizeDigit=4 \
    --minTradeSize="0.001" \
    --from captain -y \
    --home ${HOME_CLI}