
# 使用systemctl管理okdexd

本目录文件用于管理okdexd full node

**注意修改`okdexd.profile`中参数**

1. 编译
```sh
go install github.com/okchain/okdex/cmd/okdexcli
go install github.com/okchain/okdex/cmd/okdexd
```
2. 加载systemctl service 文件
```sh
cp /root/go/src/github.com/cosmos/launch/systemctl/fullnode/okdexd.service /etc/systemd/system

systemctl daemon-reload
```
3. 启动seed node okdexd服务
```sh
systemctl start okdexd
```
4. 停止seed node okdexd服务
```sh
systemctl stop okdexd
```
5. 重启seed node okdexd服务
```sh
systemctl restart okdexd
```
6. 查看systemctl状态
```sh
systemctl status okdexd
```