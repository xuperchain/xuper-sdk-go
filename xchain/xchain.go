// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package xchain is related to xchain operation
package xchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"google.golang.org/grpc"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/crypto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/txhash"
	//	"github.com/xuperchain/crypto/core/account"
)

// Xchain xchain struct
type Xchain struct {
	Cfg *config.CommConfig
	//	To     string
	//	Amount string
	ToAddressAndAmount map[string]string
	TotalToAmount      string
	Fee                string
	//	DescFile              string
	Desc                  string
	FrozenHeight          int64
	InvokeRPCReq          *pb.InvokeRPCRequest
	PreSelUTXOReq         *pb.PreExecWithSelectUTXORequest
	Initiator             string
	AuthRequire           []string
	Account               *account.Account
	PlatformAccount       *account.Account
	ChainName             string
	XchainSer             string
	ContractAccount       string
	IsNeedComplianceCheck bool
}

// PreExecWithSelecUTXO preExec and selectUTXO
func (xc *Xchain) PreExecWithSelecUTXO() (*pb.PreExecWithSelectUTXOResponse, error) {
	requestData, err := json.Marshal(xc.PreSelUTXOReq)
	if err != nil {
		log.Printf("PreExecWithSelecUTXO json marshal failed, err: %v", err)
		return nil, err
	}

	endorserRequest := &pb.EndorserRequest{
		RequestName: "PreExecWithFee",
		BcName:      xc.ChainName,
		RequestData: requestData,
	}

	conn, err := grpc.Dial(xc.Cfg.EndorseServiceHost, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("PreExecWithSelecUTXO Connect EndorseServiceHost failed, err: %v", err)
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()
	c := pb.NewXendorserClient(conn)
	endorserResponse, err := c.EndorserCall(ctx, endorserRequest)
	if err != nil {
		log.Printf("PreExecWithSelecUTXO EndorserCall failed, err: %v", err)
		return nil, fmt.Errorf("EndorserCall error! Response is: %v", err)
	}

	responseData := endorserResponse.ResponseData

	preExecWithSelectUTXOResponse := new(pb.PreExecWithSelectUTXOResponse)
	err = json.Unmarshal(responseData, preExecWithSelectUTXOResponse)
	if err != nil {
		return nil, err
	}

	log.Printf("Gas will cost: %v\n", preExecWithSelectUTXOResponse.GetResponse().GetGasUsed())
	for _, res := range preExecWithSelectUTXOResponse.GetResponse().GetResponses() {
		if res.Status >= 400 {
			return nil, fmt.Errorf("contract error status:%d message:%s", res.Status, res.Message)
		}
		log.Printf("contract response: %s\n", string(res.Body))
	}

	return preExecWithSelectUTXOResponse, nil
}

// GenComplianceCheckTx generate complianceTx to pay for compliance check
func (xc *Xchain) GenComplianceCheckTx(response *pb.PreExecWithSelectUTXOResponse) (
	*pb.Transaction, error) {

	totalNeed := new(big.Int).SetInt64(int64(xc.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee))
	txInputs, deltaTxOutput, err := xc.GenerateTxInput(response.GetUtxoOutput(), totalNeed)
	//	txInputs, err := xc.GeneratePureTxInputs(response.GetUtxoOutput())
	if err != nil {
		log.Printf("GenerateComplianceTx GenerateTxInput failed.")
		return nil, fmt.Errorf("GenerateComplianceTx GenerateTxInput err: %v", err)
	}

	checkAmount := strconv.Itoa(xc.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
	txOutputs, err := xc.GenerateTxOutput(xc.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFeeAddr, checkAmount, "0")
	if err != nil {
		log.Printf("GenerateComplianceTx GenerateTxOutput failed.")
		return nil, fmt.Errorf("GenerateComplianceTx GenerateTxOutput err: %v", err)
	}
	if deltaTxOutput != nil {
		txOutputs = append(txOutputs, deltaTxOutput)
	}
	// populates fields
	tx := &pb.Transaction{
		Desc:      []byte(""),
		Version:   common.TxVersion,
		Coinbase:  false,
		Timestamp: time.Now().UnixNano(),
		TxInputs:  txInputs,
		TxOutputs: txOutputs,
		Initiator: xc.Initiator,
	}

	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()

	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(xc.Account.PrivateKey)
	if err != nil {
		return nil, err
	}
	digestHash, dhErr := txhash.MakeTxDigestHash(tx)
	if dhErr != nil {
		return nil, dhErr
	}
	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)

	signatureInfo := &pb.SignatureInfo{
		PublicKey: xc.Account.PublicKey,
		Sign:      sign,
	}

	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)

	tx.InitiatorSigns = signatureInfos

	// make txid
	tx.Txid, _ = txhash.MakeTransactionID(tx)
	return tx, nil
}

