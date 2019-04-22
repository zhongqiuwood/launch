OKCHAIN_DAEMON=/home/ubuntu/okchain/launch/systemctl/cloud/binary/okchaind
OKCHAIN_CLI=/home/ubuntu/okchain/launch/systemctl/cloud/binary/okchaincli
HOME_DAEMON=/home/ubuntu/.okchaind
HOME_CLI=/home/ubuntu/.okchaincli

OKCHAIN_LAUNCH_TOP=/home/ubuntu/okchain/launch

# 是否使用内网IP 
# 如果false，则使用公网IP，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点公网IP
# 如果是true，则必须设置IP_PREFIX，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点内网IP
IP_INNET=true
IP_PREFIX=172.31
OKCHAIN_TESTNET_FULL_HOSTS=("172.31.29.237" "172.31.20.185" "172.31.22.135" "172.31.18.166" "172.31.28.204")
HOSTS_PREFIX=okchain_cloud

SEED_NODE_IP=172.31.26.8
SCP_TAG="-i ~/okchain-dex-test.pem ubuntu"

# OKCHAIN_DAEMON=/root/okchain/launch/systemctl/cloud/binary/okchaind
# OKCHAIN_CLI=/root/okchain/launch/systemctl/cloud/binary/okchaincli
# HOME_DAEMON=/root/.okchaind
# HOME_CLI=/root/.okchaincli

# OKCHAIN_LAUNCH_TOP=/root/okchain/launch

# IP_INNET=true
# IP_PREFIX=192.168
# OKCHAIN_TESTNET_FULL_HOSTS=("192.168.13.121" "192.168.13.122" "192.168.13.123" "192.168.13.124" "192.168.13.125")
# HOSTS_PREFIX=okchain

# SEED_NODE_IP=192.168.13.116
# SCP_TAG="root"