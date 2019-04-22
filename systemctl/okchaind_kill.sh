#!/bin/bash

NAME=okchaind
MYNAME="okchaind_kill.sh"

ps -ef|grep "$NAME"|grep -v grep |grep -v $MYNAME |awk '{print "kill -9 "$2", "$8}'
ps -ef|grep "$NAME"|grep -v grep |grep -v $MYNAME |awk '{print "kill -9 "$2}' | sudo sh
echo "All <$NAME> killed!"


