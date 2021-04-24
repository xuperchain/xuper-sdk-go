// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package transfer is related to transfer operation
package transfer

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/xuperchain/xuper-sdk-go/txhash"
	"io/ioutil"
	"log"
	"math/big"
	"path/filepath"
	"time"

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

const (
	feePlaceholder            = "$"

)

// InitTrans init a client to transfer
func InitTrans(account *account.Account,bcname string,sdkClient *xchain.SDKClient) *Trans {
	commConfig := config.GetInstance()

	return &Trans{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			//XchainSer: node,
			ChainName: bcname,
			SDKClient: sdkClient,
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
	//ctx := context.Background()
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
	if t.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
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
		t.AuthRequire = invokeRPCReq.AuthRequire

		// post
		return t.GenCompleteTxAndPost(preExeWithSelRes, hdPublicKey)
	}else{
		return t.tansferSupportAccount(to, amount, fee, desc)
	}



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

// GetBalance get your own balance
func (t *Trans) GetBalance() (string, error) {
	if t.Account == nil {
		return "", common.ErrInvalidAccount
	}
	return t.GetBalanceDetail()
}




//transfer(to, amount, fee, desc, hdPublicKey string)


func (t *Trans) tansferSupportAccount(to, amount, fee, desc string)(string,error){
	ctx := context.Background()
	opt := &TransferOptions{
		BlockchainName:t.ChainName,
		To:to,
		Amount:amount,
		Fee:fee,
		Desc:[]byte(desc),
		From:t.Account.Address,
		Version:1,						//todo
	}
	txStatus,err := t.assembleTxSupportAccount(ctx,opt)
	if err != nil {
		fmt.Println("assembleTxSupportAccount error")
		return "",err
	}

	signTx, err := t.ProcessSignTx(txStatus.Tx)
	if err != nil {
		return "",err
	}
	signInfo := &pb.SignatureInfo{
		PublicKey: t.Account.PublicKey,
		Sign:      signTx,
	}

	txStatus.Tx.InitiatorSigns = append(txStatus.Tx.InitiatorSigns, signInfo)
	txStatus.Tx.AuthRequireSigns, err = t.genAuthRequireSigns(opt,txStatus.Tx)
	if err != nil {
		return "", fmt.Errorf("Failed to genAuthRequireSigns %s", err)
	}
	txStatus.Tx.Txid, err = txhash.MakeTransactionID(txStatus.Tx)
	if err != nil {
		return "", fmt.Errorf("Failed to gen txid %s", err)
	}
	txStatus.Txid = txStatus.Tx.Txid

	client := *(t.SDKClient.XchainClient)

	// 提交
	reply, err := client.PostTx(ctx, txStatus)
	if err != nil {
		return "", fmt.Errorf("transferSupportAccount post tx err %s", err)
	}
	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return "", fmt.Errorf("Failed to post tx: %s", reply.Header.String())
	}
	return hex.EncodeToString(txStatus.GetTxid()), nil


}



func (t *Trans) genAuthRequireSigns(opt *TransferOptions, tx *pb.Transaction) ([]*pb.SignatureInfo, error) {
	authRequireSigns := []*pb.SignatureInfo{}
	signTx, err := t.ProcessSignTx(tx)
	if err != nil {
		return nil, err
	}
	signInfo := &pb.SignatureInfo{
		PublicKey: t.Account.PrivateKey,
		Sign:      signTx,
	}
	authRequireSigns = append(authRequireSigns, signInfo)
	return authRequireSigns, nil
}





func (t *Trans) ProcessSignTx(tx *pb.Transaction) ([]byte, error) {

	client := crypto.GetCryptoClient()
	privateKey, err := client.GetEcdsaPrivateKeyFromJsonStr(t.Account.PrivateKey)
	if err != nil {
		return nil, err
	}
	digestHash, dhErr := txhash.MakeTxDigestHash(tx)
	if dhErr != nil {
		return nil, dhErr
	}
	sign, sErr := client.SignECDSA(privateKey, digestHash)
	if sErr != nil {
		return nil, sErr
	}
	return sign, nil
}




func newFeeAccount(fee string) *pb.TxDataAccount {
	return &pb.TxDataAccount{
		Address: feePlaceholder,
		Amount:  fee,
	}
}


type TransferOptions struct {
	BlockchainName string
	KeyPath        string
	CryptoType     string
	To             string
	Amount         string
	Fee            string
	Desc           []byte
	FrozenHeight   int64
	Version        int32
	// 支持账户转账
	From        string
	AccountPath string
}


