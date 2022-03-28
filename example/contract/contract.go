package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/xuper"
)

func main() {
	//testEVMContract()
	//testNativeContract()
	testWasmContract()

}

func testWasmContract() {
	codePath := "example/contract/data/counter.wasm"
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		panic(err)
	}

	xuperClient, err := xuper.New("127.0.0.1:37101")
	if err != nil {
		panic("new xuper Client error:")
	}

	// account, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	acc, err := account.GetAccountFromPlainFile("example/contract/data/keys")
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("retrieveAccount address: %v\n", acc.Address)
	contractAccount := "XC1111111111111111@xuper"
	timeStamp := fmt.Sprintf("%d", time.Now().Unix())
	contractName := fmt.Sprintf("counter%s", timeStamp[1:])
	fmt.Println(contractName)
	err = acc.SetContractAccount(contractAccount)

	args := map[string]string{
		"creator": "test",
		"key":     "test",
	}

	tx, err := xuperClient.DeployWasmContract(acc, contractName, code, args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deploy wasm Success!TxID:%x\n", tx.Tx.Txid)

	tx, err = xuperClient.InvokeWasmContract(acc, contractName, "increase", args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Invoke Wasm Contract Success! TxID:%x\n", tx.Tx.Txid)

	tx, err = xuperClient.QueryWasmContract(acc, contractName, "get", args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Query Wasm Contract Success! Response:%s\n", tx.ContractResponse.Body)
}

func testEVMContract() {
	binPath := "./example/contract/Counter.bin"
	abiPath := "./example/contract/Counter.abi"

	bin, err := ioutil.ReadFile(binPath)
	if err != nil {
		panic(err)
	}
	abi, err := ioutil.ReadFile(abiPath)
	if err != nil {
		panic(err)
	}

	// 通过助记词恢复账户
	account, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("retrieveAccount address: %v\n", account.Address)
	contractAccount := "XC1234567890123456@xuper"
	err = account.SetContractAccount(contractAccount)
	if err != nil {
		panic(err)
	}

	contractName := "SDKEvmContract"
	xchainClient, err := xuper.New("127.0.0.1:37101")

	args := map[string]string{
		"key": "test",
	}
	tx, err := xchainClient.DeployEVMContract(account, contractName, abi, bin, args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("DeployEVMContract SUCCESS! %x\n", tx.Tx.Txid)

	tx, err = xchainClient.InvokeEVMContract(account, contractName, "increase", args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("InvokeEVMContract SUCCESS! %x\n", tx.Tx.Txid)

	tx, err = xchainClient.QueryEVMContract(account, contractName, "get", args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("InvokeEVMContract Success! Response:%s\n", tx.ContractResponse.Body)

}

func testNativeContract() {
	codePath := "./example/contract/counter" // 编译好的二进制文件
	code, err := ioutil.ReadFile(codePath)
	if err != nil {
		panic(err)
	}

	account, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("retrieveAccount address: %v\n", account.Address)
	contractAccount := "XC1234567890123456@xuper"
	contractName := "SDKNativeCount1"
	err = account.SetContractAccount(contractAccount)
	if err != nil {
		panic(err)
	}

	xchainClient, err := xuper.New("127.0.0.1:37101")
	if err != nil {
		panic(err)
	}
	args := map[string]string{
		"creator": "test",
		"key":     "test",
	}
	tx, err := xchainClient.DeployNativeGoContract(account, contractName, code, args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Deploy Native Go Contract Success! %x\n", tx.Tx.Txid)

	tx, err = xchainClient.InvokeNativeContract(account, contractName, "increase", args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Invoke Native Go Contract Success! %x\n", tx.Tx.Txid)

	tx, err = xchainClient.QueryNativeContract(account, contractName, "get", args)
	if err != nil {
		panic(err)
	}
	if tx != nil {
		fmt.Printf("查询结果：%s\n", tx.ContractResponse.Body)
	}
}