// GenerateTxOutput generate txoutput part
func (xc *Xchain) GenerateMultiTxOutputs(selfAmount string) ([]*pb.TxOutput, error) {
	selfAddr := xc.Account.Address
	toAddrAndAmount := xc.ToAddressAndAmount
	//	selfAmount := ""
	feeAmount := xc.Fee

	//	txOutputs := make([]TxOutput, 2)
	var txOutputs []*pb.TxOutput
	// 1.向目标账户转账
	for toAddr, toAmount := range toAddrAndAmount {
		// 向目标地址转账的txOutput
		txOutputTo := new(pb.TxOutput)
		txOutputTo.ToAddr = []byte(toAddr)
		realToAmount, isSuccess := new(big.Int).SetString(toAmount, 10)
		if isSuccess != true {
			log.Printf("toAmount convert to bigint failed")
			return nil, common.ErrInvalidAmount
		}
		txOutputTo.Amount = realToAmount.Bytes()
		//	txOutputTo.FrozenHeight = 0
		//	txOutputs[0] = *txOutputTo
		txOutputs = append(txOutputs, txOutputTo)
	}
	// 2.自己的转账地址接收差额部分的转回的txOutput
	txOutputSelf := new(pb.TxOutput)
	txOutputSelf.ToAddr = []byte(selfAddr)
	realSelfAmount, isSuccess := new(big.Int).SetString(selfAmount, 10)
	if isSuccess != true {
		log.Printf("selfAmount convert to bigint failed")
		return nil, common.ErrInvalidAmount
	}
	txOutputSelf.Amount = realSelfAmount.Bytes()
	//	txOutputs[1] = *txOutputSelf
	txOutputs = append(txOutputs, txOutputSelf)
	// 3.如果矿工手续费不是空
	if feeAmount != "" && feeAmount != "0" {
		//		realFeeAmount := new(big.Int).SetBytes(feeAmount)
		realFeeAmount, isSuccess := new(big.Int).SetString(feeAmount, 10)
		if isSuccess != true {
			log.Printf("feeAmount convert to bigint failed")
			return nil, common.ErrInvalidAmount
		}
		// 如果矿工手续费不合法
		if realFeeAmount.Cmp(big.NewInt(0)) < 0 {
			return nil, common.ErrInvalidAmount
		}
		// 支付给矿工的手续费相关的txOutput
		txOutputFee := new(pb.TxOutput)
		txOutputFee.ToAddr = []byte("$")
		txOutputFee.Amount = realFeeAmount.Bytes()
		txOutputs = append(txOutputs, txOutputFee)
	}
	//	txOutputList := new(pb.TxOutputs)
	//	txOutputList.TxOutputList = txOutputs
	//	return txOutputList, nil

	return txOutputs, nil
}

