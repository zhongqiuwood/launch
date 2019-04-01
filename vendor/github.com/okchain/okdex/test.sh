#!/bin/bash
# alice send okb to jack
okdexcli tx send $(okdexcli keys show jack -a) 100000000000000000000000okb --from=alice --yes <<EOF
asdfghjkl
EOF

# jack issue token xxb
okdexcli tx token issue --from jack --symbol xxb -n 100 --yes <<EOF
asdfghjkl
EOF

# place order
okdexcli tx order new xxb_okb SELL 10.0 1.0 --from jack --yes<<EOF
asdfghjkl
EOF

okdexcli tx order new xxb_okb SELL 10.1 1.0 --from jack --yes<<EOF
asdfghjkl
EOF

okdexcli tx order new xxb_okb BUY 9.9 2.0 --from alice --yes<<EOF
asdfghjkl
EOF

okdexcli tx order new xxb_okb BUY 9.8 10.0 --from alice --yes<<EOF
asdfghjkl
EOF

# query depth book
okdexcli query order depthbook xxb_okb

okdexcli tx order new xxb_okb BUY 10.0 1.1 --from alice --yes<<EOF
asdfghjkl
EOF

okdexcli query order depthbook xxb_okb
