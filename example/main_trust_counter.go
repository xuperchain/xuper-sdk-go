package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/contract"
	"github.com/xuperchain/xuper-sdk-go/contract_account"
	"github.com/xuperchain/xuper-sdk-go/transfer"
)

var language = 1
var contractAcc = "XC1111111111111122@xuper"
var transactionId = ""

// define blockchain node and blockchain name
var (
	contractName = "counter3"
	node         = "localhost:37101" // node ip
	bcname       = "xuper"
)

func testAccount() {
	if _, err := os.Stat("./keys"); err != nil && os.IsNotExist(err) {
	} else {
		println("existed, pass")
		return
	}
	// create an account for the user,
	// strength 1 means that the number of mnemonics is 12
	// language 1 means that mnemonics is Chinese
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		fmt.Printf("create account error: %v\n", err)
		panic(err)
	}
	// print the account, mnemonics
	fmt.Println(acc)
	fmt.Println("hello, Mnemonic: ", acc.Mnemonic)

	// retrieve the account by mnemonics
	acc, err = account.RetrieveAccount(acc.Mnemonic, 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		panic(err)
	}
	fmt.Printf("RetrieveAccount: to %v\n", acc)

	// create an account, then encrypt using password and save it to a file
	acc, err = account.CreateAndSaveAccountToFile("./keys", "123", 1, 1)
	if err != nil {
		fmt.Printf("createAndSaveAccountToFile err: %v\n", err)
		panic(err)
	}
	fmt.Printf("CreateAndSaveAccountToFile: %v\n", acc)

	// get the account from file, using password decrypt
	acc, err = account.GetAccountFromFile("keys/", "123")
	if err != nil {
		fmt.Printf("getAccountFromFile err: %v\n", err)
		panic(err)
	}
	fmt.Printf("getAccountFromFile: %v\n", acc)
	return
}
func usingAccount() (*account.Account, error) {
	// load your account from the private key and secure code you download from xuper.baidu.com
	// Note that put the downloaded private key file at path "./keys/private.key"
	acc, err := account.GetAccountFromFile("./keys/", "123")
	if err != nil {
		return nil, fmt.Errorf("create account error: %v\n", err)
	}
	// print the account, mnemonics
	fmt.Println(acc.Address)

	return acc, nil
}

func testContractAccount() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	//account, err := account.RetrieveAccount(Mnemonics, language)
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		panic(err)
	}

	// define the name of the conrtact account to be created
	// Notice !!!
	// conrtact account is (XC + 16 length digit + @xuper), like "XC1234567890123456@xuper"
	contractAccount := contractAcc

	// initialize a client to operate the contract account
	ca := contractaccount.InitContractAccount(acc, node, bcname)

	// create contract account
	txid, err := ca.CreateContractAccount(contractAccount)
	if err != nil {
		log.Printf("CreateContractAccount err: %v", err)
		os.Exit(-1)
	}
	fmt.Println(txid)
	return
}

func testTransfer() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	/*
		acc, err := account.RetrieveAccount(Mnemonics, language)
	*/
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		panic(err)
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	to := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount := "200"
	fee := "0"
	desc := ""
	// transfer
	txid, err := trans.Transfer(to, amount, fee, desc)
	if err != nil {
		fmt.Printf("Transfer err: %v\n", err)
		panic(err)
	}
	fmt.Printf("transfer tx: %v\n", txid)
	return
}

func testDeployWasmContract() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	//acc, err := account.RetrieveAccount(Mnemonics, language)
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		panic(err)
	}
	fmt.Printf("account: %v\n", acc)

	// set contract account, contract will be installed in the contract account
	// Notice !!!
	contractAccount := contractAcc

	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set init args and contract file
	args := map[string]string{
		"creator": "xchain",
	}
	codePath := "example/contract_code/trust_counter.wasm"

	// deploy wasm contract
	txid, err := wasmContract.DeployWasmContract(args, codePath, "c")
	if err != nil {
		log.Printf("DeployWasmContract err: %v", err)
		panic(err)
	}
	fmt.Printf("DeployWasmContract txid: %v\n", txid)
	return
}

func testInvokeWasmContract() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	//acc, err := account.RetrieveAccount(Mnemonics, language)
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the contract
	contractAccount := contractAcc
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// test store
	log.Println("......start store......")
	args := map[string]string{
		"duan": "25",
		"bing": "12",
	}
	args, err = wasmContract.EncryptWasmArgs(args);
	if err != nil {
		panic(err)
	}
	methodName := "store"
	txid, err := wasmContract.InvokeWasmContract(methodName, args)
	if err != nil {
		log.Printf("InvokeWasmContract PostWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	log.Printf("txid: %v", txid)
	transactionId = txid
	log.Printf("check for duan: %v", check("duan", args["duan"]))
	log.Printf("check for bing: %v", check("bing", args["bing"]))
	log.Printf("......store finished......\n\n")

	// test addition
	log.Println("start addition......")
	addArgs := map[string]string{
		"l": "duan",
		"r": "bing",
		"o": "duanbing",
	}
	methodName = "add"
	txid, err = wasmContract.InvokeWasmContract(methodName, addArgs)
	if err != nil {
		log.Printf("InvokeWasmContract PostWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	log.Printf("txid: %v", txid)
	transactionId = txid

	log.Printf("check for addition: %v", check("duanbing", "37"))
	log.Printf(" ......addition finished......\n\n")
	return
}

func queryPlainValue(key string) string{
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}

	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAcc)
	// set query function method and args
	args := map[string]string{
		"key": key,
	}
	methodName := "get"

	// 获取加密的数据
	preExeRPCRes, err := wasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	// 对加密数据解密
	resPlain,err := wasmContract.DecryptResponse(preExeRPCRes)
	if err != nil {
		log.Printf("QueryWasmContractPlain failed, err: %v", err)
		os.Exit(-1)
	}
    // 返回解密的value
	return string(resPlain.GetResponse().GetResponse()[0])
}

func check(key, plain string) bool {
    value := queryPlainValue(key)
    return value == plain
}

func testQueryWasmContract() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	//acc, err := account.RetrieveAccount(Mnemonics, language)
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAcc)

	// set query function method and args
	args := map[string]string{
		"key": "duan",
	}
	methodName := "get"

	// query contract
	preExeRPCRes, err := wasmContract.QueryWasmContract(methodName, args)
	if err != nil {
		log.Printf("QueryWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	gas := preExeRPCRes.GetResponse().GetGasUsed()
	fmt.Printf("gas used: %v\n", gas)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response: %s\n", string(res))
	}
	return
}

func testGetBalance() {
	// retrieve the account by mnemonics
	// Notice !!!
	// parameters should be Mnemonics for your account and language
	//acc, err := account.RetrieveAccount(Mnemonics, language)
	acc, err := usingAccount()
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// get balance of the account
	balance, err := trans.GetBalance()
	log.Printf("balance %v, err %v", balance, err)
	return
}

func testQueryTx() {
	// initialize a client to operate the transaction
	time.Sleep(10 * time.Second)
	trans := transfer.InitTrans(nil, node, bcname)

	// query tx by txid
	tx, err := trans.QueryTx(transactionId)
	log.Printf("query tx %v, err %v", tx, err)
	return
}

func main() {
	contractName = contractName + fmt.Sprintf("%d", time.Now().Unix()%1000000)

	//testContractAccount()
	testAccount()
	testTransfer()
	testDeployWasmContract()
	testInvokeWasmContract()
	testQueryTx()
	testQueryWasmContract()
	testGetBalance()
	println("contractname: ", contractName)
}
