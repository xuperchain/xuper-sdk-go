package query

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
)

var node = "127.0.0.1:37101"

func TestQueryClient_QueryBlockByHeight(t *testing.T) {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, xuperClient)
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

func TestQueryClient_GetAccountByAk(t *testing.T) {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, xuperClient)
	resp, err := qc.GetAccountByAk(acc.Address)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestQueryClient_GetAccountContracts(t *testing.T) {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, xuperClient)
	resp, err := qc.GetAccountContracts("XC2222222222222222@xuper")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestQueryClient_QueryUTXORecord(t *testing.T) {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, xuperClient)
	resp, err := qc.QueryUTXORecord(acc.Address, 1)
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", resp)
}

func TestQueryClient_QueryContractMethondAcl(t *testing.T) {
	xuperClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, xuperClient)
	resp, err := qc.QueryContractMethondAcl("golangcounter5", "Increase")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%+v\n", resp)

}
