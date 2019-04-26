
## 1. 测试网络部署

1. 执行```./git.sh```编译并上传最新的binaries到git, 再从git下载最新binaries
1. 执行```./start -r``` 启动系统

## 2. ```start.sh``` usage
```sh
./start.sh
```
`-c` 清理所有相关数据后，启动节点

`-s` 停止所有机器上的okchaind

`-t` 发币、上币提案以及提案投票

`-a` 上币提案Passed后，激活

```sh
./git.sh
```
更新远程机器launch代码库
`-c` git clone

## 3. 注意事项

**对于不同的机器环境，请务必首先修改以下内容：**
1. **`okchaind.profile`中的配置项**
