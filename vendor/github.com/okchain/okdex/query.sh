#!/bin/bash

printf "alice\n"
okdexcli query account $(okdexcli keys show alice -a)


printf "\njack\n"
okdexcli query account $(okdexcli keys show jack -a)

