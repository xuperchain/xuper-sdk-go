// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func TestTransfer(t *testing.T) {
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	node := "127.0.0.1:37801"
	bcname := "xuper"
	trans := InitTrans(acc, node, bcname)

	testCase := []struct {
		to     string
		amount string
		fee    string
		desc   string
	}{
		{
			to:     "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv",
			amount: "200",
			fee:    "0",
			desc:   "",
		},
		{
			to:     "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv",
			amount: "",
			fee:    "",
			desc:   "",
		},
		{
			to:     "",
			amount: "",
			fee:    "",
			desc:   "",
		},
		{
			to:     "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv",
			amount: "10",
			fee:    "",
			desc:   "",
		},
		{
			to:     "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv",
			amount: "-10",
			fee:    "-3",
			desc:   "",
		},
	}

	for _, arg := range testCase {
		tx, err := trans.Transfer(arg.to, arg.amount, arg.fee, arg.desc)
		t.Logf("transfer tx: %v, err: %v", tx, err)
	}
}

func TestGetBalace(t *testing.T) {
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	testCase := []struct {
		account *account.Account
		node    string
		bcname  string
	}{
		{
			account: acc,
			node:    "127.0.0.1:37201",
			bcname:  "xuper",
		},
		{
			account: nil,
			node:    "127.0.0.1:37201",
			bcname:  "xuper",
		},
		{
			account: acc,
			node:    "127.0.0.1:37201",
			bcname:  "",
		},
		{
			account: acc,
			node:    "",
			bcname:  "",
		},
	}

	for _, arg := range testCase {
		trans := InitTrans(arg.account, arg.node, arg.bcname)
		balance, err := trans.GetBalance()
		t.Logf("get balance: %v, err: %v", balance, err)
	}
}

func TestQueryTx(t *testing.T) {
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	node := "127.0.0.1:37201"
	bcname := "xuper"
	trans := InitTrans(acc, node, bcname)

	testCase := []struct {
		txid string
	}{
		{
			txid: "3a78d06dd39b814af113dbdc15239e675846ec927106d50153665c273f51001e",
		},
		{
			txid: "",
		},
		{
			txid: "fdsfdsa",
		},
	}

	for _, arg := range testCase {
		tx, err := trans.QueryTx(arg.txid)
		t.Logf("Querytx tx: %v, err: %v", tx, err)
	}
}
