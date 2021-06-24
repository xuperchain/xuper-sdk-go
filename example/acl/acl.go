package main

import (
	"fmt"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	// XuperChain 可以为合约账户以及合约方法设置 ACL，下面分别使用 sdk 来设置 ACL。
	setAccountACLExample()
	setMethodACLExample()
}

func setAccountACLExample() {
	// 假设你已经有一个账户，助记词为：玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即，同时已经有了对应的合约账户：XC8888888899999999@xuper，
	// 同时这个账户可以操作对应的合约账户。
	mnemonic := "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
	bob, _ := account.RetrieveAccount(mnemonic, 1)

	// 如果想修改 XC8888888899999999@xuper 的 ACL，首先需要设置普通账户的合约账户。
	contractAcc := "XC8888888899999999@xuper"
	bob.SetContractAccount(contractAcc)

	// 创建节点客户端。
	node := "127.0.0.0:37101"
	xclient, err := xuper.New(node)
	if err != nil {
		panic(err)
	}
	defer xclient.Close()

	// 设置想要的 ACL，假设想要设置两个普通账户一起签名才能使用合约账户，另外一个账户地址为：nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh
	acl := xuper.NewACL(1, 0.6)
	acl.AddAK("nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh", 0.3)
	acl.AddAK(bob.Address, 0.3)
	fmt.Println(acl)

	tx, err := xclient.SetAccountACL(bob, acl)
	if err != nil {
		panic(err)
	}
	// 查看本次交易的 gas。
	fmt.Println(tx.GasUsed)

	queryACL, err := xclient.QueryAccountACL(bob.GetContractAccount())
	if err != nil {
		panic(err)
	}
	// 查看修改后的 ACL。
	fmt.Println(queryACL)
}

func setMethodACLExample() {
	// 设置合约方法的 ACL 与设置合约账户 ACL 类似，
	// 同样假设你已经有一个账户，助记词为：玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即，同时已经有了对应的合约账户：XC8888888899999999@xuper，
	// 另外假设你也有了一个合约 counter，同时还有一个方法 increas 已经部署了。
	// 现在要设置 counter 合约的 increase 方法只能你自己调用。

	mnemonic := "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
	bob, _ := account.RetrieveAccount(mnemonic, 1)

	// 设置普通账户的合约账户。
	contractAcc := "XC8888888899999999@xuper"
	bob.SetContractAccount(contractAcc)

	// 创建节点客户端。
	node := "127.0.0.0:37101"
	xclient, err := xuper.New(node)
	if err != nil {
		panic(err)
	}
	defer xclient.Close()

	// 设置你的账户才有权限调用。
	acl := xuper.NewACL(1, 1.0)
	acl.AddAK(bob.Address, 1.0)
	fmt.Println(acl)

	tx, err := xclient.SetMethodACL(bob, "counter", "increase", acl)
	if err != nil {
		panic(err)
	}
	// 查看本次交易的 gas。
	fmt.Println(tx.GasUsed)

	queryACL, err := xclient.QueryMethodACL("counter", "increase")
	if err != nil {
		panic(err)
	}
	// 查看修改后的 ACL。
	fmt.Println(queryACL)
}
