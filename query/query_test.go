package query

import (
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"testing"
)

var node = "127.0.0.1:37101"

func TestQueryClient_QueryBlockByHeight(t *testing.T) {
	sdkClient, err := xchain.NewXuperClient(node)
	if err != nil {
		t.Errorf("New sdk error")
	}

	acc, err := account.RetrieveAccount("售 历 定 栽 护 沟 万 城 发 阵 凶 据", 1)
	if err != nil {
		t.Fatal(err)
	}
	chainName := "xuper"

	qc := InitClientWithClient(acc, chainName, sdkClient)
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
