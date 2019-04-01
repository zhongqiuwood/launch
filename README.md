# OKDEX Launch

## okdex launch 使用流程

### 1. 生成创世账户
* 在一台离线机器创建"创世账户". 该账户将保存10亿个OKB
    ```
    okdexcli keys add genesis_addr
    
    {
      "name": "genesis_addr",
      "type": "local",
      "address": "cosmos1kyh26rw89f8a4ym4p49g5z59mcj0xs4jd0wf8x",
      "pubkey": "cosmospub1addwnpepq02z7p4f7p93yzeyf4r7nk2cuajlvp523rjjdcspxzd784thpdvau57lsmu",
      "mnemonic": "******"
    }
    ```
* 将创世账户地址公开给研发人员 "cosmos1kyh26rw89f8a4ym4p49g5z59mcj0xs4jd0wf8x"

### 2. 生成Admin账户
* 在另外一台机器创建"Admin账户". 该账户为第一个Validator的委托账户, 负责启动OKChain第一个超级节点, 和生成创世块
```
okdexcli keys add admin
{
  "name": "admin",
  "type": "local",
  "address": "cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea",
  "pubkey": "cosmospub1addwnpepqwrdekewcwy6vmkl8mqu5uec29amyjqqxkt8wd3exjd0fy8pn4vqg6wquvz",
  "mnemonic": "keen border system oil inject hotel hood potato shed pumpkin legend actor"
}
```

### 3. 生成创世块文件，和创世块交易签名
* 用Admin账户生成创世交易签名: Admin账户质押自己OKB成为第一个Validator的委托人
* 初始化okdexd, 生成genesis.json
```
    okdexd init --chain-id okchain
```

* 将"Admin账户"账户信息写入genesis.json
```
    okdexd add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 1okb
```
* 用"Admin账户"的私钥和密码生成创世块交易签名
```shell
    okdexd gentx --amount 1okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okdexd tendermint show-validator) --name admin
```

* 查看交易，`$HOME/.okdexd/config/gentx/`中的文件内容


### 4. 更新launch
* 将上述两步骤中的执行结果放到okdex launch

   1. 将账户地址及发币数量写入`launch/accounts/distribution.json`中，格式如下：

      ```json
      [
      { "cosmos1kyh26rw89f8a4ym4p49g5z59mcj0xs4jd0wf8x": 1000000000},
      { "cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea": 1}
      ]

      ```

   1. 将签名的交易内容完整复制到`launch/gentx/`中，格式如下：

      ```
      okdexd collect-gentxs
      ```
* 提交更新到launch repo


### 5. 开源launch给社会, 任何人可以根据下面步骤加入OKChain网络
* 在`launch`下执行`go run main.go`生成最终的`genesis file`，即`launch/genesis.json`

* 利用launch的`genesis file`启动一个节点

   1. 初始化okdexd

      ```shell
       okdexd init --chain-id okchain
      ```

   1. 用`genesis file`覆盖`$HOME/.okdexd/config/genesis.json`，启动节点：

      ```shell
      okdexd start
      ```
