package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/balance"
	"github.com/xuperchain/xuper-sdk-go/contract"
	"github.com/xuperchain/xuper-sdk-go/contract_account"
	"github.com/xuperchain/xuper-sdk-go/crypto"
	"github.com/xuperchain/xuper-sdk-go/network"
	"github.com/xuperchain/xuper-sdk-go/transfer"

	hdapi "github.com/xuperchain/crypto/gm/hdwallet/api"
)

// define blockchain node and blockchain name
var (
	contractName = "counter"

	// node for test network of XuperOS
	// node = "14.215.179.74:37101"

	// node for official network of XuperOS
	node = "39.156.69.83:37100"

	//	node         = "127.0.0.1:37801"
	bcname = "xuper"
)

func testAccount() {
	// create an account for the user,
	// strength 1 means that the number of mnemonics is 12
	// language 1 means that mnemonics is Chinese
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		fmt.Printf("create account error: %v\n", err)
	}
	// print the account, mnemonics
	fmt.Println(acc)
	fmt.Println(acc.Mnemonic)

	// retrieve the account by mnemonics
	acc, err = account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("RetrieveAccount: to %v\n", acc)

	// create an account, then encrypt using password and save it to a file
	acc, err = account.CreateAndSaveAccountToFile("./keys", "123", 1, 1)
	if err != nil {
		fmt.Printf("createAndSaveAccountToFile err: %v\n", err)
	}
	fmt.Printf("CreateAndSaveAccountToFile: %v\n", acc)

	// get the account from file, using password decrypt
	acc, err = account.GetAccountFromFile("keys/", "123")
	if err != nil {
		fmt.Printf("getAccountFromFile err: %v\n", err)
	}
	fmt.Printf("getAccountFromFile: %v\n", acc)
	return
}

func testContractAccount() {
	// retrieve the account by mnemonics
	account, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}

	fmt.Printf("retrieveAccount address: %v\n", account.Address)

	// define the name of the conrtact account to be created
	// conrtact account is (XC + 16 length digit + @xuper), like "XC1234567890123456@xuper"
	contractAccount := "XC1234567890123456@xuper"

	// initialize a client to operate the contract account
	ca := contractaccount.InitContractAccount(account, node, bcname)

	// create contract account
	txid, err := ca.CreateContractAccount(contractAccount)
	if err != nil {
		log.Printf("CreateContractAccount err: %v", err)
		os.Exit(-1)
	}
	fmt.Println(txid)
	/*
		// the 2nd way to create contract account
		preSelectUTXOResponse, err := ca.PreCreateContractAccount(contractAccount)
		if err != nil {
			log.Printf("PreCreateContractAccount failed, err: %v", err)
			os.Exit(-1)
		}
		txid, err := ca.PostCreateContractAccount(preSelectUTXOResponse)
		if err != nil {
			log.Printf("PostCreateContractAccount failed, err: %v", err)
			os.Exit(-1)
		}
		log.Printf("txid: %v", txid)
	*/
	return
}

func testTransfer() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	to := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount := "10"
	fee := "0"
	desc := ""
	// transfer
	txid, err := trans.Transfer(to, amount, fee, desc)
	if err != nil {
		fmt.Printf("Transfer err: %v\n", err)
	}
	fmt.Printf("transfer tx: %v\n", txid)
	return
}

func testTransferByPlatform() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// retrieve the platform account by mnemonics
	accPlatform, err := account.RetrieveAccount("二 耗 逻 落 燕 死 电 卵 已 浪 教 南", 1)
	if err != nil {
		fmt.Printf("retrieve platform Account err: %v\n", err)
	}
	fmt.Printf("platform account: %v\n", accPlatform)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTransByPlatform(acc, accPlatform, node, bcname)

	// transfer destination address, amount, fee and description
	to := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount := "10"
	fee := "0"
	desc := ""
	// transfer
	txid, err := trans.Transfer(to, amount, fee, desc)
	if err != nil {
		fmt.Printf("Transfer err: %v\n", err)
	}
	fmt.Printf("transfer tx: %v\n", txid)
	return
}

func testEncryptedTransfer() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	to := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount := "10"
	fee := "0"
	desc := "encrypted transfer tx"

	cryptoClient := crypto.GetCryptoClient()
	masterKey, err := cryptoClient.GenerateMasterKeyByMnemonic("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("GenerateMasterKeyByMnemonic err: %v\n", err)
	}

	privateKey, err := cryptoClient.GenerateChildKey(masterKey, hdapi.HardenedKeyStart+1)
	if err != nil {
		fmt.Printf("GenerateChildKey err: %v\n", err)
	}

	publicKey, err := cryptoClient.ConvertPrvKeyToPubKey(privateKey)
	if err != nil {
		fmt.Printf("ConvertPrvKeyToPubKey err: %v\n", err)
	}

	// transfer
	txid, err := trans.EncryptedTransfer(to, amount, fee, desc, publicKey)
	if err != nil {
		fmt.Printf("EncryptedTransfer err: %v\n", err)
	}
	fmt.Printf("EncryptedTransfer tx: %v\n", txid)
	return
}

