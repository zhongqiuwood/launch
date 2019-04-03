
# 使用systemctl管理okdexd

本目录文件用于管理okdexd full node

**注意修改`okdexd.profile`中参数**

1. 加载systemctl service 文件
```sh
cp okdexd.service /etc/systemd/system
```
2. 初始化okdexd
```sh
okdexd init --chain-id okchain --home /root/.okdexd
```

3. 从seed node处获取genesis file
```sh
scp root@192.168.13.116:/root/.okdexd/config/genesis.json /root/.okdexd/config
```

4. 启动seed node okdexd服务
```sh
systemctl start okdexd.service
```

5. 停止seed node okdexd服务
```sh
systemctl stop okdexd.service
```

6. 重启seed node okdexd服务
```sh
systemctl restart okdexd.service
```

*说明：以上2、3两步理应写到`okdexd_start.sh`中，但测试一直有错，所以临时单独拿出来操作，以后继续完善*