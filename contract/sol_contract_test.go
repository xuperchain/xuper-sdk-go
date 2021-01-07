// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/balance"
	"github.com/xuperchain/xuper-sdk-go/config"
	contractaccount "github.com/xuperchain/xuper-sdk-go/contract_account"
	"github.com/xuperchain/xuper-sdk-go/transfer"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

var (
	storageAbi = "[{\"inputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	storageBin = "608060405234801561001057600080fd5b5060c68061001f6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820281d7d36a78b2f2787eb4af7246d847cb385e9a6ea8d944ba283bb33246bbc8864736f6c63430005110032"
)

func TestSolDeploy(t *testing.T) {
	// abi := []byte("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
	// bin := []byte("608060405234801561001057600080fd5b506040516101453803806101458339818101604052602081101561003357600080fd5b81019080805190602001909291905050505060f2806100536000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80636057361d1460375780638f88708b146062575b600080fd5b606060048036036020811015604b57600080fd5b810190808035906020019092919050505060a1565b005b608b60048036036020811015607657600080fd5b810190808035906020019092919050505060ab565b6040518082815260200191505060405180910390f35b8060008190555050565b6000600182019150600054905091905056fea265627a7a72315820d55591d88547f8fb7a64116c747399f558c48d758753158c811c93272ec5d4bf64736f6c63430005110032")
	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	fmt.Println("accounr address:", acc.Address)
	node := "10.12.199.82:8101"
	bcname := "xuper"
	cName := "storage111222"
	cAccount := "XC9999999999999999@xuper"
	// createContractAccount(acc, node, bcname)// 已经创建合约账户了
	solContract := InitSolContract(acc, node, bcname, cName, cAccount)

	// args := map[string]string{
	// 	// "num": "1",
	// }
	r, e := solContract.Deploy(nil, []byte(storageBin), []byte(storageAbi))
	if e != nil {
		panic(e)
	}
	fmt.Println(r)
}

func TestSolInvoke(t *testing.T) {
	abi := []byte("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	fmt.Println("accounr address:", acc.Address)
	node := "10.12.199.82:8101"
	bcname := "xuper"
	cName := "aaaa"
	cAccount := "XC9999999999999999@xuper"
	solContract := InitSolContract(acc, node, bcname, cName, cAccount)

	args := map[string]string{
		"num": "5",
	}
	mName := "store"
	r, e := solContract.Invoke(mName, args, abi, "1")
	if e != nil {
		panic(e)
	}
	fmt.Println("invoke sucess:", r)

	x := initXchain()
	txStatus, err := x.QueryTx(r)
	if err != nil {
		panic(err)
	}
	vv, _ := json.Marshal(txStatus)

	fmt.Println("txStatus: ", string(vv))

	x.GetBalanceDetail()
}

func TestSolQuery(t *testing.T) {
	abi := []byte("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	fmt.Println("accounr address:", acc.Address)
	node := "10.12.199.82:8101"
	bcname := "xuper"
	cName := "aaaa"
	cAccount := "XC9999999999999999@xuper"
	solContract := InitSolContract(acc, node, bcname, cName, cAccount)

	args := map[string]string{
		"s": "2",
	}
	mName := "retrieve"
	preExeRPCRes, e := solContract.Query(mName, args, abi)
	if e != nil {
		panic(e)
	}
	gas := preExeRPCRes.GetResponse().GetGasUsed()
	fmt.Printf("gas used: %v\n", gas)
	fmt.Printf("preExeRPCRes: %v \n", preExeRPCRes)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		resu, err := ParseRespWithAbiForEVM(string(abi), mName, res)
		if err != nil {
			panic(err)
		}
		for _, v := range resu {
			fmt.Println("index:", v.Index)
			fmt.Println("value:", v.Value)
		}
	}
}

func TestTransfer(t *testing.T) {

	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	node := "10.12.199.82:8101"
	bcname := "xuper"
	b1 := balance.InitBalance(acc, node, []string{bcname})
	br1, e := b1.GetBalanceDetails()
	if e != nil {
		panic(e)
	}
	fmt.Println("转账前balanceDetails:", br1)

	fmt.Println("accounr address:", acc.Address)
	tf := transfer.InitTrans(acc, node, bcname)

	toa, err := account.CreateAccount(1, 1)
	if err != nil {
		panic(err)
	}
	fmt.Println("toAddr:", toa.Address)
	fmt.Println("toMenmonic:", toa.Mnemonic)
	r, e := tf.Transfer(toa.Address, "1", "1", "")
	if err != nil {
		panic(e)
	}
	fmt.Println("txID:", r)

	b := balance.InitBalance(acc, node, []string{bcname})
	br, e := b.GetBalanceDetails()
	if e != nil {
		panic(e)
	}
	fmt.Println("转账后balanceDetails:", br)

	b2 := balance.InitBalance(toa, node, []string{bcname})
	br2, e := b2.GetBalanceDetails()
	if e != nil {
		panic(e)
	}
	fmt.Println("转账后balanceDetails:", br2)
}

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
		XchainSer: "10.12.199.82:8101",
		ChainName: "xuper",
		Account:   acc,
	}
}
