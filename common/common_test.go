// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to common variables and utils funcs
package common

import (
	"os"
	"testing"

	"github.com/xuperchain/xuperchain/core/pb"
)

func TestIsValidAmount(t *testing.T) {
	testCase := []string{
		"",
		"0",
		"345",
		"-345",
		"-34fdsafds5",
	}

	for _, arg := range testCase {
		amount, ok := IsValidAmount(arg)
		t.Logf("amount: %v, err: %v", amount, ok)
	}
}

func TestSeed(t *testing.T) {
	err := SetSeed()
	if err != nil {
		t.Error("SetSeed assert failed")
	}
	n := GetNonce()
	if n == "" {
		t.Error("GetNonce assert failed")
	}

	p := ""
	err = PathExistsAndMkdir(p)
	if err == nil {
		t.Error("PathExistsAndMkdir assert failed")
	}

	p = "./tmp"
	err = PathExistsAndMkdir(p)
	if err != nil {
		t.Error("PathExistsAndMkdir assert failed")
	}
	os.RemoveAll(p)
}

func TestVaildAmount(t *testing.T) {

	_, ok := IsValidAmount("a")
	if ok {
		t.Error("TestVaildAmount assert failed")
	}
	_, ok = IsValidAmount("-100")
	if ok {
		t.Error("TestVaildAmount assert failed")
	}
	_, ok = IsValidAmount("1")
	if !ok {
		t.Error("TestVaildAmount assert failed")
	}
}

func TestTxHash(t *testing.T) {
	_, e := MakeTxDigestHash(&pb.Transaction{})
	if e != nil {
		t.Error("MakeTxDigestHash assert failed")
	}

	_, e = MakeTransactionID(&pb.Transaction{})
	if e != nil {
		t.Error("MakeTransactionID assert failed")
	}
}