// GenerateTxOutput generate txoutput part
func (xc *Xchain) GenerateTxOutput(to, amount, fee string) ([]*pb.TxOutput, error) {
	accounts := []*pb.TxDataAccount{}
	if to != "" {
		account := &pb.TxDataAccount{
			Address:      to,
			Amount:       amount,
			FrozenHeight: 0,
		}
		accounts = append(accounts, account)
	}
	if fee != "0" {
		feeAccount := &pb.TxDataAccount{
			Address: "$",
			Amount:  fee,
		}
		accounts = append(accounts, feeAccount)
	}

	bigZero := big.NewInt(0)
	txOutputs := []*pb.TxOutput{}
	for _, acc := range accounts {
		amount, ok := big.NewInt(0).SetString(acc.Amount, 10)
		if !ok {
			return nil, common.ErrInvalidAmount
		}
		cmpRes := amount.Cmp(bigZero)
		if cmpRes < 0 {
			return nil, errors.New("Invalid negative number")
		} else if cmpRes == 0 {
			continue
		}
		txOutput := &pb.TxOutput{}
		txOutput.Amount = amount.Bytes()
		txOutput.ToAddr = []byte(acc.Address)
		txOutput.FrozenHeight = acc.FrozenHeight
		txOutputs = append(txOutputs, txOutput)
	}

	return txOutputs, nil
}

func (xc *Xchain) GeneratePureTxInputs(utxoOutputs *pb.UtxoOutput) (
	[]*pb.TxInput, error) {
	// utxoList => TxInput
	//
	// gen txInputs
	var txInputs []*pb.TxInput
	//	var txOutput *pb.TxOutput
	for _, utxo := range utxoOutputs.UtxoList {
		txInput := &pb.TxInput{}
		txInput.RefTxid = utxo.RefTxid
		txInput.RefOffset = utxo.RefOffset
		txInput.FromAddr = utxo.ToAddr
		txInput.Amount = utxo.Amount
		txInputs = append(txInputs, txInput)
	}

	return txInputs, nil
}

// GenerateTxInput generate txinput part
func (xc *Xchain) GenerateTxInput(utxoOutputs *pb.UtxoOutput, totalNeed *big.Int) (
	[]*pb.TxInput, *pb.TxOutput, error) {
	// utxoList => TxInput
	//
	// gen txInputs
	var txInputs []*pb.TxInput
	var txOutput *pb.TxOutput
	for _, utxo := range utxoOutputs.UtxoList {
		txInput := &pb.TxInput{}
		txInput.RefTxid = utxo.RefTxid
		txInput.RefOffset = utxo.RefOffset
		txInput.FromAddr = utxo.ToAddr
		txInput.Amount = utxo.Amount
		txInputs = append(txInputs, txInput)
	}

	utxoTotal, ok := big.NewInt(0).SetString(utxoOutputs.TotalSelected, 10)
	if !ok {
		return nil, nil, fmt.Errorf("GenerateTxInput totalSelected err: %v", ok)
	}

	// input > output, generate output-input to me
	if utxoTotal.Cmp(totalNeed) > 0 {
		delta := utxoTotal.Sub(utxoTotal, totalNeed)
		txOutput = &pb.TxOutput{
			ToAddr: []byte(xc.Account.Address),
			Amount: delta.Bytes(),
		}
	}

	return txInputs, txOutput, nil
}

