// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"log"
	"strconv"

	"github.com/xuperchain/xuperchain/core/pb"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

// Trans transaction structure
type Trans struct {
	xchain.Xchain
}

// InitTrans init a client to transfer
func InitTrans(account *account.Account, node, bcname string) *Trans {
	commConfig := config.GetInstance()

	return &Trans{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			XchainSer: node,
			ChainName: bcname,
		},
	}
}

// Transfer transfer 'amount' to 'to',and pay 'fee' to miner
func (t *Trans) Transfer(to, amount, fee, desc string) (string, error) {
	// (total pay amount) = (to amount + fee + checkfee)
	amount, ok := common.IsValidAmount(amount)
	if !ok {
		return "", common.ErrInvalidAmount
	}
	fee, ok = common.IsValidAmount(fee)
	if !ok {
		return "", common.ErrInvalidAmount
	}
	// generate preExe request
	invokeRequests := []*pb.InvokeRequest{}
	authRequires := []string{}
	authRequires = append(authRequires, t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      t.ChainName,
		Requests:    invokeRequests,
		Initiator:   t.Account.Address,
		AuthRequire: authRequires,
	}

	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("Transfer amount to int64 err: %v", err)
		return "", err
	}
	feeInt64, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		log.Printf("Transfer fee to int64 err: %v", err)
		return "", err
	}
	if amountInt64 < int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) {
		return "", common.ErrAmountNotEnough
	}

	needTotalAmount := amountInt64 + int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) + feeInt64

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      t.ChainName,
		Address:     t.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}
	t.PreSelUTXOReq = preSelUTXOReq

	// preExe
	preExeWithSelRes, err := t.PreExecWithSelecUTXO()
	if err != nil {
		log.Printf("Transfer PreExecWithSelecUTXO failed, err: %v", err)
		return "", err
	}

	// populates fields
	t.To = to
	t.Fee = fee
	t.DescFile = desc
	t.InvokeRPCReq = invokeRPCReq
	t.Initiator = t.Account.Address
	t.AuthRequire = authRequires
	t.Amount = strconv.FormatInt(amountInt64, 10)

	// post
	return t.GenCompleteTxAndPost(preExeWithSelRes)
}

// QueryTx query tx to get detail information
func (t *Trans) QueryTx(txid string) (*pb.TxStatus, error) {
	return t.Xchain.QueryTx(txid)
}

// GetBalance get your own balance
func (t *Trans) GetBalance() (string, error) {
	if t.Account == nil {
		return "", common.ErrInvalidAccount
	}
	return t.GetBalanceDetail()
}
