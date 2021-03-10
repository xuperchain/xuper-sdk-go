// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package xchain is related to xchain operation
package xchain

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/crypto"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/txhash"
)

func initXchain() *Xchain {
	commConfig := config.GetInstance()

	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("RetrieveAccount: %v\n", acc)

	return &Xchain{
		Cfg:       commConfig,
		XchainSer: "127.0.0.1:37801",
		ChainName: "xuper",
		Account:   acc,
	}
}

func TestTransfer(t *testing.T) {
	cli := initXchain()
	amount := "10"
	cli.Fee = "0"
	// to := "Bob"

	invokeRequests := []*pb.InvokeRequest{}
	authRequires := []string{}
	authRequires = append(authRequires, cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      cli.ChainName,
		Requests:    invokeRequests,
		Initiator:   cli.Account.Address,
		AuthRequire: authRequires,
	}

	amountInt64, _ := strconv.ParseInt(amount, 10, 64)
	feeInt64, _ := strconv.ParseInt(cli.Fee, 10, 64)
	needTotalAmount := amountInt64 + int64(cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) + feeInt64

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      cli.ChainName,
		Address:     cli.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}

	cli.PreSelUTXOReq = preSelUTXOReq

	// populates fields
	cli.InvokeRPCReq = invokeRPCReq
	cli.Initiator = cli.Account.Address
	cli.AuthRequire = authRequires

	// test PreExecWithSelecUTXO
	preExeWithSelRes, err := cli.PreExecWithSelecUTXO()
	t.Logf("PreExecWithSelecUTXO: %v, err: %v", preExeWithSelRes, err)
	// test GenComplianceCheckTx
	complianceCheckTx, err := cli.GenComplianceCheckTx(preExeWithSelRes)
	t.Logf("GenComplianceCheckTx: %v, err: %v", complianceCheckTx, err)

	// test GenRealTx
	tx, err := cli.GenRealTx(preExeWithSelRes, complianceCheckTx, "")
	t.Logf("GenRealTx: %v, err: %v", tx, err)

	// test ComplianceCheck
	endorserSign, err := cli.ComplianceCheck(tx, complianceCheckTx)
	t.Logf("ComplianceCheck: %v, err: %v", endorserSign, err)

	tx.AuthRequireSigns = append(tx.AuthRequireSigns, endorserSign)
	tx.Txid, _ = txhash.MakeTransactionID(tx)

	// test PostTx
	_, err = cli.PostTx(tx)
	t.Logf("PostTx: err: %v", err)

	txid := hex.EncodeToString(tx.Txid)
	// test QueryTx
	completeTx, err := cli.QueryTx(txid)
	t.Logf("GenCompleteTxAndPost: %v, err: %v", completeTx, err)
}

func TestTransferV2(t *testing.T) {
	cli := initXchain()

	amount := "100"
	cli.Fee = "0"
	to := "Bob"

	invokeRequests := []*pb.InvokeRequest{}
	authRequires := []string{}
	authRequires = append(authRequires, cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      cli.ChainName,
		Requests:    invokeRequests,
		Initiator:   cli.Account.Address,
		AuthRequire: authRequires,
	}

	amountInt64, _ := strconv.ParseInt(amount, 10, 64)
	feeInt64, _ := strconv.ParseInt(cli.Fee, 10, 64)
	needTotalAmount := amountInt64 + int64(cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) + feeInt64

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      cli.ChainName,
		Address:     cli.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}

	cli.PreSelUTXOReq = preSelUTXOReq

	// populates fields
	cli.InvokeRPCReq = invokeRPCReq
	cli.Initiator = cli.Account.Address
	cli.AuthRequire = authRequires

	// test PreExecWithSelecUTXO
	preExeWithSelRes, err := cli.PreExecWithSelecUTXO()
	t.Logf("PreExecWithSelecUTXO: %v, err: %v", preExeWithSelRes, err)
	// test GenComplianceCheckTx
	complianceCheckTx, err := cli.GenComplianceCheckTx(preExeWithSelRes)

	// test GenerateTxOutput
	txOutputs, err := cli.GenerateTxOutput(to, amount, cli.Fee)
	t.Logf("GenerateTxOutput: %v, err: %v", txOutputs, err)

	utxolist := []*pb.Utxo{}
	totalSelected := big.NewInt(0)
	for index, txOutput := range complianceCheckTx.TxOutputs {
		if string(txOutput.ToAddr) == cli.Initiator {
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
	amountB, _ := big.NewInt(0).SetString(amount, 10)
	fee, _ := big.NewInt(0).SetString(cli.Fee, 10)
	amountB.Add(amountB, fee)
	totalNeed.Add(totalNeed, amountB)

	// test GenerateTxInput
	txInputs, deltaTxOutput, err := cli.GenerateTxInput(utxoOutput, totalNeed)
	t.Logf("GenerateTxInput: input: %v, output: %v, err: %v", txInputs, deltaTxOutput, err)

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
		Initiator:   cli.Initiator,
		AuthRequire: cli.AuthRequire,
	}

	tx.TxInputsExt = preExeWithSelRes.GetResponse().GetInputs()
	tx.TxOutputsExt = preExeWithSelRes.GetResponse().GetOutputs()
	tx.ContractRequests = preExeWithSelRes.GetResponse().GetRequests()
	common.SetSeed()
	tx.Nonce = common.GetNonce()
	cryptoClient := crypto.GetCryptoClient()
	privateKey, _ := cryptoClient.GetEcdsaPrivateKeyFromJsonStr((cli.Account.PrivateKey))

	digestHash, _ := txhash.MakeTxDigestHash(tx)
	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)
	signatureInfo := &pb.SignatureInfo{
		PublicKey: cli.Account.PublicKey,
		Sign:      sign,
	}
	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)
	tx.InitiatorSigns = signatureInfos
	if cli.ContractAccount != "" {
		tx.AuthRequireSigns = signatureInfos
	}
	// make txid
	tx.Txid, _ = txhash.MakeTransactionID(tx)

	// test ComplianceCheck
	endorserSign, err := cli.ComplianceCheck(tx, complianceCheckTx)
	t.Logf("ComplianceCheck: %v, err: %v", endorserSign, err)

	tx.AuthRequireSigns = append(tx.AuthRequireSigns, endorserSign)
	tx.Txid, _ = txhash.MakeTransactionID(tx)

	// test PostTx
	_, err = cli.PostTx(tx)
	t.Logf("PostTx: err: %v", err)

	txid := hex.EncodeToString(tx.Txid)
	// test QueryTx
	completeTx, err := cli.QueryTx(txid)
	t.Logf("GenCompleteTxAndPost: %v, err: %v", completeTx, err)

}

