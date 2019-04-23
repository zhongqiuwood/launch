# OKCHAIN Launch

## okchain launch 使用流程

### 1. 生成Captain账户 *<Boss操作>*
* 在一台离线机器创建`Captain账户`. 该账户将保存10亿个OKB
   ```sh
   okchaincli keys add captain --passwd 12345678
   ```
   将生成类似以下内容：
   ```
   NAME:   TYPE:   ADDRESS:                                        PUBKEY:
   captain local   okchain1gxsegjq6xp4a30tc8kc8h7w9ej7vfx8a2zrr05  okchainpub1addwnpepqwa2v88qt4mgee5p2s0jl7dppw32jzt6jp98txk88q2x9ens94vcy6uawvt

   **Important** write this mnemonic phrase in a safe place.
   It is the only way to recover your account if you ever forget your password.

   pulse whip pelican between bring decorate laptop abuse spend avoid pyramid judge
   ```
   ***务必保存好助记词***

* 将账户`ADDRESS`公开给研发人员，如`okchain1gxsegjq6xp4a30tc8kc8h7w9ej7vfx8a2zrr05`

### 2. 生成Admin账户 *<工程院操作>*
* 在另外一台机器创建`Admin账户`. 该账户为第一个Validator的委托账户, 负责启动OKChain第一个超级节点和生成创世块

   在初始块中为`Admin账户`分配1000000okb，用于创建Validator，日后将1000000okb返还给`Captain账户`
   ```
   okchaincli keys add admin --passwd 12345678
   ```
   ***务必保存好助记词***

### 3. 生成创世块文件和创世块交易 *<工程院操作>*
* 初始化okchaind, 生成genesis.json
   ```
    okchaind init --chain-id okchain
   ```
* 将"Admin账户"账户信息写入genesis.json
   ```
    okchaind add-genesis-account cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea 1000000okb
   ```
* 用"Admin账户"的私钥和密码生成创世块交易
   ```shell
    okchaind gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okchaind tendermint show-validator) --name admin
   ```
* 查看交易文件，`$HOME/.okchaind/config/gentx/`中的内容

### 4. 更新launch *<工程院操作>*
* 将上述两步骤中的执行结果放到okchain launch

   1. 将Captain账户地址及发币数量写入`launch/accounts/captain.json`中，格式如下：
      ```json
      [
      { "cosmos1kyh26rw89f8a4ym4p49g5z59mcj0xs4jd0wf8x": 1000000000}
      ]

      ```
   2. 将Admin账户地址及发币数量写入`launch/accounts/admin.json`中，格式如下：
      ```json
      [
      { "cosmos1m3gmu4zlnv2hmqfu2jwr97r2653w9yshxkhfea": 1}
      ]

      ```
   3. 将签名的交易文件完整复制到`launch/gentx/`中

* 提交更新到launch repo

### 5. 开源launch给社会, 任何人可以根据下面步骤加入OKChain网络
* 在`launch`下执行`go run main.go`生成最终的`genesis file`，即`launch/genesis.json`

* 利用launch的`genesis file`启动一个节点

   1. 初始化okchaind

      ```shell
       okchaind init --chain-id okchain
      ```

   1. 用`genesis file`覆盖`$HOME/.okchaind/config/genesis.json`，启动节点：

      ```shell
      okchaind start
      ```
