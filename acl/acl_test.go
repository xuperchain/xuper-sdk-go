package acl

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
)

const (
	node     = "10.13.32.249:37101"
	mnemonic = "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"
)

func init() {
	config.SetConfig(node, "", "", "", false, false, "")
}

// test CreateContractAccount
func TestCreateContractAccount(t *testing.T) {
	testAccount, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var tests = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC123456789012345@xuper",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC" + (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16] + "@xuper",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "1234567890123456",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "123456789012345a",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "",
		},
	}

	for _, test := range tests {
		init := InitAcl(
			test.account,
			test.node,
			test.bcname,
		)
		if init == nil {
			t.Fatal("InitACL error")
		}

		contractAccount, err := init.CreateContractAccount(test.contractAccount)
		t.Logf("create contract account, res: %v, err: %v", contractAccount, err)
	}
}

// test SetContractAccountAcl
func TestSetContractAccountAcl(t *testing.T) {
	testAccount, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var tests = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "1234567812345678",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC1234567890123456@xuper",
		},
	}

	for _, test := range tests {
		init := InitAcl(
			test.account,
			test.node,
			test.bcname,
		)
		if init == nil {
			t.Fatal("InitACL error")
		}

		acls, err := init.QueryAccountAcl(test.contractAccount)
		t.Logf("query contract account acl, res: %v, err: %v", acls.GetAcl(), err)

		addrs := []string{"nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh", "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpNr"}
		res, err := init.SetContractAccountAcl(test.contractAccount, addrs)
		t.Logf("set contract account acl, res: %v, err: %v", res, err)

		acls, err = init.QueryAccountAcl(test.contractAccount)
		t.Logf("query contract account acl, res: %v, err: %v", acls.GetAcl(), err)
	}
}

// test SetContractMethodAcl
func TestSetContractMethodAcl(t *testing.T) {
	testAccount, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var tests = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
		contractName    string
		methodName      string
	}{
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC1234567890123456@xuper",
			contractName:    "counter1",
			methodName:      "increase",
		},
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC1234567890123456@xuper",
			contractName:    "counter",
			methodName:      "get",
		},
	}

	for _, test := range tests {
		init := InitAcl(
			test.account,
			test.node,
			test.bcname,
		)
		if init == nil {
			t.Fatal("InitACL error")
		}

		acls, err := init.QueryMethodAcl(test.contractName, test.methodName)
		t.Logf("query contract method acl, res: %v, err: %v", acls.GetAcl(), err)

		addrs := []string{"nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh", "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpNr"}
		res, err := init.SetContractMethodAcl(test.contractAccount, test.contractName, test.methodName, addrs)
		t.Logf("set contract method acl, res: %v, err: %v", res, err)

		acls, err = init.QueryMethodAcl(test.contractName, test.methodName)
		t.Logf("query contract method acl, res: %v, err: %v", acls.GetAcl(), err)
	}
}

// test SetAccountAcl
func TestSetAccountAcl(t *testing.T) {
	testAccount, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var tests = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
	}{
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC1234567890123456@xuper",
		},
	}

	for _, test := range tests {
		init := InitAclSet(
			test.account,
			test.node,
			test.bcname,
			test.contractAccount,
			"", "",
		)
		if init == nil {
			t.Fatal("InitACL error")
		}

		acls, err := init.QueryAccountAcl(test.contractAccount)
		t.Logf("query contract account acl, res: %v, err: %v", acls.GetAcl(), err)

		addrs := map[string]float32{
			"nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh":  1.1,
			"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpNr": 1.2,
		}
		res, err := init.SetAccountAcl(addrs)
		t.Logf("set contract account acl, res: %v, err: %v", res, err)

		acls, err = init.QueryAccountAcl(test.contractAccount)
		t.Logf("query contract account acl, res: %v, err: %v", acls.GetAcl(), err)
	}
}

// test SetMethodAcl
func TestSetMethodAcl(t *testing.T) {
	testAccount, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		t.Fatalf("retrieveAccount err: %v\n", err)
	}
	t.Logf("RetrieveAccount: %v", testAccount)

	var tests = []struct {
		account         *account.Account
		node            string
		bcname          string
		contractAccount string
		contractName    string
		methodName      string
	}{
		{
			account:         testAccount,
			node:            node,
			bcname:          "xuper",
			contractAccount: "XC1234567890123456@xuper",
			contractName:    "counter",
			methodName:      "increase",
		},
	}

	for _, test := range tests {
		init := InitAclSet(
			test.account,
			test.node,
			test.bcname,
			test.contractAccount,
			test.contractName,
			test.methodName,
		)
		if init == nil {
			t.Fatal("InitACL error")
		}

		acls, err := init.QueryMethodAcl(test.contractName, test.methodName)
		t.Logf("query contract method acl, res: %v, err: %v", acls.GetAcl(), err)

		addrs := map[string]float32{
			"nuSMPvo6UUoTaT8mMQmHbfiRbJNbAymGh":  1.1,
			"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpNr": 1.2,
		}
		res, err := init.SetMethodAcl(addrs)
		t.Logf("set contract method acl, res: %v, err: %v", res, err)

		acls, err = init.QueryMethodAcl(test.contractName, test.methodName)
		t.Logf("query contract method acl, res: %v, err: %v", acls.GetAcl(), err)
	}
}
