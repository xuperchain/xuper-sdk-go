package main

import (
	"fmt"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	// 多签场景有两种方式可以选择，主要流程都是先创建好交易，然后收集到足够的签名再发送到链上。
	multisign1()
	multisign2()
}

// 构建 Request 再构造 Transaction，然后签名再 post 到链。
func multisign1() {
	// 假设有一个合约账户，两个普通地址（AK），需要两个地址同时签名才可以使用合约账户，此时多签的示例如下。
	// 已有账户 alice 和 bob，两个账户同时签名才可以使用合约账户 XC1234567812345678@xuper。
	contractAccount := "XC1234567812345678@xuper"
	alice, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	alice.SetContractAccount(contractAccount)

	bob, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	bob.SetContractAccount(contractAccount)

	// 创建链的客户端。
	node := "127.0.0.1:37101"
	xclient, _ := xuper.New(node)

	// 创建本次交易的请求数据。
	code := []byte{}
	args := map[string]string{
		"creator": "bob",
	}

	// 首先使用 alice 构造交易，由于知道还需要 bob 账户签名，因此需要增加 bob 的 AuthRequire。
	authRequire := []string{bob.GetAuthRequire()}
	request, err := xuper.NewDeployContractRequest(alice, "counter", nil, code, args, "wasm", "c", xuper.WithOtherAuthRequires(authRequire))
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
	tx.Sign(bob)

	// 收集到足够的签名后，将交易发送出去。
	tx, err = xclient.PostTx(tx)
	if err != nil {
		panic(err)
	}

	// 打印出交易的 ID。
	fmt.Printf("%x\n", tx.Tx.Txid)
}

// 使用 xuperclient 构造交易。
func multisign2() {
	// 假设有一个合约账户，两个普通地址（AK），需要两个地址同时签名才可以使用合约账户，此时多签的示例如下。
	// 已有账户 alice 和 bob，两个账户同时签名才可以使用合约账户 XC1234567812345678@xuper。
	contractAccount := "XC1234567812345678@xuper"
	alice, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	alice.SetContractAccount(contractAccount)

	bob, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	bob.SetContractAccount(contractAccount)

	// 创建链的客户端。
	node := "127.0.0.1:37101"
	xclient, _ := xuper.New(node)

	// 创建本次交易的请求数据。
	code := []byte{}
	args := map[string]string{
		"creator": "bob",
	}

	// 首先使用 alice 构造交易，由于知道还需要 bob 账户签名，因此需要增加 bob 的 AuthRequire。
	authRequire := []string{bob.GetAuthRequire()}

	// 使用 xuper.WithNotPost() 表明只构造交易，不将交易 post 到链上。
	// xuper.WithOtherAuthRequires() 表明还需要增加 bob 的签名。
	tx, err := xclient.DeployWasmContract(alice, "counter", code, args, xuper.WithNotPost(), xuper.WithOtherAuthRequires(authRequire))
	if err != nil {
		panic(err)
	}

	// 可以将 tx 数据通过网络传输给其他服务，也可以直接使用账户对 tx 进行签名。最后只需要将收集到签名的 tx 返回即可。

	// 使用 bob 账户对 tx 进行签名。
	tx.Sign(bob)

	// 收集到足够的签名后，将交易发送出去。
	tx, err = xclient.PostTx(tx)
	if err != nil {
		panic(err)
	}

	// 打印出交易的 ID。
	fmt.Printf("%x\n", tx.Tx.Txid)
}
