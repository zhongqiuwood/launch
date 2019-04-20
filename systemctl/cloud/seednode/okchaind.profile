
OKCHAIN_DAEMON=/root/okchain/launch/systemctl/cloud/binary/okchaind
OKCHAIN_CLI=/root/okchain/launch/systemctl/cloud/binary/okchaincli
HOME_DAEMON=/root/.okchaind
HOME_CLI=/root/.okchaincli

OKCHAIN_LAUNCH_TOP=/root/okchain/launch

# 是否使用内网IP 
# 如果false，则使用公网IP，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点公网IP
# 如果是true，则必须设置IP_PREFIX，OKCHAIN_TESTNET_FULL_HOSTS设置其他节点内网IP
IP_INNET=false
IP_PREFIX=192.168
OKCHAIN_TESTNET_FULL_HOSTS=("3.112.83.58" "13.230.27.11" "13.231.103.82" "3.112.62.181" "3.112.95.199")

HOSTS_PREFIX=okchain_cloud

