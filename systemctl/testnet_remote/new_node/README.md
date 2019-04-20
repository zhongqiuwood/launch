## 新节点加入网络（非验证节点）

1. 修改`okchaind.profile`中的配置项
    * `SEED_NODE_ID` seed node id，可在seed node机器上执行`okchaind tendermint show-node-id`查看
    * `SEED_NODE_IP` seed node ip，注意区分内外网IP
    * `SEED_NODE_URL` seed node url，\<ip\>:\<port\>
    * `SEED_NODE_GENESIS` seed node的genesis file路径，必须可以`scp`到本机
    * `LOCAL_IP` 本机IP，注意区分内外网IP
    * `OKCHAIN_DAEMON` 本机`okchaind`可执行程序的路径
    * `HOME_DAEMON` 本机`okchaind`的`home`目录
2. 执行`./start.sh`