// GenRealTx generate really effective transaction
func (xc *Xchain) GenRealTx(response *pb.PreExecWithSelectUTXOResponse,
	complianceCheckTx *pb.Transaction, hdPublicKey string) (*pb.Transaction, error) {
	utxolist := []*pb.Utxo{}
	totalSelected := big.NewInt(0)
	for index, txOutput := range complianceCheckTx.TxOutputs {
		if string(txOutput.ToAddr) == xc.Initiator {
			utxo := &pb.Utxo{
				Amount:    txOutput.Amount,
				ToAddr:    txOutput.ToAddr,
				RefTxid:   complianceCheckTx.Txid,
				RefOffset: int32(index),
			}
			utxolist = append(utxolist, utxo)

			utxoAmount := big.NewInt(0).SetBytes(utxo.Amount)
			totalSelected.Add(totalSelected, utxoAmount)
		}
	}
	utxoOutput := &pb.UtxoOutput{
		UtxoList:      utxolist,
		TotalSelected: totalSelected.String(),
	}

	totalNeed := big.NewInt(0)

	// no need to double check
	amount, ok := big.NewInt(0).SetString(xc.TotalToAmount, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	fee, ok := big.NewInt(0).SetString(xc.Fee, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	amount.Add(amount, fee)
	totalNeed.Add(totalNeed, amount)

	//	txOutputs, err := xc.GenerateTxOutput(xc.To, xc.Amount, xc.Fee)
	selfAmount := totalSelected.Sub(totalSelected, totalNeed)
	txOutputs, err := xc.GenerateMultiTxOutputs(selfAmount.String())
	if err != nil {
		log.Printf("GenRealTx GenerateTxOutput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxOutput err: %v", err)
	}

	//	txInputs, deltaTxOutput, err := xc.GenerateTxInput(utxoOutput, totalNeed)
	txInputs, err := xc.GeneratePureTxInputs(utxoOutput)
	if err != nil {
		log.Printf("GenRealTx GenerateTxInput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxInput err: %v", err)
	}

	//	if deltaTxOutput != nil {
	//		txOutputs = append(txOutputs, deltaTxOutput)
	//	}

	tx := &pb.Transaction{
		Desc:      []byte("Maybe common transfer transaction"),
		Version:   common.TxVersion,
		Coinbase:  false,
		Timestamp: time.Now().UnixNano(),
		TxInputs:  txInputs,
		TxOutputs: txOutputs,
		Initiator: xc.Initiator,
		//		AuthRequire: xc.AuthRequire,
	}

	if len(xc.AuthRequire) != 0 {
		tx.AuthRequire = xc.AuthRequire
	}

	//	tx.Desc = []byte(xc.Desc)

	cryptoClient := crypto.GetCryptoClient()

	if len(hdPublicKey) == 0 {
		// 如果不需要HD分层加密功能
		tx.Desc = []byte(xc.Desc)
	} else {
		// 如果需要HD分层加密功能
		cypherText, err := cryptoClient.EncryptByHdKey(hdPublicKey, xc.Desc)
		if err != nil {
			return nil, err
		}

		tx.Desc = []byte(cypherText)

		// 继续组装HDInfo
		originalHash := cryptoClient.HashUsingDoubleSha256([]byte(xc.Desc))

		hdInfo := &pb.HDInfo{
			HdPublicKey:  []byte(hdPublicKey),
			OriginalHash: originalHash,
		}

		tx.HDInfo = hdInfo
	}

	tx.TxInputsExt = response.GetResponse().GetInputs()
	tx.TxOutputsExt = response.GetResponse().GetOutputs()
	tx.ContractRequests = response.GetResponse().GetRequests()

	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()

	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(xc.Account.PrivateKey)
	if err != nil {
		return nil, err
	}
	digestHash, dhErr := txhash.MakeTxDigestHash(tx)
	if dhErr != nil {
		return nil, dhErr
	}
	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)
	signatureInfo := &pb.SignatureInfo{
		PublicKey: xc.Account.PublicKey,
		Sign:      sign,
	}
	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)
	tx.InitiatorSigns = signatureInfos
	if xc.ContractAccount != "" {
		tx.AuthRequireSigns = signatureInfos
	}
	// make txid
	tx.Txid, _ = txhash.MakeTransactionID(tx)
	return tx, nil

}

