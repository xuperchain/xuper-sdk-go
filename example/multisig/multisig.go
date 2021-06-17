package main

import (
	"fmt"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

// 假设有一个合约账户，两个普通地址（AK），需要两个地址同时签名才可以使用合约账户，此时多签的示例如下。
func main() {
	me, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}

	// 如果没有合约账户需要先创建合约账户。
	contractAccount := "XC1234567812345678@xuper"
	me.SetContractAccount(contractAccount)
	fmt.Println(me.Address)

	// 创建链的客户端。
	node := "127.0.0.1:37101"
	xclient, _ := xuper.New(node)

	// 创建本次交易的请求数据。
	code := []byte{}
	args := map[string]string{
		"creator": "bob",
	}
	request, err := xuper.NewDeployContractRequest(me, "counter", nil, code, args, "wasm", "c")
	if err != nil {
		panic(err)
	}

	// 构造本次交易的数据结构，此时交易还未发送到链上。
	// 由于使用多签的形式，我们还需要其他账户对此交易结构进行签名。
	tx, err := xclient.GenerateTx(request)
	if err != nil {
		panic(err)
	}

	// 可以将 tx 数据通过网络传输给其他服务，也可以直接使用账户对 tx 进行签名。最后只需要将收集到签名的 tx 返回即可。

	// 使用 bob 账户对 tx 进行签名。
	bob, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	bob.SetContractAccount(contractAccount)
	tx.Sign(bob)

	// 收集到足够的签名后，将交易发送出去。
	tx, err = xclient.PostTx(tx)
	if err != nil {
		panic(err)
	}

	// 打印出交易的 ID。
	fmt.Printf("%x\n", tx.Tx.Txid)
}
