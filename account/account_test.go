// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package account is related to account operation
package account

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/xuperchain/xuper-sdk-go/v2/common"
)

func TestCreateAccount(t *testing.T) {
	testCase := []struct {
		strength uint8
		language int
	}{
		{
			strength: 1,
			language: 1,
		},
		{
			strength: 2,
			language: 1,
		},
		{
			strength: 3,
			language: 1,
		},
		{
			strength: 1,
			language: 2,
		},
		{
			strength: 2,
			language: 2,
		},
		{
			strength: 3,
			language: 2,
		},
		{
			strength: 0,
			language: 5,
		},
	}

	for _, arg := range testCase {
		acc, err := CreateAccount(arg.strength, arg.language)
		t.Logf("create account: %v, err: %v", acc, err)
	}
}

func TestRetrieveAccount(t *testing.T) {
	testCase := []struct {
		mnemonic string
		language int
	}{
		{
			mnemonic: "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即",
			language: 1,
		},
		{
			mnemonic: "玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即",
			language: 2,
		},
		{
			mnemonic: "",
			language: 1,
		},
	}

	for _, arg := range testCase {
		acc, err := RetrieveAccount(arg.mnemonic, arg.language)
		t.Logf("RetrieveAccount: %v, err: %v", acc, err)
	}
}

func TestCreateAndSaveAccountToFile(t *testing.T) {
	testCase := []struct {
		path     string
		passwd   string
		strength uint8
		language int
	}{
		{
			path:     "./keys",
			passwd:   "123",
			strength: 1,
			language: 1,
		},
		{
			path:     "./keys",
			passwd:   "123",
			strength: 1,
			language: 2,
		},
		{
			path:     "./aaa",
			passwd:   "123",
			strength: 1,
			language: 2,
		},
	}

	for _, arg := range testCase {
		acc, err := CreateAndSaveAccountToFile(arg.path, arg.passwd, arg.strength, arg.language)
		t.Logf("CreateAndSaveAccountToFile: %v, err: %v", acc, err)
		fmt.Println(os.RemoveAll(arg.path))
	}
}

func TestGetAccountFromFile(t *testing.T) {
	testCase := []struct {
		path   string
		passwd string
	}{
		{
			path:   "./keys/",
			passwd: "123",
		},
		{
			path:   "./aaa/",
			passwd: "123",
		},
	}

	for _, arg := range testCase {
		CreateAndSaveAccountToFile(arg.path, arg.passwd, 1, 1)

		acc, err := GetAccountFromFile(arg.path, arg.passwd)
		if err != nil {
			t.Error(err)
		}
		if acc == nil {
			t.Error("GetAccountFromFile assert failed")
		}
		os.RemoveAll(arg.path)
	}
}

func TestSetContractAccount(t *testing.T) {
	acc, _ := CreateAccount(1, 1)
	err := acc.SetContractAccount("123")
	if !errors.Is(err, common.ErrInvalidContractAccount) {
		t.Error(err)
	}

	err = acc.SetContractAccount("XC123@xuper")
	if !errors.Is(err, common.ErrInvalidContractAccount) {
		t.Error(err)
	}

	err = acc.SetContractAccount("1234567812345678@xuper")
	if !errors.Is(err, common.ErrInvalidContractAccount) {
		t.Error(err)
	}

	err = acc.SetContractAccount("XC1234567812345678@xuper")
	if err != nil {
		t.Error(err)
	}

	ar := acc.GetAuthRequire()
	if ar != "XC1234567812345678@xuper/"+acc.Address {
		t.Error("account authRequire assert failed")
	}

	acc.RemoveContractAccount()
	if acc.HasContractAccount() {
		t.Error("Remove contract account test failed")
	}
	ar = acc.GetAuthRequire()
	if ar != acc.Address {
		t.Error("account authRequire assert failed")
	}
}
