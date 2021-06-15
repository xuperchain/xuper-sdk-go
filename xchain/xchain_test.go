// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package xchain is related to xchain operation
package xchain

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"testing"
)

func initXchain() *Xchain {
	commConfig := config.GetInstance()
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RetrieveAccount: %v\n", acc)

	return &Xchain{
		Cfg:       commConfig,
		XchainSer: "127.0.0.1:37101",
		ChainName: "xuper",
		Account:   acc,
	}
}

func TestNewXuperClient(t *testing.T) {
	node := "127.0.0.1:37101"
	xuperClient, err := NewXuperClient(node)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v", xuperClient)
}

func initXchainWithClient() *Xchain {
	node := "127.0.0.1:37101"
	xuperClient, err := NewXuperClient(node)
	if err != nil {
		panic(err)
	}
	commConfig := config.GetInstance()

	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RetrieveAccount: %v\n", acc)

	return &Xchain{
		Cfg:         commConfig,
		ChainName:   "xuper",
		Account:     acc,
		XuperClient: xuperClient,
	}
}

func TestXchain_QueryBlockByHeight(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.QueryBlockByHeight(12)
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("block is nil")
	}
}

func TestXchain_GetAccountByAk(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.GetAccountByAk(qc.Account.Address)
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("account is nil")
	}
}

func TestXchain_GetAccountContracts(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.GetAccountContracts("XC2222222222222222@xuper")
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("contracts is nil")
	}
}

func TestXchain_QueryUTXORecord(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.QueryUTXORecord(qc.Account.Address, 1)
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("record is nil")
	}
}

func TestXchain_QueryContractMethondAcl(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.QueryContractMethondAcl("golangcounter5", "Increase")
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("method is nil")
	}
}

func TestQueryTx(t *testing.T) {
	qc := initXchainWithClient()
	b, err := qc.QueryTx("c7545ca8c6f604aa9eec972c64ff9c098dcde86a08decb4424c187e5200d4f17")
	if err != nil {
		t.Error(err)
	}
	if b != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("tx is nil")
	}

	b2, err := qc.QueryTx("c7545ca8c6f604aa9eec972c64ff9c098dcde86a08decb4424c187e5200d4f19")
	if b2 != nil {
		fmt.Printf("%+v\n", b)
	} else {
		fmt.Println("tx is nil")
	}
	err = qc.CloseClient()
	if err != nil {
		t.Error(err)
	}
}
