// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"sync"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
)

const (
	node     = "10.13.32.249:37101"
	mnemonic = "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
)

func init() {
	config.SetConfig(node, "", "", "", false, false, "")
}

func TestDeployWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", acc)

	var testDeployWasmContracts = []struct {
		account         *account.Account
		bcname          string
		node            string
		contractName    string
		contractAccount string
	}{
		{
			account:         acc,
			bcname:          "xuper",
			node:            node,
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            node,
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "",
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "",
			contractName:    "",
			contractAccount: "XC8888888888888888@xuper",
		},
	}

	args := map[string]string{
		"creator": "xchain",
	}
	codePath := "../example/contract_code/counter.wasm"
	runtime := "c"
	for _, arg := range testDeployWasmContracts {
		wasmContract := InitWasmContract(arg.account, arg.node, arg.bcname, arg.contractName, arg.contractAccount)
		txid, err := wasmContract.DeployWasmContract(args, codePath, runtime)
		t.Logf("DeployWasmContract txid: %v, err: %v", txid, err)
	}
}

func TestInvokeWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", acc)

	var testInvokeWasmContracts = []struct {
		account         *account.Account
		bcname          string
		node            string
		contractName    string
		contractAccount string
	}{
		{
			account:         acc,
			bcname:          "xuper",
			node:            node,
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            node,
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "",
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "",
			contractName:    "",
			contractAccount: "XC8888888888888888@xuper",
		},
	}

	args := map[string]string{
		"key": "counter",
	}
	for _, arg := range testInvokeWasmContracts {
		for _, method := range []string{"increase", "get", "query"} {
			wasmContract := InitWasmContract(arg.account, arg.node, arg.bcname, arg.contractName, arg.contractAccount)
			txid, err := wasmContract.InvokeWasmContract(method, args)
			t.Logf("InvokeWasmContract txid: %v, err: %v", txid, err)
		}
	}
}

func TestUpgradeWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", acc)

	var tests = []struct {
		account         *account.Account
		bcname          string
		node            string
		contractName    string
		contractAccount string
	}{
		{
			account:         acc,
			bcname:          "xuper",
			node:            node,
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "xuper",
			node:            node,
			contractName:    "countermmm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
	}

	args := map[string]string{
		"creator": "xchain",
	}
	codePath := "../example/contract_code/counter.wasm"
	runtime := "c"
	for _, arg := range tests {
		wasmContract := InitWasmContract(arg.account, arg.node, arg.bcname, arg.contractName, arg.contractAccount)
		txid, err := wasmContract.UpgradeWasmContract(args, codePath, runtime)
		t.Logf("UpgradeWasmContract txid: %v, err: %v", txid, err)
	}
}

func TestInvokeWasmContract2(t *testing.T) {
	acc, _ := account.RetrieveAccount(mnemonic, 1)
	wasmContract := InitWasmContract(acc, node, "xuper", "counter", "XC1111111111111111@xuper")

	wg := &sync.WaitGroup{}
	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			txid, err := wasmContract.InvokeWasmContract("increase", map[string]string{"key": "counter"})
			if err != nil {
				t.Log("err:", err)
				return
			}
			t.Log("txid:", txid)
		}()
	}
	wg.Wait()
}
