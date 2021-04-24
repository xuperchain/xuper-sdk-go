// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
)

//
func TestTransfer(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	node := "127.0.0.1:37101"
	bcname := "xuper"

	sdkClient,err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}

	trans := InitTrans(acc, bcname, sdkClient)

	testCase := []struct {
		to     string
		amount string
		fee    string
		desc   string
	}{
		{
			to:     "jRGSGzpkWLcVBhxbLbdKLuc2drW55kLsf",
			amount: "2",
			fee:    "0",
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
	}{
		{
			account: acc,
		},
		{
			account: nil,
		},
		{
			account: &account.Account{Address:"W7UuhkbGXbCx4BrFamTZaK96QN6rAbREk"},
		},
		{
			account: &account.Account{Address:"XC1111111111111111@xuper"},
		},
	}

	node := "127.0.0.1:37101"
	bcname := "xuper"
	sdkClient,err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	for _, arg := range testCase {
		trans := InitTrans(arg.account,bcname,sdkClient)
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

	node := "127.0.0.1:37101"
	bcname := "xuper"
	sdkClient,err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	trans := InitTrans(acc, bcname,sdkClient)

	testCase := []struct {
		txid string
	}{
		{
			txid: "fb225328af683506b36e9b5f8b389e3c4c4e8759bafe5330f0aca9b753183536",
		},
		{
			txid: "",
		},
		{
			txid: "264afe6c55000277e9449ab2711dc53bb3754cf9125587adba1c7db8afc7eec8",
		},
	}

	for _, arg := range testCase {
		tx, err := trans.QueryTx(arg.txid)
		t.Logf("Querytx tx: %v, err: %v", tx, err)
	}
}