// GenRealTxOnly generate really effective transaction
func (xc *Xchain) GenRealTxOnly(response *pb.PreExecWithSelectUTXOResponse, hdPublicKey string) (*pb.Transaction, error) {
	//	txOutputs, err := xc.GenerateTxOutput(xc.To, xc.Amount, xc.Fee)
	//	if err != nil {
	//		log.Printf("GenRealTx GenerateTxOutput failed.")
	//		return nil, fmt.Errorf("GenRealTx GenerateTxOutput err: %v", err)
	//	}

	//	utxolist := []*pb.Utxo{}
	//	totalSelected := big.NewInt(0)
	//	for index, txOutput := range complianceCheckTx.TxOutputs {
	//		if string(txOutput.ToAddr) == xc.Initiator {
	//			utxo := &pb.Utxo{
	//				Amount:    txOutput.Amount,
	//				ToAddr:    txOutput.ToAddr,
	//				RefTxid:   complianceCheckTx.Txid,
	//				RefOffset: int32(index),
	//			}
	//			utxolist = append(utxolist, utxo)
	//
	//			utxoAmount := big.NewInt(0).SetBytes(utxo.Amount)
	//			totalSelected.Add(totalSelected, utxoAmount)
	//		}
	//	}

	utxoOutput := &pb.UtxoOutput{
		//		UtxoList: utxolist,
		//		TotalSelected: totalSelected.String(),
		UtxoList:      response.UtxoOutput.UtxoList,
		TotalSelected: response.UtxoOutput.TotalSelected,
	}
	totalNeed := big.NewInt(0)
	amount, ok := big.NewInt(0).SetString(xc.TotalToAmount, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	fee, ok := big.NewInt(0).SetString(xc.Fee, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	amount.Add(amount, fee)
	totalNeed.Add(totalNeed, amount)

	totalSelected, ok := big.NewInt(0).SetString(response.UtxoOutput.TotalSelected, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}

	selfAmount := totalSelected.Sub(totalSelected, totalNeed)
	txOutputs, err := xc.GenerateMultiTxOutputs(selfAmount.String())
	if err != nil {
		log.Printf("GenRealTx GenerateTxOutput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxOutput err: %v", err)
	}

	//	txInputs, deltaTxOutput, err := xc.GenerateTxInput(utxoOutput, totalNeed)
	txInputs, err := xc.GeneratePureTxInputs(utxoOutput)
	if err != nil {
		log.Printf("GenRealTx GenerateTxInput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxInput err: %v", err)
	}

	//	if deltaTxOutput != nil {
	//		txOutputs = append(txOutputs, deltaTxOutput)
	//	}

	tx := &pb.Transaction{
		Desc:      []byte("Maybe common transfer transaction"),
		Version:   common.TxVersion,
		Coinbase:  false,
		Timestamp: time.Now().UnixNano(),
		TxInputs:  txInputs,
		TxOutputs: txOutputs,
		Initiator: xc.Initiator,
		//		AuthRequire: xc.AuthRequire,
	}

	if len(xc.AuthRequire) != 0 {
		tx.AuthRequire = xc.AuthRequire
	}

	cryptoClient := crypto.GetCryptoClient()

	if len(hdPublicKey) == 0 {
		// 如果不需要HD分层加密功能
		tx.Desc = []byte(xc.Desc)
	} else {
		// 如果需要HD分层加密功能
		cypherText, err := cryptoClient.EncryptByHdKey(hdPublicKey, xc.Desc)
		if err != nil {
			return nil, err
		}

		tx.Desc = []byte(cypherText)

		// 继续组装HDInfo
		originalHash := cryptoClient.HashUsingDoubleSha256([]byte(xc.Desc))

		hdInfo := &pb.HDInfo{
			HdPublicKey:  []byte(hdPublicKey),
			OriginalHash: originalHash,
		}

		tx.HDInfo = hdInfo
	}

	tx.TxInputsExt = response.GetResponse().GetInputs()
	tx.TxOutputsExt = response.GetResponse().GetOutputs()
	tx.ContractRequests = response.GetResponse().GetRequests()

	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()

	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(xc.Account.PrivateKey)
	if err != nil {
		return nil, err
	}

	digestHash, dhErr := txhash.MakeTxDigestHash(tx)
	if dhErr != nil {
		return nil, dhErr
	}

	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)
	if err != nil {
		return nil, err
	}
	signatureInfo := &pb.SignatureInfo{
		PublicKey: xc.Account.PublicKey,
		Sign:      sign,
	}

	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)

	tx.InitiatorSigns = signatureInfos
	if xc.ContractAccount != "" {
		tx.AuthRequireSigns = signatureInfos
	}
	// make txid
	tx.Txid, _ = txhash.MakeTransactionID(tx)
	return tx, nil

}