func (t *Trans) assembleTxSupportAccount(ctx context.Context,opt *TransferOptions)(*pb.TxStatus,error){
	bigZero := big.NewInt(0)
	totalNeed := big.NewInt(0)
	tx := &pb.Transaction{
		Version:   opt.Version,
		Coinbase:  false,
		Desc:     []byte(t.Desc),
		Nonce:     common.GetNonce(),
		Timestamp: time.Now().UnixNano(),
		Initiator: t.Initiator,
	}
	account := &pb.TxDataAccount{
		Address:      opt.To,
		Amount:       opt.Amount,
		FrozenHeight: opt.FrozenHeight,
	}
	accounts := []*pb.TxDataAccount{account}
	if opt.Fee != "" && opt.Fee != "0" {
		accounts = append(accounts, newFeeAccount(opt.Fee))
	}
	// 组装output
	for _, acc := range accounts {
		amount, ok := big.NewInt(0).SetString(acc.Amount, 10)
		if !ok {
			return nil,  errors.New("Invalid amount number")
		}
		if amount.Cmp(bigZero) < 0 {
			return nil, errors.New("Amount in transaction can not be negative number")
		}
		totalNeed.Add(totalNeed, amount)
		txOutput := &pb.TxOutput{}
		txOutput.ToAddr = []byte(acc.Address)
		txOutput.Amount = amount.Bytes()
		txOutput.FrozenHeight = acc.FrozenHeight
		tx.TxOutputs = append(tx.TxOutputs, txOutput)
	}
	// 组装input 和 剩余output
	txInputs, deltaTxOutput, err := t.assembleTxInputsSupportAccount(ctx,opt, totalNeed)
	if err != nil {
		return nil, err
	}
	tx.TxInputs = txInputs
	if deltaTxOutput != nil {
		tx.TxOutputs = append(tx.TxOutputs, deltaTxOutput)
	}
	// 设置auth require
	tx.AuthRequire, err = genAuthRequire(opt.From, opt.AccountPath)
	if err != nil {
		return nil, err
	}

	preExeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      opt.BlockchainName,
		Requests:    []*pb.InvokeRequest{},
		//Header:      global.GHeader(),
		Initiator:   t.Initiator,
		AuthRequire: tx.AuthRequire,
	}
	client := *(t.SDKClient.XchainClient)

	preExeRes, err := client.PreExec(ctx, preExeRPCReq)
	if err != nil {
		return nil, err
	}

	tx.ContractRequests = preExeRes.GetResponse().GetRequests()
	tx.TxInputsExt = preExeRes.GetResponse().GetInputs()
	tx.TxOutputsExt = preExeRes.GetResponse().GetOutputs()


	txStatus := &pb.TxStatus{
		Bcname: opt.BlockchainName,
		Status: pb.TransactionStatus_UNCONFIRM,
		Tx:     tx,
	}
	//txStatus.Header = &pb.Header{
	//	Logid: global.Glogid(),
	//}
	return txStatus, nil
}


func genAuthRequire(from, path string) ([]string, error) {				// 此处需要改造
	authRequire := []string{}
	if path == "" {
		authRequire = append(authRequire, from)
		return authRequire, nil
	}
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, fi := range dir {
		if fi.IsDir() {
			addr, err := readAddress(path + "/" + fi.Name())
			if err != nil {
				return nil, err
			}
			authRequire = append(authRequire, from+"/"+addr)
		}
	}
	return authRequire, nil
}

func readAddress(keypath string) (string, error) {
	return readKeys(filepath.Join(keypath, "address"))
}

func readKeys(file string) (string, error) {
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	buf = bytes.TrimSpace(buf)
	return string(buf), nil
}

func (t *Trans) assembleTxInputsSupportAccount(ctx context.Context, opt *TransferOptions, totalNeed *big.Int) ([]*pb.TxInput, *pb.TxOutput, error) {
	ui := &pb.UtxoInput{
		Bcname:    opt.BlockchainName,
		Address:   opt.From,
		TotalNeed: totalNeed.String(),
		NeedLock:  true,
	}

	//client := t.
	//
	//	cryptoClient := crypto.GetCryptoClient()
	//privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(xc.Account.PrivateKey)


	client := *(t.SDKClient.XchainClient)

	utxoRes, selectErr := client.SelectUTXO(ctx, ui)
	if selectErr != nil || utxoRes.Header.Error != pb.XChainErrorEnum_SUCCESS {
		fmt.Println("Select utxo error:")
		fmt.Println(utxoRes.Header.Error)

		return nil, nil, errors.New("Select utxo error")
	}
	var txTxInputs []*pb.TxInput
	var txOutput *pb.TxOutput
	for _, utxo := range utxoRes.UtxoList {
		txInput := new(pb.TxInput)
		txInput.RefTxid = utxo.RefTxid
		txInput.RefOffset = utxo.RefOffset
		txInput.FromAddr = utxo.ToAddr
		txInput.Amount = utxo.Amount
		txTxInputs = append(txTxInputs, txInput)
	}
	utxoTotal, ok := big.NewInt(0).SetString(utxoRes.TotalSelected, 10)
	if !ok {
		return nil, nil, errors.New("Select utxo error")
	}
	// 多出来的utxo需要再转给自己
	if utxoTotal.Cmp(totalNeed) > 0 {
		delta := utxoTotal.Sub(utxoTotal, totalNeed)
		txOutput = &pb.TxOutput{
			ToAddr: []byte(opt.From), // 收款人就是汇款人自己
			Amount: delta.Bytes(),
		}
	}
	return txTxInputs, txOutput, nil
}













