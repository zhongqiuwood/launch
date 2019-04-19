
OKCHAIN_DAEMON=/usr/local/go/bin/okchaind
OKCHAIN_CLI=/usr/local/go/bin/okchaincli
HOME_DAEMON=/root/tmp/.okchaind
HOME_CLI=/root/tmp/.okchaincli

LAUNCH_TOP=/root/go/src/github.com/cosmos/launch

# 是否使用内网IP 如果false，则使用公网IP；如果是true，则必须设置IP_PREFIX
IP_INNET=true
IP_PREFIX=192.168
HOSTS_PREFIX=okchain

SEED_NODE_IP=192.168.13.116
