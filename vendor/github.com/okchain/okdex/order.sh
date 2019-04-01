#!/bin/bash

okdexcli tx order new xxb_okb SELL 10.0 1.0 --from alice --yes<<EOF
asdfghjkl
EOF

okdexcli tx order new xxb_okb SELL 10.1 1.0 --from alice --yes<<EOF
asdfghjkl
EOF

okdexcli tx order new xxb_okb BUY 9.9 1.5 --from jack --yes<<EOF
asdfghjkl
EOF

