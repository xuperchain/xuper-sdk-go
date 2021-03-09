package contract

import (
	"fmt"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func Example_wasm() {
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	node := "127.0.0.1:37101"
	bcName := "xuper"
	contractName := "counter"
	contractAccount := "XC8888888888888888@xuper"
	wasmContract := InitWasmContract(acc, node, bcName, contractName, contractAccount)

	args := map[string]string{
		"key": "counter",
	}
	txid, err := wasmContract.DeployWasmContract(args, "./counter.wasm", "c")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("部署合约交易 ID：", txid)

	method := "increase"
	txid, err = wasmContract.InvokeWasmContract(method, args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("调用合约交易 ID：", txid)
}

func Example_evm() {

	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	node := "127.0.0.1:37101"
	bcName := "xuper"
	contractName := "counter"
	contractAccount := "XC8888888888888888@xuper"
	evmContract := InitEVMContract(acc, node, bcName, contractName, contractAccount)

	args := map[string]string{
		"num": "1",
	}

	abi := `[{"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"constant":false,"inputs":[],"name":"retrieve","outputs":[{"internalType":"uint256","name":"","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"store","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"internalType":"uint256","name":"num","type":"uint256"}],"name":"storepay","outputs":[],"payable":true,"stateMutability":"payable","type":"function"}]`
	bin := `608060405234801561001057600080fd5b5060405161016c38038061016c8339818101604052602081101561003357600080fd5b810190808051906020019092919050505080600081905550506101118061005b6000396000f3fe60806040526004361060305760003560e01c80632e64cec11460355780636057361d14605d5780638995db74146094575b600080fd5b348015604057600080fd5b50604760bf565b6040518082815260200191505060405180910390f35b348015606857600080fd5b50609260048036036020811015607d57600080fd5b810190808035906020019092919050505060c8565b005b60bd6004803603602081101560a857600080fd5b810190808035906020019092919050505060d2565b005b60008054905090565b8060008190555050565b806000819055505056fea265627a7a723158209500c3e12321b837819442c0bc1daa92a4f4377fc7b59c41dbf9c7620b2f961064736f6c63430005110032`

	txid, err := evmContract.Deploy(args, []byte(bin), []byte(abi))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	args = map[string]string{
		"num": "2",
	}
	txid, err = evmContract.Invoke("store", args, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println(txid)

	preExeRPCRes, err := evmContract.Query("retrieve", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	gas := preExeRPCRes.GetResponse().GetGasUsed()
	fmt.Printf("gas used: %v\n", gas)
	fmt.Printf("preExeRPCRes: %v \n", preExeRPCRes)
	for _, res := range preExeRPCRes.GetResponse().GetResponse() {
		fmt.Printf("contract response: %s\n", string(res))
	}
}
