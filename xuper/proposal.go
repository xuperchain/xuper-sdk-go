package xuper

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuper-sdk-go/v2/crypto"
	"github.com/xuperchain/xuperchain/core/pb"
)

const (
	defaultChainName = "xuper"
)

type Proposal struct {
	xclient *XClient
	request *Request

	preResp           *pb.PreExecWithSelectUTXOResponse
	tx                *Transaction
	complianceCheckTx *pb.Transaction
}

func (p *Proposal) Build() (*Transaction, error) {
	preResp, err := p.PreExecWithSelectUtxo()
	if err != nil {
		return nil, err
	}

	tx, err := p.GenCompleteTx(preResp)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Proposal) PreExecWithSelectUtxo() (*pb.PreExecWithSelectUTXOResponse, error) {

	req, err := p.genPreExecUtxoRequest()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	preExecWithSelectUTXOResponse := new(pb.PreExecWithSelectUTXOResponse)

	if config.GetInstance().ComplianceCheck.IsNeedComplianceCheck {
		requestData, err := json.Marshal(req)
		if err != nil {
			return nil, err
		}
		endorserRequest := &pb.EndorserRequest{
			RequestName: "PreExecWithFee",
			BcName:      req.Bcname,
			RequestData: requestData,
		}
		c := *p.xclient.ec
		endorserResponse, err := c.EndorserCall(ctx, endorserRequest)
		if err != nil {
			return nil, errors.Wrap(err, "EndorserCall failed") // todo error wrap?
		}
		responseData := endorserResponse.ResponseData
		err = json.Unmarshal(responseData, preExecWithSelectUTXOResponse)
		if err != nil {
			return nil, err
		}
	} else {
		c := *p.xclient.xc
		var err error
		preExecWithSelectUTXOResponse, err = c.PreExecWithSelectUTXO(ctx, req)
		if err != nil {
			return nil, errors.Wrap(err, "PreExecWithSelectUTXO failed") // todo error wrap?
		}
	}

	for _, res := range preExecWithSelectUTXOResponse.GetResponse().GetResponses() {
		if res.Status >= 400 {
			return nil, fmt.Errorf("contract invoke error status:%d message:%s", res.Status, res.Message)
		}
	}

	return preExecWithSelectUTXOResponse, nil
}

func (p *Proposal) GenCompleteTx(preResp *pb.PreExecWithSelectUTXOResponse) (*Transaction, error) {
	var (
		tx  *pb.Transaction
		err error
	)

	if config.GetInstance().ComplianceCheck.IsNeedComplianceCheck {
		tx, err = p.genTxWithComplianceCheck(preResp)
		if err != nil {
			return nil, err
		}
	} else {
		tx, err = p.genTx(preResp, nil)
		if err != nil {
			return nil, err
		}
	}

	var ContractResponse *pb.ContractResponse
	if len(preResp.GetResponse().GetResponses()) != 0 {
		// 如果没有背书，那么一个合约调用应该有一个 response。
		// 有背书或者有 reserved contract 时，会有多个 response，最后一个 response 为本次交易的合约执行结果。
		// server 端实现代码在 xuperchain 项目：core/utxo/utxo.go:PreExec 接口。
		ContractResponse = preResp.GetResponse().GetResponses()[len(preResp.GetResponse().GetResponses())-1]
	}

	transaction := &Transaction{
		Tx: tx,
		// ComplianceCheckTx: nil,
		ContractResponse: ContractResponse,
		Fee:              p.request.opt.fee,
		GasUsed:          preResp.GetResponse().GetGasUsed(),
		DigestHash:       nil, // todo
	}

	return transaction, nil
}

