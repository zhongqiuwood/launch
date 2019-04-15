
# 使用systemctl管理okchaind

本目录文件用于管理okchaind seed node

**注意修改`okchaind.profile`中参数**

1. 编译
```sh
go install github.com/okchain/okchain/cmd/okchaincli
go install github.com/okchain/okchain/cmd/okchaind
go build -o /root/go/src/github.com/cosmos/launch/launch /root/go/src/github.com/cosmos/launch/main.go
```
2. 加载systemctl service 文件
```sh
cp /root/go/src/github.com/cosmos/launch/systemctl/seednode/okchaind.service /etc/systemd/system

systemctl daemon-reload
```
3. 启动seed node okchaind服务
```sh
systemctl start okchaind
```
4. 停止seed node okchaind服务
```sh
systemctl stop okchaind
```
5. 重启seed node okchaind服务
```sh
systemctl restart okchaind
```
6. 查看systemctl状态
```sh
systemctl status okchaind
```

