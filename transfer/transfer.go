// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"io/ioutil"
	"log"
	//	"math/big"
	"errors"
	"strconv"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/crypto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"

	"github.com/xuperchain/crypto/core/utils"
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

// InitTrans init a client to transfer
func InitTransByPlatform(account, platformAccount *account.Account, node, bcname string) *Trans {
	commConfig := config.GetInstance()

	return &Trans{
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			PlatformAccount: platformAccount,
			XchainSer:       node,
			ChainName:       bcname,
		},
	}
}

func (t *Trans) TransferWithDescFile(to, amount, fee, descFilePath string) (string, error) {
	desc := ""
	if descFilePath != "" {
		descBytes, err := ioutil.ReadFile(descFilePath)
		if err != nil {
			return "", err
		}

		desc = string(descBytes)
	}

	return t.Transfer(to, amount, fee, desc)
}

func (t *Trans) EncryptedTransfer(to, amount, fee, desc, hdPublicKey string) (string, error) {
	if len(desc) == 0 {
		hdPublicKey = ""
	}
	return t.transfer(to, amount, fee, desc, hdPublicKey)
}

// Transfer transfer 'amount' to 'to',and pay 'fee' to miner
func (t *Trans) Transfer(to, amount, fee, desc string) (string, error) {
	return t.transfer(to, amount, fee, desc, "")
}

// Transfer transfer 'amount' to 'to',and pay 'fee' to miner
func (t *Trans) transfer(to, amount, fee, desc, hdPublicKey string) (string, error) {
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

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:    t.ChainName,
		Requests:  invokeRequests,
		Initiator: t.Account.Address,
		//		AuthRequire: authRequires,
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

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if t.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires := []string{}
		authRequires = append(authRequires, t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

		// 如果是平台发起的转账
		if t.Xchain.PlatformAccount != nil {
			authRequires = append(authRequires, t.Xchain.PlatformAccount.Address)
		}

		invokeRPCReq.AuthRequire = authRequires

		// 是否需要支付合规性背书费用
		if t.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	needTotalAmount := amountInt64 + extraAmount + feeInt64

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
	if preExeWithSelRes.Response == nil {
		return "", errors.New("preExe return nil")
	}

	// populates fields
	//	t.To = to
	t.Fee = fee
	t.Desc = desc
	t.InvokeRPCReq = invokeRPCReq
	t.Initiator = t.Account.Address
	//	t.Amount = strconv.FormatInt(amountInt64, 10)
	toAddressAndAmount := make(map[string]string)
	toAddressAndAmount[to] = amount
	t.ToAddressAndAmount = toAddressAndAmount
	t.TotalToAmount = amount //strconv.FormatInt(amountInt64, 10)

	// if ComplianceCheck is needed
	//	if t.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
	//		authRequires := []string{}
	//		authRequires = append(authRequires, t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	//		// 如果是平台发起的转账
	//		if t.Xchain.PlatformAccount != nil {
	//			authRequires = append(authRequires, t.Xchain.PlatformAccount.Address)
	//		}
	//		t.AuthRequire = authRequires
	//	}
	t.AuthRequire = invokeRPCReq.AuthRequire

	// post
	return t.GenCompleteTxAndPost(preExeWithSelRes, hdPublicKey)
}

// Transfer transfer 'amount' to 'to',and pay 'fee' to miner
func (t *Trans) BatchTransfer(toAddressAndAmount map[string]string, fee, desc string) (string, error) {
	//	var txOutputs []*pb.TxOutput

	// 求转出和
	//	realToAmountSum := new(big.Int)
	amountInt64 := int64(0)
	for _, toAmount := range toAddressAndAmount {
		singleAmountInt64, err := strconv.ParseInt(toAmount, 10, 64)
		if err != nil {
			return "", err
		}

		if singleAmountInt64 < 0 {
			return "", errors.New("Transfer amount is negative")
		}
		//		realToAmount, isSuccess := new(big.Int).SetString(amount, 10)
		//		if isSuccess != true {
		//			return "", errors.New("toAmount convert to bigint failed")
		//		}
		//		realToAmountSum = new(big.Int).Add(realToAmountSum, realToAmount)

		amountInt64 += singleAmountInt64
	}

	// generate preExe request
	invokeRequests := []*pb.InvokeRequest{}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:    t.ChainName,
		Requests:  invokeRequests,
		Initiator: t.Account.Address,
		//		AuthRequire: authRequires,
	}

	//	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	//	if err != nil {
	//		log.Printf("Transfer amount to int64 err: %v", err)
	//		return "", err
	//	}

	// get fee amount
	feeInt64, err := strconv.ParseInt(fee, 10, 64)
	if err != nil {
		log.Printf("Transfer fee to int64 err: %v", err)
		return "", err
	}
	if feeInt64 < 0 {
		return "", errors.New("fee amount is negative")
	}

	// get extra amount
	extraAmount := int64(0)
	// if ComplianceCheck is needed
	if t.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires := []string{}
		authRequires = append(authRequires, t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

		// 如果是平台发起的转账
		if t.Xchain.PlatformAccount != nil {
			authRequires = append(authRequires, t.Xchain.PlatformAccount.Address)
		}

		invokeRPCReq.AuthRequire = authRequires

		//		if amountInt64 < int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) {
		//			return "", common.ErrAmountNotEnough
		//		}

		//		extraAmount = int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if t.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	needTotalAmount := amountInt64 + extraAmount + feeInt64

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
	//	t.To = to
	//	t.Amount = strconv.FormatInt(amountInt64, 10)
	t.ToAddressAndAmount = toAddressAndAmount
	t.TotalToAmount = strconv.FormatInt(amountInt64, 10) //string(amountInt64)
	t.Fee = fee
	t.Desc = desc
	t.InvokeRPCReq = invokeRPCReq
	t.Initiator = t.Account.Address

	// if ComplianceCheck is needed
	//	if t.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
	//		authRequires := []string{}
	//		authRequires = append(authRequires, t.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	//
	//		// 如果是平台发起的转账
	//		if t.Xchain.PlatformAccount != nil {
	//			authRequires = append(authRequires, t.Xchain.PlatformAccount.Address)
	//		}
	//
	//		t.AuthRequire = authRequires
	//	}
	t.AuthRequire = invokeRPCReq.AuthRequire

	// post
	return t.GenCompleteTxAndPost(preExeWithSelRes, "")
}

// QueryTx query tx to get detail information
func (t *Trans) QueryTx(txid string) (*pb.TxStatus, error) {
	return t.Xchain.QueryTx(txid)
}

// QueryTx query tx to get detail information
func (t *Trans) DecryptedTx(tx *pb.Transaction, privateAncestorKey string) (string, error) {
	cryptoClient := crypto.GetCryptoClient()

	originalDesc, err := cryptoClient.DecryptByHdKey(string(tx.HDInfo.HdPublicKey), privateAncestorKey, string(tx.Desc))
	if err != nil {
		return "", err
	}

	originalHash := cryptoClient.HashUsingDoubleSha256([]byte(originalDesc))
	if !utils.BytesCompare(originalHash, tx.HDInfo.OriginalHash) {
		return "", errors.New("originalHash doesn't match the originalDesc")
	}

	return originalDesc, nil
}

// GetBalance get your own balance
func (t *Trans) GetBalance() (string, error) {
	if t.Account == nil {
		return "", common.ErrInvalidAccount
	}
	return t.GetBalanceDetail()
}
