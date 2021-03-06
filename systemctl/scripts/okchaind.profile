OKCHAIN_DAEMON=/home/ubuntu/okchain/launch/systemctl/binary/okchaind
OKCHAIN_CLI=/home/ubuntu/okchain/launch/systemctl/binary/okchaincli
HOME_DAEMON=/home/ubuntu/.okchaind
HOME_CLI=/home/ubuntu/.okchaincli

CAPTAIN_MNEMONIC="puzzle glide follow cruel say burst deliver wild tragic galaxy lumber offer"
ADMIN_MNEMONIC="keen border system oil inject hotel hood potato shed pumpkin legend actor"

ADMIN1_MNEMONIC="shine left lumber budget elegant margin aunt truly prize snap shy claw"
ADMIN2_MNEMONIC="tiny sudden coyote idea name thought consider jump occur aerobic approve media"
ADMIN3_MNEMONIC="hole galaxy armed garlic casino tumble fitness six jungle success tissue jaguar"
ADMIN4_MNEMONIC="breeze real effort sail deputy spray life real injury universe praise common"
ADMIN5_MNEMONIC="action verb surge exercise order pause wait special account kid hard devote"

OKCHAIN_TESTNET_FULL_MNEMONIC=(
"${ADMIN1_MNEMONIC}"
"${ADMIN2_MNEMONIC}"
"${ADMIN3_MNEMONIC}"
"${ADMIN4_MNEMONIC}"
"${ADMIN5_MNEMONIC}"
)

OKCHAIN_TESTNET_ALL_ADMIN_MNEMONIC=(
"${ADMIN_MNEMONIC}"
"${ADMIN1_MNEMONIC}"
"${ADMIN2_MNEMONIC}"
"${ADMIN3_MNEMONIC}"
"${ADMIN4_MNEMONIC}"
"${ADMIN5_MNEMONIC}"
)

TESTNET_RPC_INTERFACE=okchain_cloud22:26657
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

# OKCHAIN_DAEMON=/root/okchain/launch/systemctl/binary/okchaind
# OKCHAIN_CLI=/root/okchain/launch/systemctl/binary/okchaincli
# HOME_DAEMON=/root/.okchaind
# HOME_CLI=/root/.okchaincli

# OKCHAIN_LAUNCH_TOP=/root/okchain/launch

# IP_INNET=true
# IP_PREFIX=192.168
# OKCHAIN_TESTNET_FULL_HOSTS=("192.168.13.121" "192.168.13.122" "192.168.13.123" "192.168.13.124" "192.168.13.125")
# HOSTS_PREFIX=okchain

# SEED_NODE_IP=192.168.13.116
# SCP_TAG="root"