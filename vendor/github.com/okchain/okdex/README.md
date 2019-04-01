## 版本说明

1. `cosmos-sdk`使用`v0.33.0_patched` [Github](https://github.com/okblockchainlab/cosmos-sdk/tree/v0.33.0_patched)
2. `tendermint`使用`v0.31.0-rc0_patched` [Github](https://github.com/okblockchainlab/tendermint/tree/v0.31.0-rc0_patched)
3. 如果编译时遇到错误
    ```
    github.com/tendermint/iavl/nodedb.go:333:11: ndb.batch.Close undefined (type db.Batch has no field or method Close)
    ```
   需收到将`github.com/tendermint/iavl`切换到`v0.12.1`
   ```
   > go get github.com/tendermint/iavl

   > cd $GOPATH/src/github.com/tendermint/iavl

   > git checkout v0.12.1
   ```
4. 