func testBatchTransfer() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transfer transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// transfer destination address, amount, fee and description
	to1 := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount1 := "10"
	to2 := "jingbo"
	amount2 := "20"

	toAddressAndAmount := make(map[string]string)
	toAddressAndAmount[to1] = amount1
	toAddressAndAmount[to2] = amount2

	fee := "0"
	desc := "multi transfer test"

	// transfer
	txid, err := trans.BatchTransfer(toAddressAndAmount, fee, desc)
	if err != nil {
		fmt.Printf("Transfer err: %v\n", err)
	}
	fmt.Printf("transfer tx: %v\n", txid)
	return
}

func testBatchTransferByPlatform() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// retrieve the platform account by mnemonics
	accPlatform, err := account.RetrieveAccount("二 耗 逻 落 燕 死 电 卵 已 浪 教 南", 1)
	if err != nil {
		fmt.Printf("retrieve platform Account err: %v\n", err)
	}
	fmt.Printf("platform account: %v\n", accPlatform)

	// initialize a client to operate the transfer transaction
	//	trans := transfer.InitTrans(acc, node, bcname)
	trans := transfer.InitTransByPlatform(acc, accPlatform, node, bcname)

	// transfer destination address, amount, fee and description
	to1 := "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv"
	amount1 := "10"
	to2 := "jingbo"
	amount2 := "20"

	toAddressAndAmount := make(map[string]string)
	toAddressAndAmount[to1] = amount1
	toAddressAndAmount[to2] = amount2

	fee := "0"
	desc := "multi transfer test"

	// transfer
	txid, err := trans.BatchTransfer(toAddressAndAmount, fee, desc)
	if err != nil {
		fmt.Printf("Transfer err: %v\n", err)
	}
	fmt.Printf("transfer tx: %v\n", txid)
	return
}

func testCreateChain() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transfer transaction
	chain := network.InitChain(acc, node, bcname)

	// desc for creating a new blockchain

	// ./xchain-cli status -H 127.0.0.1:37801
	// ./xchain-cli account balance dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN -H 127.0.0.1:37801 --name TestChain

	// tdpos的desc demo
	desc := `{
	  "Module": "kernel",
	  "Method": "CreateBlockChain",
	  "Args": {
	    "name": "HelloChain",
	    "data": "{\"maxblocksize\": \"128\", \"award_decay\": {\"height_gap\": 31536000, \"ratio\": 1}, \"version\": \"1\", \"predistribution\": [{\"quota\": \"1000000000000000\", \"address\": \"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\"}], \"decimals\": \"8\", \"period\": \"3000\",\"award\": \"1000000\", \"genesis_consensus\": {\"config\": {\"init_proposer\": {\"1\": [\"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\", \"nYoKRf3jX7vhfSn4jUwHzUf5v5eVxdaNQ\", \"kGXLu6Kex54AJZcp5QPTQ5Hz4ebcUXLLB\"]}, \"timestamp\": \"1534928070000000000\", \"period\": \"500\", \"alternate_interval\": \"3000\", \"term_interval\": \"3000\", \"block_num\": \"10\", \"vote_unit_price\": \"1\", \"proposer_num\": \"3\"}, \"name\": \"tdpos\", \"type\":\"tdpos\"}}"
	    }
	}`

	// single的desc demo
	//	desc := `{
	//    "Module": "kernel",
	//    "Method": "CreateBlockChain",
	//    "Args": {
	//        "name": "TestChain",
	//    	"data": "{\"version\": \"1\", \"consensus\": {\"miner\":\"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\", \"type\":\"single\"},\"predistribution\":[{\"address\": \"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN\",\"quota\": \"1000000000000000\"}],\"maxblocksize\": \"128\",\"period\": \"3000\",\"award\": \"1000000\"}"
	//	    }
	//	}`

	// transfer
	txid, err := chain.CreateChain(desc)
	if err != nil {
		fmt.Printf("create new blockchain err: %v\n", err)
	}
	fmt.Printf("create new blockchain tx: %v\n", txid)
	return
}

