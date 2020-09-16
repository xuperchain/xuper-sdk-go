// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func TestDeployWasmContract(t *testing.T) {
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
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
			node:            "127.0.0.1:37201",
			contractName:    "countermmm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "127.0.0.1:37201",
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
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
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
			node:            "127.0.0.1:37201",
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "",
			node:            "127.0.0.1:37201",
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
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
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
			node:            "127.0.0.1:37201",
			contractName:    "counterm.wasm",
			contractAccount: "XC8888888888888888@xuper",
		},
		{
			account:         acc,
			bcname:          "xuper",
			node:            "127.0.0.1:37201",
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
