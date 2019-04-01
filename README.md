# OKDEX Launch

## okdex launch 使用流程

1. 创建账户

   ```shell
   okdexcli keys add boss
   ```

   查看账户地址

   ```shell
   okdexcli keys show boss -a
   ```

1. 生成签名交易

   1. 初始化okdexd

      ```shell
       okdexd init --chain-id okchain
      ```

   1. 将账户信息写入genesis file

      ```shell
      okdexd add-genesis-account $(okdexcli keys show boss -a) 1000000000okb
      ```

   1. 生成签名交易

      ```shell
      okdexd gentx --amount 1000000okb --min-self-delegation 1 --commission-rate 0.1 --commission-max-rate 0.5 --commission-max-change-rate 0.001 --pubkey $(okdexd tendermint show-validator) --name boss
      ```

   1. 查看交易，`$HOME/.okdexd/config/gentx/`中的文件内容

1. 将上述两步骤中的执行结果放到okdex launch

   1. 将账户地址及发币数量写入`launch/accounts/initOKB.json`中，格式如下：

      ```json
      {
        "cosmos14s3dfqterut5flk9py9yurve7kvjwrp52e2ufe": 1000000000
      }
      ```

   1. 将签名的交易内容完整复制到`launch/gentx/`中，格式如下：

      ```json
      {"type":"auth/StdTx","value":{"msg":[{"type":"cosmos-sdk/MsgCreateValidator","value":{"description":{"moniker":"yulinshengdeMacBook-Pro.local","identity":"","website":"","details":""},"commission":{"rate":"0.100000000000000000","max_rate":"0.500000000000000000","max_change_rate":"0.001000000000000000"},"min_self_delegation":"1","delegator_address":"cosmos14s3dfqterut5flk9py9yurve7kvjwrp52e2ufe","validator_address":"cosmosvaloper14s3dfqterut5flk9py9yurve7kvjwrp50d7f92","pubkey":"cosmosvalconspub1zcjduepqa9ad9ksej6ywkzne3vcle4vewglq5xcan4km7x4vp5uw45qcsdkqsxskrv","value":{"denom":"okb","amount":"1000000"}}}],"fee":{"amount":null,"gas":"200000"},"signatures":[{"pub_key":{"type":"tendermint/PubKeySecp256k1","value":"Aw6QKYLwNVyVHPofaxKOUTtOkUy3nO0OiKjEtpqsslxW"},"signature":"dnk3828ZvLWYw76WTqXTzeD2CjR8TJJndelCY6R3XGU9cXyCw2wCu6/pB7e6Xu8++Y/vjjHt0VdmZstHliDHmQ=="}],"memo":"5aa3315b66480b9a0575dd0c67e2469c00388be9@192.168.26.129:26656"}}
      ```

1. 在`launch`下执行`go run main.go`生成最终的`genesis file`，即`launch/genesis.json`

1. 利用launch的`genesis file`启动一个节点

   1. 初始化okdexd

      ```shell
       okdexd init --chain-id okchain
      ```

   1. 用`genesis file`覆盖`$HOME/.okdexd/config/genesis.json`，启动节点：

      ```shell
      okdexd start
      ```