func (p *Proposal) genTxWithComplianceCheck(preResp *pb.PreExecWithSelectUTXOResponse) (*pb.Transaction, error) {
	var (
		complianceCheckTx *pb.Transaction
		err               error
	)

	if config.GetInstance().ComplianceCheck.IsNeedComplianceCheckFee {
		complianceCheckTx, err = p.genComplianceCheckTx(preResp)
		if err != nil {
			return nil, err
		}
	}

	tx, err := p.genTx(preResp, complianceCheckTx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Proposal) genComplianceCheckTx(preResp *pb.PreExecWithSelectUTXOResponse) (*pb.Transaction, error) {
	complianceCheckFee := config.GetInstance().ComplianceCheck.ComplianceCheckEndorseServiceFee
	complianceCheckFeeAddr := config.GetInstance().ComplianceCheck.ComplianceCheckEndorseServiceAddr
	utxoOutput := preResp.GetUtxoOutput()

	checkTxOutput, err := p.generateComplianceCheckTxOutput(complianceCheckFeeAddr, strconv.Itoa(complianceCheckFee))
	if err != nil {
		return nil, err
	}

	complianceCheckFeeBigInt := new(big.Int).SetInt64(int64(complianceCheckFee)) // 不会有这么多的手续费，所以暂不考虑溢出。
	txInputs, deltaTxOutput, err := p.generateComplianceCheckTxInput(utxoOutput, complianceCheckFeeBigInt)

	if deltaTxOutput != nil {
		checkTxOutput = append(checkTxOutput, deltaTxOutput)
	}

	initiator := p.request.initiatorAccount
	tx := &pb.Transaction{
		Desc:      []byte(""),
		Version:   common.TxVersion,
		Coinbase:  false,
		Timestamp: time.Now().UnixNano(),
		TxInputs:  txInputs,
		TxOutputs: checkTxOutput,
		Initiator: initiator.Address,
	}

	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()

	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(initiator.PrivateKey)
	if err != nil {
		return nil, err
	}
	// digestHash, dhErr := txhash.MakeTxDigestHash(tx)
	// if dhErr != nil {
	// 	return nil, dhErr
	// }

	// sign, err := cryptoClient.SignECDSA(privateKey, digestHash)
	sign, err := cryptoClient.SignECDSA(privateKey, nil) // todo

	signatureInfo := &pb.SignatureInfo{
		PublicKey: initiator.PublicKey,
		Sign:      sign,
	}

	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)

	tx.InitiatorSigns = signatureInfos

	// make txid
	// tx.Txid, _ = txhash.MakeTransactionID(tx) // todo
	return nil, nil
}

func (p *Proposal) populateTx(desc string, output []*pb.TxOutput, input []*pb.TxInput) (*pb.Transaction, error) {

	return nil, nil
}

// generateTxOutput generate txoutput part
func (p *Proposal) generateComplianceCheckTxOutput(to, amount string) ([]*pb.TxOutput, error) {
	accounts := []*pb.TxDataAccount{}
	if to != "" {
		account := &pb.TxDataAccount{
			Address:      to,
			Amount:       amount,
			FrozenHeight: 0,
		}
		accounts = append(accounts, account)
	}
	// if fee != "0" {
	// 	feeAccount := &pb.TxDataAccount{
	// 		Address: "$",
	// 		Amount:  fee,
	// 	}
	// 	accounts = append(accounts, feeAccount)
	// }

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

// generateTxInput generate txinput part
func (p *Proposal) generateComplianceCheckTxInput(utxoOutputs *pb.UtxoOutput, totalNeed *big.Int) ([]*pb.TxInput, *pb.TxOutput, error) {
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
		return nil, nil, fmt.Errorf("Invalid utxoOutputs.TotalSelected: %s", utxoOutputs.TotalSelected)
	}

	// input > output, generate output-input to me
	if utxoTotal.Cmp(totalNeed) > 0 {
		delta := utxoTotal.Sub(utxoTotal, totalNeed)
		txOutput = &pb.TxOutput{
			ToAddr: []byte(p.request.initiatorAccount.Address),
			Amount: delta.Bytes(),
		}
	}

	return txInputs, txOutput, nil
}

