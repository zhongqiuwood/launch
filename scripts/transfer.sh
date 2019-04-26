#!/bin/bash

okchaincli tx send $(okchaincli keys show alice -a) 0.000000002okb --from jack  -y

