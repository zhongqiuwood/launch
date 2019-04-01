#!/bin/bash

okdexcli tx send  $(okdexcli keys show alice -a) 1mycoin --generate-only --from jack > unsignedtx.json <<EOF
y
asdfghjkl
EOF
