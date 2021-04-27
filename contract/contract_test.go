// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
)

var wasmNode = "127.0.0.1:37101"

func TestDeployWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %+v", acc)
	t.Logf("public: %+v", acc.PublicKey)

	var testDeployWasmContracts = []struct {
		account         *account.Account
		bcname          string
		contractName    string
		contractAccount string
	}{
		{
			account:         acc,
			bcname:          "xuper",
			contractName:    "counterwasm5",
			contractAccount: "XC2222222222222222@xuper",
		},
	}
	sdkClient, err := xchain.NewXuperClient(wasmNode)
	if err != nil {
		t.Error(err)
	}

	args := map[string]string{
		"creator": "xchain",
	}
	//codePath := "../example/contract_code/counter.wasm"
	codePath := "./counter2.wasm"
	runtime := "c"
	for _, arg := range testDeployWasmContracts {
		wasmContract := InitWasmContractWithClient(arg.account, arg.bcname, arg.contractName, arg.contractAccount, sdkClient)
		txid, err := wasmContract.DeployWasmContract(args, codePath, runtime)
		t.Logf("DeployWasmContract txid: %v, err: %v", txid, err)
	}
}

func TestInvokeWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", acc)

	var testInvokeWasmContracts = []struct {
		account         *account.Account
		bcname          string
		contractName    string
		contractAccount string
	}{
		{
			account:         acc,
			bcname:          "xuper",
			contractName:    "counterwasm2",
			contractAccount: "XC2222222222222222@xuper",
		},
	}
	sdkClient, err := xchain.NewXuperClient(wasmNode)
	if err != nil {
		t.Error(err)
	}
	args := map[string]string{
		"key": "xchain",
	}
	for _, arg := range testInvokeWasmContracts {
		for _, method := range []string{"increase", "get"} {
			wasmContract := InitWasmContractWithClient(arg.account, arg.bcname, arg.contractName, arg.contractAccount, sdkClient)
			txid, err := wasmContract.InvokeWasmContract(method, args)
			t.Logf("InvokeWasmContract txid: %v, err: %v", txid, err)
		}
	}
}