func (p *Proposal) genTx(preResp *pb.PreExecWithSelectUTXOResponse, lastTx *pb.Transaction) (*pb.Transaction, error) {
	utxoOutput := &pb.UtxoOutput{}
	totalSelected := big.NewInt(0)
	totalNeed := big.NewInt(0)
	utxolist := []*pb.Utxo{}

	if lastTx != nil {
		for index, txOutput := range lastTx.TxOutputs {
			if string(txOutput.ToAddr) == p.request.initiatorAccount.Address {
				utxo := &pb.Utxo{
					Amount:    txOutput.Amount,
					ToAddr:    txOutput.ToAddr,
					RefTxid:   lastTx.Txid,
					RefOffset: int32(index),
				}
				utxolist = append(utxolist, utxo)

				utxoAmount := big.NewInt(0).SetBytes(utxo.Amount)
				totalSelected.Add(totalSelected, utxoAmount)
			}
		}
	} else {
		if preResp.UtxoOutput != nil {
			utxoOutput.UtxoList = preResp.UtxoOutput.UtxoList
			utxoOutput.TotalSelected = preResp.UtxoOutput.TotalSelected
		}
	}

	amount := big.NewInt(0)

	// amount
	if p.request.opt.contractInvokeAmount != "" {
		invokeAmount, ok := big.NewInt(0).SetString(p.request.opt.contractInvokeAmount, 10) // todo
		if !ok {
			return nil, common.ErrInvalidAmount
		}
		amount.Add(amount, invokeAmount)
	}

	if p.request.transferAmount != "" {
		transferAmount, ok := big.NewInt(0).SetString(p.request.transferAmount, 10) // todo
		if !ok {
			return nil, common.ErrInvalidAmount
		}
		amount.Add(amount, transferAmount)
	}

	// fee
	if p.request.opt.fee != "" {
		fee, ok := big.NewInt(0).SetString(p.request.opt.fee, 10)
		if !ok {
			return nil, common.ErrInvalidAmount
		}
		amount.Add(amount, fee)
	}

	// gas
	gasUsed := big.NewInt(preResp.GetResponse().GetGasUsed())
	amount.Add(amount, gasUsed)

	// total
	totalNeed.Add(totalNeed, amount)

	selfAmount := totalSelected.Sub(totalSelected, totalNeed)

	txOutputs, err := p.generateMultiTxOutputs(selfAmount.String(), gasUsed)
	if err != nil {
		return nil, err
	}

	txInputs := p.genPureTxInputs(utxoOutput)

	tx := &pb.Transaction{
		Desc:             []byte(p.request.opt.desc),
		Version:          common.TxVersion,
		Coinbase:         false,
		Timestamp:        time.Now().UnixNano(),
		TxInputs:         txInputs,
		TxOutputs:        txOutputs,
		Initiator:        p.request.initiatorAccount.Address,
		AuthRequire:      []string{p.request.initiatorAccount.GetAuthRequire()},
		TxInputsExt:      preResp.GetResponse().GetInputs(),
		TxOutputsExt:     preResp.GetResponse().GetOutputs(),
		ContractRequests: preResp.GetResponse().GetRequests(),
	}
	err = common.SetSeed()
	if err != nil {
		return nil, err
	}
	tx.Nonce = common.GetNonce()

	// todo signature
	return tx, nil
}

func (p *Proposal) genPureTxInputs(utxoOutputs *pb.UtxoOutput) []*pb.TxInput {
	var txInputs []*pb.TxInput
	for _, utxo := range utxoOutputs.UtxoList {
		txInput := &pb.TxInput{}
		txInput.RefTxid = utxo.RefTxid
		txInput.RefOffset = utxo.RefOffset
		txInput.FromAddr = utxo.ToAddr
		txInput.Amount = utxo.Amount
		txInputs = append(txInputs, txInput)
	}

	return txInputs
}

func (p *Proposal) genPreExecUtxoRequest() (*pb.PreExecWithSelectUTXORequest, error) {
	utxoAddr := p.request.initiatorAccount.Address
	if p.request.initiatorAccount.HasContractAccount() {
		utxoAddr = p.request.initiatorAccount.GetContractAccount()
	}

	totalAmount, err := p.calcTotalAmount()
	if err != nil {
		return nil, err
	}

	invokeRPCReq, err := p.genInvokeRPCRequest()
	if err != nil {
		return nil, err
	}

	req := &pb.PreExecWithSelectUTXORequest{
		Bcname:      p.getChainName(),
		Address:     utxoAddr,
		TotalAmount: totalAmount,
		Request:     invokeRPCReq,
	}

	return req, nil
}

