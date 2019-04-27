#!/bin/bash

okecho() {
 echo "shell exec: [$@]"
 $@
}

okecho okchaincli query account $(okchaincli keys show alice -a)


#printf "\njack\n"
okecho okchaincli query account $(okchaincli keys show jack -a)

