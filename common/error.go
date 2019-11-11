// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to common variables and utils funcs
package common

import (
	"errors"
)

var (
	// ErrInvalidAmount invalid amount
	ErrInvalidAmount = errors.New("invalid amount")
	// ErrTxNotFound tx is not found
	ErrTxNotFound = errors.New("tx not found")
	// ErrInvalidAccount invalid account
	ErrInvalidAccount = errors.New("invalid account")
	// ErrInvalidContractAccount
	ErrInvalidContractAccount = errors.New("conrtact account must be numbers of length 16")
	// ErrAmountNotEnough
	ErrAmountNotEnough = errors.New("Amount must be bigger than compliancecheck fee which is 10")
)