func (p *Proposal) generateMultiTxOutputs(selfAmount string, gasUsed *big.Int) ([]*pb.TxOutput, error) {
	realSelfAmount, ok := new(big.Int).SetString(selfAmount, 10)
	if !ok {
		return nil, common.ErrInvalidAmount
	}

	var txOutputs []*pb.TxOutput
	req := p.request

	// 1. transfer
	if req.transferTo != "" {
		txOutput, err := p.makeTxOutput(req.transferTo, req.transferAmount)
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// 2. transfer to contract
	if req.opt.contractInvokeAmount != "" {
		txOutput, err := p.makeTxOutput(req.contractName, req.opt.contractInvokeAmount)
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// 3. self
	if realSelfAmount.Cmp(big.NewInt(0)) > 0 {
		txOutput, err := p.makeTxOutput(req.initiatorAccount.Address, selfAmount)
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// 4. fee
	if req.opt.fee != "" && req.opt.fee != "0" {
		txOutput, err := p.makeTxOutput("$", req.opt.fee)
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// 5. gasUsed
	if gasUsed.Cmp(big.NewInt(0)) > 0 {
		txOutput, err := p.makeTxOutput("$", gasUsed.String())
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	return txOutputs, nil
}

func (p *Proposal) makeTxOutput(addr, amount string) (*pb.TxOutput, error) {
	txOutput := new(pb.TxOutput)
	txOutput.ToAddr = []byte(addr)
	realToAmount, isSuccess := new(big.Int).SetString(amount, 10)
	if isSuccess != true {
		return nil, common.ErrInvalidAmount
	}
	txOutput.Amount = realToAmount.Bytes()
	return txOutput, nil
}

func (p *Proposal) getChainName() string {
	chainName := defaultChainName
	if p.request.opt.bcname != "" {
		chainName = p.request.opt.bcname
	}
	return chainName
}

func (p *Proposal) getInitiator() string {
	initiator := p.request.initiatorAccount.Address
	if p.request.initiatorAccount.HasContractAccount() {
		initiator = p.request.initiatorAccount.GetContractAccount()
	}
	return initiator
}

func (p *Proposal) genInvokeRequests() ([]*pb.InvokeRequest, error) {
	r := p.request
	invokeReq := &pb.InvokeRequest{
		ModuleName:   r.module,
		ContractName: r.contractName,
		MethodName:   r.methodName,
		Args:         r.args,
		Amount:       r.opt.contractInvokeAmount,
	}
	return []*pb.InvokeRequest{invokeReq}, nil
}

func (p *Proposal) genInvokeRPCRequest() (*pb.InvokeRPCRequest, error) {
	invokeRequests, err := p.genInvokeRequests()
	if err != nil {
		return nil, err
	}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      p.getChainName(),
		Requests:    invokeRequests,
		Initiator:   p.getInitiator(),
		AuthRequire: []string{p.request.initiatorAccount.GetAuthRequire()},
	}

	return invokeRPCReq, nil
}

func (p *Proposal) calcTotalAmount() (int64, error) {
	var totalAmount int64
	req := p.request
	if req.transferAmount != "" {
		if amount, err := strconv.ParseInt(req.transferAmount, 10, 64); err == nil {
			totalAmount += amount
		} else {
			return 0, err
		}
	}

	if req.opt.fee != "" {
		if amount, err := strconv.ParseInt(req.opt.fee, 10, 64); err == nil {
			totalAmount += amount
		} else {
			return 0, err
		}
	}

	if req.opt.contractInvokeAmount != "" {
		if amount, err := strconv.ParseInt(req.opt.contractInvokeAmount, 10, 64); err == nil {
			totalAmount += amount
		} else {
			return 0, err
		}
	}

	// endorser logic
	cfg := config.GetInstance()
	if cfg.ComplianceCheck.IsNeedComplianceCheck &&
		cfg.ComplianceCheck.IsNeedComplianceCheckFee {

		totalAmount += int64(cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
	}

	return totalAmount, nil
}
