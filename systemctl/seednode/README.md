
# 使用systemctl管理okdexd

本目录文件用于管理okdexd seed node

**注意修改`okdexd.profile`中参数**

1. 编译
```sh
go install github.com/okchain/okdex/cmd/okdexcli
go install github.com/okchain/okdex/cmd/okdexd
go build -o /root/go/src/github.com/cosmos/launch/launch /root/go/src/github.com/cosmos/launch/main.go
```
1. 加载systemctl service 文件
```sh
cp /root/go/src/github.com/cosmos/launch/systemctl/seednode/okdexd.service /etc/systemd/system
```

2. 启动seed node okdexd服务
```sh
systemctl start okdexd.service
```

3. 停止seed node okdexd服务
```sh
systemctl stop okdexd.service
```

4. 重启seed node okdexd服务
```sh
systemctl restart okdexd.service
```

5. 如果未能正常启动，可查看systemctl状态
```sh
systemctl status okdexd.service
```