func testDeployWasmContract() {
	// retrieve the account by mnemonics
	//	acc, err := account.RetrieveAccount("既 站 冒 装 沈 会 硝 街 储 贡 袁 席", 1)
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)

	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// set contract account, contract will be installed in the contract account
	//	contractAccount := "XC8888888888888888@xuper"
	contractAccount := "XC1234567890123456@xuper"

	// initialize a client to operate the contract
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set init args and contract file
	args := map[string]string{
		"creator": "xchain",
	}
	codePath := "example/contract_code/counter.wasm"

	// deploy wasm contract
	txid, err := wasmContract.DeployWasmContract(args, codePath, "c")
	if err != nil {
		log.Printf("DeployWasmContract err: %v", err)
		os.Exit(-1)
	}
	fmt.Printf("DeployWasmContract txid: %v\n", txid)

	/*
		// the 2nd way to deploy wasm contract, preDeploy and Post
		preSelectUTXOResponse, err := wasmContract.PreDeployWasmContract(args, codePath, "c")
		if err != nil {
			log.Printf("DeployWasmContract GetPreDeployWasmContractRes failed, err: %v", err)
			os.Exit(-1)
		}
		txid, err := wasmContract.PostWasmContract(preSelectUTXOResponse)
		if err != nil {
			log.Printf("DeployWasmContract PostWasmContract failed, err: %v", err)
			os.Exit(-1)
		}
		log.Printf("txid: %v", txid)
	*/
	return
}

func testInvokeWasmContract() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the contract
	contractAccount := "XC1234567890123456@xuper"
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set invoke function method and args
	args := map[string]string{
		"key": "counter",
	}
	methodName := "increase"

	// invoke contract
	txid, err := wasmContract.InvokeWasmContract(methodName, args)
	if err != nil {
		log.Printf("InvokeWasmContract PostWasmContract failed, err: %v", err)
		os.Exit(-1)
	}
	log.Printf("txid: %v", txid)
	/*
		// the 2nd way to invoke wasm contract, preInvoke and Post
		preSelectUTXOResponse, err := wasmContract.PreInvokeWasmContract(methodName, args)
		if err != nil {
			log.Printf("InvokeWasmContract GetPreMethodWasmContractRes failed, err: %v", err)
			os.Exit(-1)
		}
		txid, err := wasmContract.PostWasmContract(preSelectUTXOResponse)
		if err != nil {
			log.Printf("InvokeWasmContract PostWasmContract failed, err: %v", err)
			os.Exit(-1)
		}
		log.Printf("txid: %v", txid)
	*/
	return
}

func testQueryWasmContract() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("既 站 冒 装 沈 会 硝 街 储 贡 袁 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the contract
	contractAccount := ""
	wasmContract := contract.InitWasmContract(acc, node, bcname, contractName, contractAccount)

	// set query function method and args
	args := map[string]string{
		"key": "counter",
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
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// get balance of the account
	balance, err := trans.GetBalance()
	log.Printf("balance %v, err %v", balance, err)
	return
}

func testGetMultiChainBalance() {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
		return
	}
	fmt.Printf("account: %v\n", acc)

	bcNames := []string{}
	bcNames = append(bcNames, "xuper")
	bcNames = append(bcNames, "HelloChain8")

	// initialize a client to operate the transaction
	balanceUtil := balance.InitBalance(acc, node, bcNames)

	// get balance of the account
	balances, err := balanceUtil.GetBalanceDetails()
	log.Printf("balances %v, err %v", balances, err)
	return
}

func testQueryTx() {
	// initialize a client to operate the transaction
	trans := transfer.InitTrans(nil, node, bcname)
	txid := "3a78d06dd39b814af113dbdc15239e675846ec927106d50153665c273f51001e"

	// query tx by txid
	tx, err := trans.QueryTx(txid)
	log.Printf("query tx %v, err %v", tx, err)
	return
}

func testDecryptedTx() {
	// initialize a client to operate the transaction
	trans := transfer.InitTrans(nil, node, bcname)
	//	txid := "b59a83d9ade65ef2e0e50bfbcb497c6310c527a59f1f4b2ba66a24518b43cd03"
	txid := "4cf794fe7de9a497147859019fdb01a1c9e09a2abf7d5afd0265604eee8517ca"

	// query tx by txid
	TxStatus, err := trans.QueryTx(txid)
	if err != nil {
		fmt.Printf("QueryTx err: %v\n", err)
	}
	encryptedTx := TxStatus.Tx

	cryptoClient := crypto.GetCryptoClient()
	masterKey, err := cryptoClient.GenerateMasterKeyByMnemonic("浓 玉 寿 元 杭 仗 堡 呈 昨 逐 抚 席", 1)
	if err != nil {
		fmt.Printf("GenerateMasterKeyByMnemonic err: %v\n", err)
	}

	decryptedDesc, err := trans.DecryptedTx(encryptedTx, masterKey)
	log.Printf("decrypted tx desc [%v], err %v", decryptedDesc, err)
	return
}

func main() {
	testAccount()
	testTransfer()
	//	testTransferByPlatform()
	// TODO 广播交易，同时对desc使用分层确定性技术对desc进行加密
	//	testEncryptedTransfer()
	//	testBatchTransfer()
	//	testBatchTransferByPlatform()
	//	testContractAccount()
	//	testDeployWasmContract()
	//	testInvokeWasmContract()
	//	testQueryWasmContract()
	testGetBalance()
	//	testGetMultiChainBalance()
	//	testQueryTx()
	// TODO 查询交易，同时对desc使用分层确定性技术对desc进行加密
	//	testDecryptedTx()
	//	testCreateChain()

	return
}
