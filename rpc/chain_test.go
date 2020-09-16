package rpc

import (
	"fmt"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/util"
)

const (
	node     = "10.13.32.249:37101"
	bcname   = "xuper"
	mnemonic = "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
)

func TestQueryStatus(t *testing.T) {
	reply, err := QueryStatus(node)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)

	reply, err = QueryStatus(node, bcname)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryBlockByHeight(t *testing.T) {
	reply, err := QueryBlockByHeight(node, bcname, 0)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryBlockById(t *testing.T) {
	id := "8df6544cc5242f4eedede3b50c58b06941b300c499812daaf93ab7a6b366a433"
	reply, err := QueryBlockById(node, bcname, id)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryTxById(t *testing.T) {
	id := "9e3fb90dd6b69b063ab3ed649a13b4d65f0f9b41324137eec96187120aa97ad7"
	reply, err := QueryTxById(node, bcname, id)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryAccountAcl(t *testing.T) {
	account := "XC1234567890123456@xuper"
	reply, err := QueryAccountAcl(node, bcname, account)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryMethodAcl(t *testing.T) {
	contract := "counter"
	method := "increase"
	reply, err := QueryMethodAcl(node, bcname, contract, method)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryBalance(t *testing.T) {
	account := "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN"
	reply, err := QueryBalance(node, bcname, account)
	if err != nil {
		t.Fatal(err)
	}
	util.Print(reply)
}

func TestQueryNetUrl(t *testing.T) {
	reply, err := QueryNetUrl(node)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(reply)
}
