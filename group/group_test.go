package group

import (
	"testing"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
)

const (
	node     = "10.13.32.249:37101"
	mnemonic = "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
	bcname   = "xuper"
)

func init() {
	config.SetConfig(node, "", "", "", false, false, "")
}

func TestGroup(t *testing.T) {
	acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", acc)

	var tests = []struct {
		account *account.Account
		bcname  string
		neturl  string
		address string
	}{
		{
			account: acc,
			bcname:  "xuper",
			neturl:  "/ip4/127.0.0.1/tcp/47101/p2p/QmVxeNubpg1ZQjQT8W5yZC9fD7ZB1ViArwvyGUB53sqf8e",
			address: "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN",
		},
		{
			account: acc,
			bcname:  "hello",
			neturl:  "/ip4/127.0.0.1/tcp/47101/p2p/QmVxeNubpg1ZQjQT8W5yZC9fD7ZB1ViArwvyGUB53sqf8e",
			address: "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN",
		},
	}

	wasm := InitGroup(acc, node, bcname)
	for _, test := range tests {

		res, err := wasm.AddChain(test.bcname)
		t.Logf("add res: %v, err: %v", res, err)
		res, err = wasm.DelChain(test.bcname)
		t.Logf("del res: %v, err: %v", res, err)
		resp, err := wasm.ListChain()
		t.Logf("list res: %v, err: %v", resp, err)

		res, err = wasm.AddNode(test.bcname, test.neturl, test.address)
		t.Logf("add res: %v, err: %v", res, err)
		res, err = wasm.DelNode(test.bcname, test.neturl)
		t.Logf("del res: %v, err: %v", res, err)
		resp, err = wasm.ListNode(test.bcname)
		t.Logf("list res: %v, err: %v", resp, err)

	}
}