// ComplianceCheck whether the transaction complies with the rule
func (xc *Xchain) ComplianceCheck(tx *pb.Transaction, fee *pb.Transaction) (
	*pb.SignatureInfo, error) {
	txStatus := &pb.TxStatus{
		Bcname: xc.ChainName,
		Tx:     tx,
	}

	requestData, err := json.Marshal(txStatus)
	if err != nil {
		log.Printf("json encode txStatus failed: %v", err)
		return nil, err
	}

	endorserRequest := &pb.EndorserRequest{
		RequestName: "ComplianceCheck",
		BcName:      xc.ChainName,
		Fee:         fee,
		RequestData: requestData,
	}

	conn, err := grpc.Dial(xc.Cfg.EndorseServiceHost, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("ComplianceCheck connect EndorseServiceHost err: %v", err)
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	c := pb.NewXendorserClient(conn)
	endorserResponse, err := c.EndorserCall(ctx, endorserRequest)
	if err != nil {
		log.Printf("EndorserCall failed and err is: %v", err)
		return nil, fmt.Errorf("EndorserCall error! Response is: %v", err)
	}

	return endorserResponse.GetEndorserSign(), nil
}

// PostTx posttx
func (xc *Xchain) PostTx(tx *pb.Transaction) (string, error) {
	posttx := func(tx *pb.Transaction) error {
		conn, err := grpc.Dial(xc.XchainSer, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
		if err != nil {
			log.Printf("Posttx connect xchain err: %v", err)
			return err
		}
		defer conn.Close()
		ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
		defer cancel()

		txStatus := &pb.TxStatus{
			Bcname: xc.ChainName,
			Status: pb.TransactionStatus_UNCONFIRM,
			Tx:     tx,
			Txid:   tx.Txid,
		}

		c := pb.NewXchainClient(conn)
		res, err := c.PostTx(ctx, txStatus)
		if err != nil {
			return err
		}
		if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
			return fmt.Errorf("Failed to post tx: %s", res.Header.Error.String())
		}

		return nil
	}
	err := posttx(tx)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(tx.Txid), nil
}

// GenCompleteTxAndPost generate comlete tx and post tx
func (xc *Xchain) GenCompleteTxAndPost(preExeResp *pb.PreExecWithSelectUTXOResponse, hdPublicKey string) (string, error) {
	tx := &pb.Transaction{}
	complianceCheckTx := &pb.Transaction{}
	var err error
	// if ComplianceCheck is needed
	// 如果需要进行合规性背书
	if xc.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		// 如果需要进行合规性背书，但不需要支付合规性背书费用
		if xc.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == false {
			tx, err = xc.GenRealTxOnly(preExeResp, hdPublicKey)
			if err != nil {
				log.Printf("GenRealTxOnly failed, err: %v", err)
				return "", err
			}

			complianceCheckTx = nil
		} else {
			// 如果需要进行合规性背书，且需要支付合规性背书费用
			complianceCheckTx, err = xc.GenComplianceCheckTx(preExeResp)
			if err != nil {
				log.Printf("GenCompleteTxAndPost GenComplianceCheckTx failed, err: %v", err)
				return "", err
			}
			log.Printf("ComplianceCheck txid: %v\n", hex.EncodeToString(complianceCheckTx.Txid))

			tx, err = xc.GenRealTx(preExeResp, complianceCheckTx, hdPublicKey)
			if err != nil {
				log.Printf("GenRealTx failed, err: %v", err)
				return "", err
			}
		}

		endorserSign, err := xc.ComplianceCheck(tx, complianceCheckTx)
		if err != nil {
			return "", err
		}

		tx.AuthRequireSigns = append(tx.AuthRequireSigns, endorserSign)

		// 如果是平台发起的转账
		if xc.PlatformAccount != nil {
			cryptoClient := crypto.GetCryptoClient()
			platformPrivateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(xc.PlatformAccount.PrivateKey)
			if err != nil {
				return "", err
			}
			digestHash, dhErr := txhash.MakeTxDigestHash(tx)
			if dhErr != nil {
				return "", dhErr
			}
			platformSign, err := cryptoClient.SignECDSA(platformPrivateKey, digestHash)
			if err != nil {
				return "", err
			}
			platformSignatureInfo := &pb.SignatureInfo{
				PublicKey: xc.PlatformAccount.PublicKey,
				Sign:      platformSign,
			}

			tx.AuthRequireSigns = append(tx.AuthRequireSigns, platformSignatureInfo)
		}

		tx.Txid, _ = txhash.MakeTransactionID(tx)

	} else {
		// only GenRealTx is needed
		// 如果不需要进行合规性背书
		tx, err = xc.GenRealTxOnly(preExeResp, hdPublicKey)
		if err != nil {
			log.Printf("GenRealTxOnly failed, err: %v", err)
			return "", err
		}

	}

	//	txJSON, _ := json.Marshal(tx)
	//	log.Printf("tx is: %s", txJSON)

	return xc.PostTx(tx)

}

