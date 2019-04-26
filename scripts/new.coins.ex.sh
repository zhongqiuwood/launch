#!/bin/bash


asset=$1
quote_asset=$2
from=$3

okchaincli tx gov submit-dex-list-proposal --title="list ${asset}/${quote_asset}" \
    --description="list ${asset}/${quote_asset}" --type=DexList --deposit="5000000okb"  \
    --list-asset="${asset}" --quote-asset="${quote_asset}" --expire-time="9000000" \
    --init-price="10000000" --max-price-digit=6 --max-size-digit=6 \
    --min-trade-size="0.0001" \
    --from ${from} -y

exit

okchaincli tx gov submit-dex-list-proposal --title="list ${asset}/${quote_asset}" \
    --description="list ${asset}/${quote_asset}" --type=DexList --deposit="5000000okb"  \
    --list-asset="${asset}" --quote-asset="${quote_asset}" --expire-time="200000" \
    --init-price="10000000" --max-price-digit=6 --max-size-digit=6 \
    --merge-types="0.1, 1, 10" \
    --min-trade-size="0.0001" \
    --from ${from}