func TestTransferV3(t *testing.T) {
	cli := initXchain()

	amount := "10"
	cli.Fee = "0"
	// cli.To = "Bob"

	invokeRequests := []*pb.InvokeRequest{}
	authRequires := []string{}
	authRequires = append(authRequires, cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      cli.ChainName,
		Requests:    invokeRequests,
		Initiator:   cli.Account.Address,
		AuthRequire: authRequires,
	}
	amountInt64, _ := strconv.ParseInt(amount, 10, 64)
	feeInt64, _ := strconv.ParseInt(cli.Fee, 10, 64)
	needTotalAmount := amountInt64 + int64(cli.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee) + feeInt64
	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      cli.ChainName,
		Address:     cli.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}
	cli.PreSelUTXOReq = preSelUTXOReq

	// populates fields
	cli.InvokeRPCReq = invokeRPCReq
	cli.Initiator = cli.Account.Address
	cli.AuthRequire = authRequires

	// test PreExecWithSelecUTXO
	preExeWithSelRes, err := cli.PreExecWithSelecUTXO()
	t.Logf("PreExecWithSelecUTXO: %v, err: %v", preExeWithSelRes, err)

	// test GenCompleteTxAndPost
	tx, err := cli.GenCompleteTxAndPost(preExeWithSelRes, "")
	t.Logf("GenCompleteTxAndPost: %v, err: %v", tx, err)

	// test QueryTx
	completeTx, err := cli.QueryTx(tx)
	t.Logf("GenCompleteTxAndPost: %v, err: %v", completeTx, err)
}

func TestGetBalanceDetail(t *testing.T) {
	cli := initXchain()

	testCase := []string{
		"Sw5kwvaf3PAwozXxMdFuBrd9UiqXuXhVF",
		"XC8888888888888888@xuper",
		"xxxx",
	}

	for _, arg := range testCase {
		cli.Account = &account.Account{
			Address: arg,
		}
		balance, err := cli.GetBalanceDetail()
		t.Logf("GetBalanceDetail address: %v, err: %v", balance, err)
	}
}

func TestQueryTx(t *testing.T) {
	acc, err := account.RetrieveAccount("江 西 伏 物 十 勘 峡 环 初 至 赏 给", 1)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("RetrieveAccount: %v\n", acc)

	cli := initXchain()

	testCase := []struct {
		txid string
	}{
		{
			txid: "3a78d06dd39b814af113dbdc15239e675846ec927106d50153665c273f51001e",
		},
		{
			txid: "",
		},
		{
			txid: "fdsfdsa",
		},
	}

	for _, arg := range testCase {
		tx, err := cli.QueryTx(arg.txid)
		t.Logf("Querytx tx: %v, err: %v", tx, err)
	}
}
