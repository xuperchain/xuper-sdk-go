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

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	trans := InitTransWithClient(acc, bcname, sdkClient)

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

//测试并发
func TestGoroutine(t *testing.T) { //并发执行
	accList1 := []*account.Account{}
	accList2 := []*account.Account{}
	bz1, err := ioutil.ReadFile("./wallet3")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bz1, &accList1)
	if err != nil {
		t.Error(err)
	}

	bz2, err := ioutil.ReadFile("./wallet4")
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

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(accList1))

	for i := 0; i < len(accList1); i++ {
		trans := InitTransWithClient(accList1[i], bcname, sdkClient)
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

func TestGoroutine1(t *testing.T) { //并发执行
	accList1 := []*account.Account{}
	accList2 := []*account.Account{}
	bz1, err := ioutil.ReadFile("./wallet2")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bz1, &accList1)
	if err != nil {
		t.Error(err)
	}

	bz2, err := ioutil.ReadFile("./wallet3")
	if err != nil {
		t.Error(err)
	}
	err = json.Unmarshal(bz2, &accList2)
	if err != nil {
		t.Error(err)
	}
	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	wg := sync.WaitGroup{}
	wg.Add(len(accList2))

	for j := 0; j < 100; j++ {
		for i := 0; i < len(accList1); i++ {
			trans := InitTransWithClient(accList1[i], bcname, sdkClient)
			go func(index int) {
				txid, err := trans.transfer(accList2[index].Address, "100", "0", "", "")
				if err != nil {
					t.Error(err)
				}
				fmt.Println("交易成功：", txid)
				wg.Done()
			}(j*100 + i)
		}
		time.Sleep(time.Second * 2)
	}
	fmt.Println("等待中。。。")
	wg.Wait()
	fmt.Println("任务完成")
}

// 辅助并发测试，创建多个地址用来测试并发
func TestCreateAddress(t *testing.T) {
	accList, err := createAddress(10000)
	if err != nil {
		t.Error(err)
	}

	bz, err := json.Marshal(accList)
	if err != nil {
		t.Error(err)
	}
	err = ioutil.WriteFile("./wallet4", bz, 777)
	if err != nil {
		t.Error(err)
	}
}

func createAddress(n int) ([]*account.Account, error) {
	accList := []*account.Account{}
	var language int = 1

	for i := 0; i < n; i++ {
		acc, err := account.CreateAccount(uint8(1), language)
		if err != nil {
			return nil, err
		}
		fmt.Printf("%+v\n", acc)
		accList = append(accList, acc)
	}
	return accList, nil
}

// 辅助并发测试，用来给测试并发的from钱包转账进行转账
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

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	trans := InitTransWithClient(accFrom, bcname, sdkClient)

	for i := 0; i < len(accList); i++ {
		tx, err := trans.Transfer(accList[i].Address, "100000", "", "")
		if err != nil {
			panic(err)
		}
		fmt.Println("执行成功：", tx)
		time.Sleep(time.Second * 1)
	}
}

// 辅助并发测试，输出钱包中文件长度
func TestLength(t *testing.T) { //并发执行
	accList1 := []*account.Account{}
	accList2 := []*account.Account{}
	accList3 := []*account.Account{}
	accList4 := []*account.Account{}
	bz1, err := ioutil.ReadFile("./wallet21")
	if err != nil {
		t.Error(err)
	}

	bz2, err := ioutil.ReadFile("./wallet2")
	if err != nil {
		t.Error(err)
	}

	bz3, err := ioutil.ReadFile("./wallet3")
	if err != nil {
		t.Error(err)
	}

	bz4, err := ioutil.ReadFile("./wallet4")
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bz1, &accList1)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bz2, &accList2)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bz3, &accList3)
	if err != nil {
		t.Error(err)
	}

	err = json.Unmarshal(bz4, &accList4)
	if err != nil {
		t.Error(err)
	}

	fmt.Printf("length:\nwallet1:%d\nwallet2:%d\nwallet3:%d\nwallet4:%d\n", len(accList1), len(accList2), len(accList3), len(accList4))
}

//测试要点，1，交易要能成功 ，不能panic, 2 发送方的确没有付fee
func TestNoFee(t *testing.T) {
	accList, err := createAddress(2)
	if err != nil {
		t.Error(err)
	}
	accFrom, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", accFrom)

	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	trans := InitTransWithClient(accFrom, bcname, sdkClient)
	txId1, err := trans.transfer(accList[0].Address, "100", "", "", "")
	if err != nil {
		t.Error(err)
	}
	acc1Client := InitTransWithClient(accList[0], bcname, sdkClient)
	acc2Client := InitTransWithClient(accList[1], bcname, sdkClient)
	balance1, err := acc1Client.GetBalance()
	if err != nil {
		t.Error()
	}
	balance2, err := acc2Client.GetBalance()
	if err != nil {
		t.Error()
	}
	fmt.Printf("balances account1: %s, account2: %s\n", balance1, balance2)
	txid2, err := acc1Client.transfer(accList[1].Address, "1", "", "", "")
	if err != nil {
		t.Error(err)
	}

	balance1, err = acc1Client.GetBalance()
	if err != nil {
		t.Error()
	}
	balance2, err = acc2Client.GetBalance()
	if err != nil {
		t.Error()
	}
	fmt.Printf("after nofee transfer balances \naccount1: %s, account2: %s\n", balance1, balance2)
	fmt.Printf("tx1:%s\ntx2:%s\n", txId1, txid2)
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
//		trans := InitTransWithClient(arg.account, arg.node, arg.bcname)
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
//	trans := InitTransWithClient(acc, node, bcname)
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
