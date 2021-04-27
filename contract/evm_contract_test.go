// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"fmt"
	"log"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/balance"
	"github.com/xuperchain/xuper-sdk-go/config"
	contractaccount "github.com/xuperchain/xuper-sdk-go/contract_account"
	"github.com/xuperchain/xuper-sdk-go/transfer"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

var (
	//storageAbi = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"storepay","outputs":[],"payable":true,"stateMutability":"payable","type":"function"}]`
	storageAbi = `[{"constant":true,"inputs":[{"name":"key","type":"string"}],"name":"get","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[],"name":"getOwner","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"key","type":"string"}],"name":"increase","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"inputs":[],"payable":false,"stateMutability":"nonpayable","type":"constructor"}]`

	storageBin = `608060405234801561001057600080fd5b50336000806101000a81548173ffffffffffffffffffffffffffffffffffffffff021916908373ffffffffffffffffffffffffffffffffffffffff1602179055506103d1806100606000396000f3fe6080604052600436106100345760003560e01c8063693ec85e14610039578063893d20e814610115578063ae896c871461016c575b600080fd5b34801561004557600080fd5b506100ff6004803603602081101561005c57600080fd5b810190808035906020019064010000000081111561007957600080fd5b82018360208201111561008b57600080fd5b803590602001918460018302840111640100000000831117156100ad57600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f820116905080830192505050505050509192919290505050610227565b6040518082815260200191505060405180910390f35b34801561012157600080fd5b5061012a61029c565b604051808273ffffffffffffffffffffffffffffffffffffffff1673ffffffffffffffffffffffffffffffffffffffff16815260200191505060405180910390f35b6102256004803603602081101561018257600080fd5b810190808035906020019064010000000081111561019f57600080fd5b8201836020820111156101b157600080fd5b803590602001918460018302840111640100000000831117156101d357600080fd5b91908080601f016020809104026020016040519081016040528093929190818152602001838380828437600081840152601f19601f8201169050808301925050505050505091929192905050506102c5565b005b60006001826040518082805190602001908083835b602083101515610261578051825260208201915060208101905060208303925061023c565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020549050919050565b60008060009054906101000a900473ffffffffffffffffffffffffffffffffffffffff16905090565b600180826040518082805190602001908083835b6020831015156102fe57805182526020820191506020810190506020830392506102d9565b6001836020036101000a038019825116818451168082178552505050505050905001915050908152602001604051809103902054036001826040518082805190602001908083835b60208310151561036b5780518252602082019150602081019050602083039250610346565b6001836020036101000a0380198251168184511680821785525050505050509050019150509081526020016040518091039020819055505056fea165627a7a723058201bd728e661baca0ac724c3636d4319e6d0287cf46df74c9d69282ffd03eb43540029`
	// storageAbi = `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// storageBin = "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	node = "127.0.0.1:37101"
)

func TestEVMDeploy(t *testing.T) {
	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)

	fmt.Println("account address:", acc.Address)
	// node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "evmCounter4"
	cAccount := "XC2222222222222222@xuper"
	//createContractAccount(acc, node, bcname) // 已经创建合约账户了

	EVMContract := InitEVMContract(acc, node, bcname, cName, cAccount)

	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	r, e := EVMContract.Deploy(args, []byte(storageBin), []byte(storageAbi))
	if e != nil {
		t.Error(e)
	}
	fmt.Println(r)
}

func TestEVMDeployWithClient(t *testing.T) {
	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)

	fmt.Println("account address:", acc.Address)
	// node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "evmCounter3"
	cAccount := "XC2222222222222222@xuper"
	//createContractAccount(acc, node, bcname) // 已经创建合约账户了

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}

	EVMContract := InitEVMContractWithClient(acc, bcname, cName, cAccount, sdkClient)

	args := map[string]string{
		"creator": "XC2222222222222222@xuper",
	}
	r, e := EVMContract.Deploy(args, []byte(storageBin), []byte(storageAbi))
	if e != nil {
		t.Error(e)
	}
	fmt.Println(r)
}

func testGetBalance(account *account.Account) {
	// 实例化一个交易相关的客户端对象

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		panic(err)
	}
	trans := transfer.InitTransWithClient(account, "xuper", sdkClient)
	// 查询账户的余额
	balance, err := trans.GetBalance()
	log.Printf("balance %v, err %v", balance, err)
}

func TestEVMInvoke(t *testing.T) {
	// abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`
	// bin := "608060405234801561001057600080fd5b506040516101203803806101208339818101604052602081101561003357600080fd5b8101908080519060200190929190505050806000819055505060c68061005a6000396000f3fe6080604052348015600f57600080fd5b506004361060325760003560e01c80632e64cec11460375780636057361d146053575b600080fd5b603d607e565b6040518082815260200191505060405180910390f35b607c60048036036020811015606757600080fd5b81019080803590602001909291905050506087565b005b60008054905090565b806000819055505056fea265627a7a72315820deacba9b51787b987df74d6ecd3bd463204d72726c7d7d97da0b0a8c62e8ccc364736f6c63430005110032"
	acc, _ := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	fmt.Println("accounr address:", acc.Address)
	// node := "127.0.0.1:37101"
	bcname := "xuper"
	cName := "evmCounter1"
	cAccount := "XC2222222222222222@xuper"
	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}

	EVMContract := InitEVMContractWithClient(acc, bcname, cName, cAccount, sdkClient)
	testGetBalance(acc)

	args := map[string]string{
		"key": "test",
	}
	mName := "increase"
	r, e := EVMContract.Invoke(mName, args, "111")
	if e != nil {
		t.Error(e)
	}
	fmt.Println("invoke sucess:", r)
	testGetBalance(acc)

	// x := initXchain()
	// txStatus, err := x.QueryTx(r)
	// if err != nil {
	// 	panic(err)
	// }
	// vv, _ := json.Marshal(txStatus)

	// fmt.Println("txStatus: ", string(vv))

	// x.GetBalanceDetail()
}

//func TestEVMQuery(t *testing.T) {
//	// abi := []byte(`[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"}]`)
//	// abi := []byte("[{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"constant\":true,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"s\",\"type\":\"uint256\"}],\"name\":\"retrieve\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"num\",\"type\":\"uint256\"}],\"name\":\"store\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]")
//	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	fmt.Println("accounr address111:", acc.Address)
//
//	// node := "127.0.0.1:37101"
//	bcname := "xuper"
//	cName := "storageA"
//	cAccount := "XC2222222222222222@xuper"
//	EVMContract := InitEVMContractWithClient(acc, node, bcname, cName, cAccount)
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

func TestTransfer(t *testing.T) {

	acc, _ := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	node := "127.0.0.1:37101"
	bcname := "xuper"
	b1 := balance.InitBalance(acc, node, []string{bcname})
	br1, e := b1.GetBalanceDetails()
	if e != nil {
		panic(e)
	}
	fmt.Println("转账前balanceDetails:", br1)

	fmt.Println("accounr address:", acc.Address)
	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}

	tf := transfer.InitTransWithClient(acc, bcname, sdkClient)

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
		XchainSer: node,
		ChainName: "xuper",
		Account:   acc,
	}
}
