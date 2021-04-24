// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"fmt"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	contractaccount "github.com/xuperchain/xuper-sdk-go/contract_account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

var (
	storageAbi = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"storepay","outputs":[],"payable":true,"stateMutability":"payable","type":"function"}]`
	storageBin = `608060405234801561001057600080fd5b5060405161016c38038061016c8339818101604052602081101561003357600080fd5b810190808051906020019092919050505080600081905550506101118061005b6000396000f3fe60806040526004361060305760003560e01c80632e64cec11460355780636057361d14605d5780638995db74146094575b600080fd5b348015604057600080fd5b50604760bf565b6040518082815260200191505060405180910390f35b348015606857600080fd5b50609260048036036020811015607d57600080fd5b810190808035906020019092919050505060c8565b005b60bd6004803603602081101560a857600080fd5b810190808035906020019092919050505060d2565b005b60008054905090565b8060008190555050565b806000819055505056fea265627a7a723158209500c3e12321b837819442c0bc1daa92a4f4377fc7b59c41dbf9c7620b2f961064736f6c63430005110032`
	// storageAbi = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// storageBin = "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	node = "127.0.0.1:37101"
)

func TestEVMDeploy(t *testing.T) {
	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)

	fmt.Println("account address:", acc.Address)
	node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "storageA1"
	cAccount := "XC2222222222222222@xuper"
	//createContractAccount(acc, node, bcname) // 已经创建合约账户了

	sdkClient,err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}

	EVMContract := InitEVMContract(acc, bcname, cName, cAccount,sdkClient)
	args := map[string]string{
		"num": "1",
	}
	r, e := EVMContract.Deploy(args, []byte(storageBin), []byte(storageAbi))
	if e != nil {
		t.Error(e)
	}
	fmt.Println(r)
}

//func testGetBalance(account *account.Account) {
//	// 实例化一个交易相关的客户端对象
//	trans := transfer.InitTrans(account, node, "xuper")
//	// 查询账户的余额
//	balance, err := trans.GetBalance()
//	log.Printf("balance %v, err %v", balance, err)
//}

//func TestEVMInvoke(t *testing.T) {
//	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
//	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
//	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	fmt.Println("accounr address:", acc.Address)
//	// node := "127.0.0.1:37101"
//	bcname := "xuper"
//	cName := "storageA"
//	cAccount := "XC9999999999999999@xuper"
//	EVMContract := InitEVMContract(acc, node, bcname, cName, cAccount)
//	testGetBalance(acc)
//
//	args := map[string]string{
//		"num": "5882",
//	}
//	mName := "store"
//	r, e := EVMContract.Invoke(mName, args, "111")
//	if e != nil {
//		t.Error(e)
//	}
//	fmt.Println("invoke sucess:", r)
//	testGetBalance(acc)
//
//	// x := initXchain()
//	// txStatus, err := x.QueryTx(r)
//	// if err != nil {
//	// 	panic(err)
//	// }
//	// vv, _ := json.Marshal(txStatus)
//
//	// fmt.Println("txStatus: ", string(vv))
//
//	// x.GetBalanceDetail()
//}

//func TestEVMQuery(t *testing.T) {
//	// abi := []byte(`[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`)
//	// abi := []byte("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
//	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	fmt.Println("accounr address111:", acc.Address)
//
//	// node := "127.0.0.1:37101"
//	bcname := "xuper"
//	cName := "storageA"
//	cAccount := "XC9999999999999999@xuper"
//	EVMContract := InitEVMContract(acc, node, bcname, cName, cAccount)
//
//	// args := map[string]string{
//	// 	"s": "2",
//	// }
//	mName := "retrieve"
//	preExeRPCRes, e := EVMContract.Query(mName, nil)
//	if e != nil {
//		t.Error(e)
//	}
//	gas := preExeRPCRes.GetResponse().GetGasUsed()
//	fmt.Printf("gas used: %v\n", gas)
//	fmt.Printf("preExeRPCRes: %v \n", preExeRPCRes)
//	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
//		fmt.Printf("contract response: %s\n", string(res))
//	}
//}
//
//func TestTransfer(t *testing.T) {
//
//	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	node := "127.0.0.1:37101"
//	bcname := "xuper"
//	b1 := balance.InitBalance(acc, node, []string{bcname})
//	br1, e := b1.GetBalanceDetails()
//	if e != nil {
//		panic(e)
//	}
//	fmt.Println("转账前balanceDetails:", br1)
//
//	fmt.Println("accounr address:", acc.Address)
//	tf := transfer.InitTrans(acc, node, bcname)
//
//	toa, err := account.CreateAccount(1, 1)
//	if err != nil {
//		panic(err)
//	}
//	fmt.Println("toAddr:", toa.Address)
//	fmt.Println("toMenmonic:", toa.Mnemonic)
//	r, e := tf.Transfer(toa.Address, "1", "1", "")
//	if err != nil {
//		panic(e)
//	}
//	fmt.Println("txID:", r)
//
//	b := balance.InitBalance(acc, node, []string{bcname})
//	br, e := b.GetBalanceDetails()
//	if e != nil {
//		panic(e)
//	}
//	fmt.Println("转账后balanceDetails:", br)
//
//	b2 := balance.InitBalance(toa, node, []string{bcname})
//	br2, e := b2.GetBalanceDetails()
//	if e != nil {
//		panic(e)
//	}
//	fmt.Println("转账后balanceDetails:", br2)
//}

func createContractAccount(acc *account.Account, node, bcname string) {
	ca := contractaccount.InitContractAccount(acc, node, bcname)
	r, e := ca.CreateContractAccount("XC9999999999999999@xuper")
	if e != nil {
		panic(e)
	}
	fmt.Println("createCoutractAccount SUCCESS:", r)
}

func initXchain() *xchain.Xchain {
	commConfig := config.GetInstance()

	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RetrieveAccount: %v\n", acc)

	return &xchain.Xchain{
		Cfg:       commConfig,
		XchainSer: node,
		ChainName: "xuper",
		Account:   acc,
	}
}