// GetBalanceDetail get unfrozen balance and frozen balance
func (xc *Xchain) GetBalanceDetail() (string, error) {
	conn, err := grpc.Dial(xc.XchainSer, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("GetBalance connect xchain err: %v", err)
		return "", err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	tfds := []*pb.TokenFrozenDetails{{Bcname: xc.ChainName}}
	addStatus := &pb.AddressBalanceStatus{
		Address: xc.Account.Address,
		Tfds:    tfds,
	}

	c := pb.NewXchainClient(conn)
	res, err := c.GetBalanceDetail(ctx, addStatus)
	if err != nil {
		return "", err
	}
	if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return "", errors.New(res.Header.Error.String())
	}

	balanceJSON, err := json.Marshal(res.Tfds[0].Tfd)
	return string(balanceJSON), err
}

// QueryTx get tx's status
func (xc *Xchain) QueryTx(txid string) (*pb.TxStatus, error) {
	conn, err := grpc.Dial(xc.XchainSer, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("QueryTx connect xchain err: %v", err)
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	rawTxid, err := hex.DecodeString(txid)
	if err != nil {
		return nil, fmt.Errorf("txid format is wrong: %s", txid)
	}
	txStatus := &pb.TxStatus{
		Bcname: xc.ChainName,
		Txid:   rawTxid,
	}

	c := pb.NewXchainClient(conn)
	res, err := c.QueryTx(ctx, txStatus)
	if err != nil {
		return nil, err
	}
	if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(res.Header.Error.String())
	}
	if res.Tx == nil {
		return nil, common.ErrTxNotFound
	}
	return res, nil
}

// PreExec pre exec
func (xc *Xchain) PreExec() (*pb.InvokeRPCResponse, error) {
	conn, err := grpc.Dial(xc.XchainSer, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("Posttx connect xchain err: %v", err)
		return nil, err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	c := pb.NewXchainClient(conn)
	preExeRPCRes, err := c.PreExec(ctx, xc.InvokeRPCReq)
	if err != nil {
		return nil, err
	}
	for _, res := range preExeRPCRes.GetResponse().GetResponses() {
		if res.Status >= 400 {
			return nil, fmt.Errorf("contract error status:%d message:%s", res.Status, res.Message)
		}
		log.Printf("contract response: %s\n", string(res.Body))
	}
	log.Printf("Gas will cost: %v\n", preExeRPCRes.GetResponse().GetGasUsed())
	return preExeRPCRes, nil
}
