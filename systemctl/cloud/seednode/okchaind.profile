
OKCHAIN_DAEMON=/usr/local/go/bin/okchaind
OKCHAIN_CLI=/usr/local/go/bin/okchaincli
HOME_DAEMON=/root/tmp/.okchaind
HOME_CLI=/root/tmp/.okchaincli

LAUNCH_TOP=/root/go/src/github.com/cosmos/launch

# 是否使用内网IP 
# 如果false，则使用公网IP，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点公网IP
# 如果是true，则必须设置IP_PREFIX，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点内网IP
IP_INNET=true
IP_PREFIX=192.168
OKCHAIN_TESTNET_FULL_HOSTS=("192.168.13.121" "192.168.13.122" "192.168.13.123" "192.168.13.124" "192.168.13.125")

HOSTS_PREFIX=okchain

