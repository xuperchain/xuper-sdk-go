// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contractaccount is related to contract account operation
package contractaccount

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
)

// test CreateContractAccount
func TestCreateContractAccount(t *testing.T) {
	testAccount, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var testContractAccounts = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "XC123456789012345@xuper",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37101",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
	}

	for _, testContractAccount := range testContractAccounts {
		initContractAccount := InitContractAccount(testContractAccount.account, testContractAccount.node,
			testContractAccount.bcname)
		if initContractAccount == nil {
			t.Fatal("initContractAccount error")
		}

		contractAccount, err := initContractAccount.CreateContractAccount(testContractAccount.contractAccount)
		t.Logf("create contract account, res: %v, err: %v", contractAccount, err)
	}
}

// test PreCreateContractAccount and PostCreateContractAccount
func TestCreateContractAccount2(t *testing.T) {
	testAccount, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var testContractAccounts = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "123456789012345",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37101",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
	}

	for _, testContractAccount := range testContractAccounts {
		initContractAccount := InitContractAccount(testContractAccount.account, testContractAccount.node,
			testContractAccount.bcname)
		if initContractAccount == nil {
			t.Fatal("initContractAccount error")
		}

		preExeResp, err := initContractAccount.PreCreateContractAccount(testContractAccount.contractAccount)
		if err != nil {
			t.Logf("preExe error:%v", err)
		} else {
			txid, err := initContractAccount.PostCreateContractAccount(preExeResp)
			t.Logf("create contract account, res: %v, err: %v", txid, err)
		}
	}
}

// test PreCreateContractAccount and PostCreateContractAccount
func TestCreateContractAccount3(t *testing.T) {
	testAccount, err := account.RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var testContractAccounts = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "XC123456789012345@xuper",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37101",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
		{
			account:         testAccount,
			node:            "127.0.0.1:37201",
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
	}

	for _, testContractAccount := range testContractAccounts {
		initContractAccount := InitContractAccount(testContractAccount.account, testContractAccount.node,
			testContractAccount.bcname)
		if initContractAccount == nil {
			t.Fatal("initContractAccount error")
		}
		preExeResp, err := initContractAccount.PreCreateContractAccount(testContractAccount.contractAccount)
		if err != nil {
			t.Logf("preExe error:%v", err)
		} else {

			initContractAccount = InitContractAccount(testContractAccount.account, testContractAccount.node,
				testContractAccount.bcname)
			if initContractAccount == nil {
				t.Fatal("initContractAccount error")
			}
			txid, err := initContractAccount.PostCreateContractAccount(preExeResp)
			t.Logf("create contract account, res: %v, err: %v", txid, err)
		}
	}
}
