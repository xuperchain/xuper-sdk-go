// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to common variables and utils funcs
package common

import (
	"testing"
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
