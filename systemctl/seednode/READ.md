
# 使用systemctl管理okdexd

本目录文件用于管理okdexd seed node

**注意修改`okdexd.profile`中参数**

1. 加载systemctl service 文件
```sh
cp okdexd.service /etc/systemd/system
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