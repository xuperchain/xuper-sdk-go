// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package transfer 转账
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

// InitTrans 创建转账客户端实例
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

// InitTransByPlatform 创建具有平台账户的转账客户端实例
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

// TransferWithDescFile 根据描述文件进行转账
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

// EncryptedTransfer 加密转账
func (t *Trans) EncryptedTransfer(to, amount, fee, desc, hdPublicKey string) (string, error) {
	if len(desc) == 0 {
		hdPublicKey = ""
	}
	return t.transfer(to, amount, fee, desc, hdPublicKey)
}

// Transfer 转账
//
//Parameters:
//   - `to`：转账目标地址。
//   - `amount`：转账金额。
//   - `fee`：转账交易给矿工的手续费。
//   - `desc`：交易的描述文件，可以为空。
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

// BatchTransfer 批量转账
//
//Parameters:
//   - `toAddressAndAmount`：转账目标地址已经对应金额的 map。
//   - `fee`：转账交易给矿工的手续费。
//   - `desc`：交易的描述文件，可以为空。
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

// QueryTx 根据交易 ID 查询交易
func (t *Trans) QueryTx(txid string) (*pb.TxStatus, error) {
	return t.Xchain.QueryTx(txid)
}

// DecryptedTx 解密交易的 desc
func (t *Trans) DecryptedTx(tx *pb.Transaction, privateAncestorKey string) (string, error) {
	cryptoClient := crypto.GetXchainCryptoClient()

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

// GetBalance 查询当前账户余额
func (t *Trans) GetBalance() (string, error) {
	if t.Account == nil {
		return "", common.ErrInvalidAccount
	}
	return t.GetBalanceDetail()
}
