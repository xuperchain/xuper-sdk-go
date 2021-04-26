// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"encoding/json"
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"io/ioutil"
	"sync"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
)

var (
	node   = "127.0.0.1:37101"
	bcname = "xuper"
)

func TestTransfer(t *testing.T) {
	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	sdkClient, err := xchain.NewSDKClient(node)
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
			to:     "UgdxaYwTzUjkvQnmeB3VgnGFdXfrsiwFv",
			amount: "200",
			fee:    "0",
			desc:   "",
		},
	}

	for _, arg := range testCase {
		tx, err := trans.Transfer(arg.to, arg.amount, arg.fee, arg.desc)
		t.Logf("transfer tx: %v, err: %v", tx, err)
	}
}

func TestGoroutine(t *testing.T) { //并发执行
	accList1 := []*account.Account{}
	accList2 := []*account.Account{}
	bz1, err := ioutil.ReadFile("./wallet1")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bz1, &accList1)
	if err != nil {
		t.Error(err)
	}

	bz2, err := ioutil.ReadFile("./wallet2")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bz2, &accList2)
	if err != nil {
		t.Error(err)
	}

	//accFrom, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Logf("RetrieveAccount: %v\n", accFrom)

	sdkClient, err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(accList1))

	for i := 0; i < len(accList1); i++ {
		trans := InitTrans(accList1[i], bcname, sdkClient)
		go func(index int) {
			txid, err := trans.transfer(accList2[index].Address, "1", "0", "", "")
			if err != nil {
				t.Error(err)
			}
			fmt.Println("交易成功：", txid)
			wg.Done()
		}(i)
	}
	fmt.Println("等待中。。。")
	wg.Wait()
	fmt.Println("任务完成")
}

func TestCreateAddress(t *testing.T) {
	accList := []*account.Account{}
	var language int = 1
	times := 100

	for i := 0; i < times; i++ {
		acc, err := account.CreateAccount(uint8(1), language)
		if err != nil {
			t.Error(err)
		}
		fmt.Printf("%+v\n", acc)
		accList = append(accList, acc)
	}
	bz, err := json.Marshal(accList)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("./wallet2", bz, 777)
	if err != nil {
		t.Error(err)
	}
}

func TestSendTokenToWallet(t *testing.T) {
	bz, err := ioutil.ReadFile("./wallet2")
	if err != nil {
		t.Error(err)
	}
	accList := []*account.Account{}
	err = json.Unmarshal(bz, &accList)
	if err != nil {
		t.Error(err)
	}

	accFrom, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", accFrom)

	sdkClient, err := xchain.NewSDKClient(node)
	if err != nil {
		t.Error(err)
	}
	trans := InitTrans(accFrom, bcname, sdkClient)

	for i := 0; i < len(accList); i++ {
		tx, err := trans.Transfer(accList[i].Address, "100", "", "")
		if err != nil {
			panic(err)
		}
		fmt.Println("执行成功：", tx)
		time.Sleep(time.Second * 2)

	}

}

//func TestGetBalace(t *testing.T) {
//	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("RetrieveAccount: %v\n", acc)
//
//	testCase := []struct {
//		account *account.Account
//		node    string
//		bcname  string
//	}{
//		{
//			account: acc,
//			node:    "127.0.0.1:37201",
//			bcname:  "xuper",
//		},
//		{
//			account: nil,
//			node:    "127.0.0.1:37201",
//			bcname:  "xuper",
//		},
//		{
//			account: acc,
//			node:    "127.0.0.1:37201",
//			bcname:  "",
//		},
//		{
//			account: acc,
//			node:    "",
//			bcname:  "",
//		},
//	}
//
//	for _, arg := range testCase {
//		trans := InitTrans(arg.account, arg.node, arg.bcname)
//		balance, err := trans.GetBalance()
//		t.Logf("get balance: %v, err: %v", balance, err)
//	}
//}

//func TestQueryTx(t *testing.T) {
//	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
//	if err != nil {
//		t.Fatal(err)
//	}
//	t.Logf("RetrieveAccount: %v\n", acc)
//
//	node := "127.0.0.1:37201"
//	bcname := "xuper"
//	trans := InitTrans(acc, node, bcname)
//
//	testCase := []struct {
//		txid string
//	}{
//		{
//			txid: "3a78d06dd39b814af113dbdc15239e675846ec927106d50153665c273f51001e",
//		},
//		{
//			txid: "",
//		},
//		{
//			txid: "fdsfdsa",
//		},
//	}
//
//	for _, arg := range testCase {
//		tx, err := trans.QueryTx(arg.txid)
//		t.Logf("Querytx tx: %v, err: %v", tx, err)
//	}
//}
