## Usage

### 可信账本测试流程
使用默认背书.
1. 根据xuper.json作为创始块配置，搭建dpos网络；
2. 使用main.go的testAccount创建ak, 使用account.json创建合约账号,并且注意给ak充值 
```
./xchain-cli account new --desc /root/xuper-sdk-go/example/account.json  --fee=1000
```
3. 运行 `main_trust_counter`
