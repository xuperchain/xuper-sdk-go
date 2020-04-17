// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package xchain is related to xchain operation
package xchain

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"time"

	"github.com/xuperchain/xuperchain/core/pb"
	"github.com/xuperchain/xuperchain/core/utxo/txhash"
	"google.golang.org/grpc"

	"strconv"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/crypto"
)

// Xchain xchain struct
type Xchain struct {
	Cfg             *config.CommConfig
	To              string
	Amount          string
	Fee             string
	DescFile        string
	FrozenHeight    int64
	InvokeRPCReq    *pb.InvokeRPCRequest
	PreSelUTXOReq   *pb.PreExecWithSelectUTXORequest
	Initiator       string
	AuthRequire     []string
	Account         *account.Account
	ChainName       string
	XchainSer       string
	ContractAccount string
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

	log.Printf("Fee will cost: %v\n", preExecWithSelectUTXOResponse.GetResponse().GetGasUsed())
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
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJSON([]byte(xc.Account.PrivateKey))
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
			return nil, errors.New("Invalid amount number")
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
	complianceCheckTx *pb.Transaction) (*pb.Transaction, error) {
	txOutputs, err := xc.GenerateTxOutput(xc.To, xc.Amount, xc.Fee)
	if err != nil {
		log.Printf("GenRealTx GenerateTxOutput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxOutput err: %v", err)
	}

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
	amount, ok := big.NewInt(0).SetString(xc.Amount, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	fee, ok := big.NewInt(0).SetString(xc.Fee, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}
	amount.Add(amount, fee)
	totalNeed.Add(totalNeed, amount)
	txInputs, deltaTxOutput, err := xc.GenerateTxInput(utxoOutput, totalNeed)
	if err != nil {
		log.Printf("GenRealTx GenerateTxInput failed.")
		return nil, fmt.Errorf("GenRealTx GenerateTxInput err: %v", err)
	}

	if deltaTxOutput != nil {
		txOutputs = append(txOutputs, deltaTxOutput)
	}

	tx := &pb.Transaction{
		Desc:        []byte("Maybe common transfer transaction"),
		Version:     common.TxVersion,
		Coinbase:    false,
		Timestamp:   time.Now().UnixNano(),
		TxInputs:    txInputs,
		TxOutputs:   txOutputs,
		Initiator:   xc.Initiator,
		AuthRequire: xc.AuthRequire,
	}

	if xc.DescFile != "" {
		tx.Desc, err = ioutil.ReadFile(xc.DescFile)
		if err != nil {
			return nil, err
		}
	}
	tx.TxInputsExt = response.GetResponse().GetInputs()
	tx.TxOutputsExt = response.GetResponse().GetOutputs()
	tx.ContractRequests = response.GetResponse().GetRequests()

	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()
	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJSON([]byte(xc.Account.PrivateKey))
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
func (xc *Xchain) PostTx(tx *pb.Transaction) error {
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

// GenCompleteTxAndPost generate comlete tx and post tx
func (xc *Xchain) GenCompleteTxAndPost(preExeResp *pb.PreExecWithSelectUTXOResponse) (string, error) {
	complianceCheckTx, err := xc.GenComplianceCheckTx(preExeResp)
	if err != nil {
		log.Printf("GenCompleteTxAndPost GenComplianceCheckTx failed, err: %v", err)
		return "", err
	}
	log.Printf("ComplianceCheck txid: %v\n", hex.EncodeToString(complianceCheckTx.Txid))

	tx, err := xc.GenRealTx(preExeResp, complianceCheckTx)
	if err != nil {
		log.Printf("GenRealTx failed, err: %v", err)
		return "", err
	}

	endorserSign, err := xc.ComplianceCheck(tx, complianceCheckTx)
	if err != nil {
		log.Println("call ComplianceCheck failed, err=", err.Error())
		return "", err
	}

	tx.AuthRequireSigns = append(tx.AuthRequireSigns, endorserSign)

	tx.Txid, _ = txhash.MakeTransactionID(tx)

	err = xc.PostTx(tx)
	if err != nil {
		return "", err
	}

	log.Printf("Real txid: %v\n", hex.EncodeToString(tx.Txid))
	return hex.EncodeToString(tx.Txid), nil
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
	log.Printf("Fee will cost: %v\n", preExeRPCRes.GetResponse().GetGasUsed())
	return preExeRPCRes, nil
}
