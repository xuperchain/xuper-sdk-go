package xuper

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuper-sdk-go/v2/crypto"
	"github.com/xuperchain/xuperchain/service/pb"
)

const (
	defaultChainName = "xuper"
)

// Proposal 代表单个请求，构造交易，但不 post。
type Proposal struct {
	xclient *XClient
	request *Request

	preResp           *pb.PreExecWithSelectUTXOResponse
	feePreResp        *pb.UtxoOutput
	tx                *Transaction
	complianceCheckTx *pb.Transaction

	cfg *config.CommConfig

	txVersion int32
}

// NewProposal new Proposal instance.
func NewProposal(xclient *XClient, request *Request, cfg *config.CommConfig) (*Proposal, error) {
	if xclient == nil || request == nil || cfg == nil {
		return nil, errors.New("new proposal failed, parameters can not be nil")
	}

	// 开放网络交易版本根据配置文件来，非开放网络交易使用 common 中的 TxVersion 也就是版本3.
	v := int32(common.TxVersion)
	if cfg.ComplianceCheck.IsNeedComplianceCheck {
		v = cfg.TxVersion
	}

	return &Proposal{
		xclient:   xclient,
		request:   request,
		cfg:       cfg,
		txVersion: v,
	}, nil
}

// Build 发起预执行，构造交易。
func (p *Proposal) Build() (*Transaction, error) {
	err := p.PreExecWithSelectUtxo() // T_T!，开放网络所有交易都是通过 AK 支付手续费，除了开放网络，其他的根据是否设置了合约账户，以及是否只是合约账户支付手续费来判断。
	if err != nil {
		return nil, err
	}

	tx, err := p.GenCompleteTx()
	if err != nil {
		return nil, err
	}
	p.tx = tx

	return tx, nil
}

// PreExecWithSelectUtxo 预执行并选择 utxo，如果有背书则调用 EndorserCall。
func (p *Proposal) PreExecWithSelectUtxo() error {

	req, err := p.genPreExecUtxoRequest()
	if err != nil {
		return err
	}

	ctx := context.Background()
	preExecWithSelectUTXOResponse := new(pb.PreExecWithSelectUTXOResponse)

	if p.cfg.ComplianceCheck.IsNeedComplianceCheck {
		requestData, err := json.Marshal(req)
		if err != nil {
			return err
		}
		endorserRequest := &pb.EndorserRequest{
			RequestName: "PreExecWithFee",
			BcName:      req.Bcname,
			RequestData: requestData,
		}
		c := p.xclient.ec
		endorserResponse, err := c.EndorserCall(ctx, endorserRequest)
		if err != nil {
			return errors.Wrap(err, "EndorserCall PreExecWithFee failed")
		}
		responseData := endorserResponse.ResponseData
		err = json.Unmarshal(responseData, preExecWithSelectUTXOResponse)
		if err != nil {
			return err
		}
	} else {
		c := p.xclient.xc
		var err error
		preExecWithSelectUTXOResponse, err = c.PreExecWithSelectUTXO(ctx, req)
		if err != nil {
			return errors.Wrap(err, "PreExecWithSelectUTXO failed")
		}

		// AK 发起交易，仅使用合约账户支付手续费时，需要选择 utxo。
		if p.request.opt.onlyFeeFromAccount {
			amount, ok := big.NewInt(0).SetString(p.request.opt.fee, 10)
			if !ok {
				return errors.Wrap(common.ErrInvalidAmount, "invalid request fee")
			}
			amount.Add(amount, big.NewInt(preExecWithSelectUTXOResponse.GetResponse().GetGasUsed()))

			feeReq := p.genSelectUtxoRequest(p.request.initiatorAccount.GetContractAccount(), amount.String())

			p.feePreResp, err = c.SelectUTXO(ctx, feeReq)
			if err != nil {
				return errors.Wrap(err, "SelectUTXO from contract account failed")
			}
		}
	}

	for _, res := range preExecWithSelectUTXOResponse.GetResponse().GetResponses() {
		if res.Status >= 400 {
			return fmt.Errorf("contract invoke error status:%d message:%s", res.Status, res.Message)
		}
	}

	p.preResp = preExecWithSelectUTXOResponse
	return nil
}

// GenCompleteTx 根据预执行结果构造完整的交易。
func (p *Proposal) GenCompleteTx() (*Transaction, error) {
	var (
		tx         *pb.Transaction
		digestHash []byte
		err        error

		preResp = p.preResp
	)

	// public method should check proposal's preResp.
	if preResp == nil {
		return nil, errors.New("proposal preResp can not be nil")
	}

	if p.cfg.ComplianceCheck.IsNeedComplianceCheck {
		tx, err = p.genTxWithComplianceCheck()
		if err != nil {
			return nil, err
		}
	} else {
		tx, err = p.genTx()
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

	// initiator sign tx and calc tx ID.
	digestHash, err = p.signTx(tx)
	if err != nil {
		return nil, err
	}

	transaction := &Transaction{
		Tx:               tx,
		ContractResponse: ContractResponse,
		Bcname:           p.getChainName(),
		Fee:              p.request.opt.fee,
		GasUsed:          preResp.GetResponse().GetGasUsed(),
		DigestHash:       digestHash,
	}

	return transaction, nil
}

func (p *Proposal) genTxWithComplianceCheck() (*pb.Transaction, error) {
	var (
		complianceCheckTx *pb.Transaction
		err               error
	)

	if p.cfg.ComplianceCheck.IsNeedComplianceCheckFee {
		complianceCheckTx, err = p.genComplianceCheckTx()
		if err != nil {
			return nil, err
		}
		p.complianceCheckTx = complianceCheckTx
	}

	tx, err := p.genTx()
	if err != nil {
		return nil, err
	}

	endorserSign, err := p.complianceCheck(tx)
	if err != nil {
		return nil, err
	}
	tx.AuthRequireSigns = append(tx.AuthRequireSigns, endorserSign)

	return tx, nil
}

func (p *Proposal) complianceCheck(tx *pb.Transaction) (*pb.SignatureInfo, error) {
	ctx := context.Background()
	txStatus := &pb.TxStatus{
		Bcname: p.getChainName(),
		Tx:     tx,
	}

	requestData, err := json.Marshal(txStatus)
	if err != nil {
		log.Printf("json encode txStatus failed: %v", err)
		return nil, err
	}

	endorserRequest := &pb.EndorserRequest{
		RequestName: "ComplianceCheck",
		BcName:      p.getChainName(),
		Fee:         p.complianceCheckTx,
		RequestData: requestData,
	}

	endorserResponse, err := p.xclient.ec.EndorserCall(ctx, endorserRequest)
	if err != nil {
		return nil, errors.Wrap(err, "EndorserCall ComplianceCheck failed")
	}
	return endorserResponse.GetEndorserSign(), nil
}

func (p *Proposal) genComplianceCheckTx() (*pb.Transaction, error) {
	complianceCheckFee := p.cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee
	complianceCheckFeeAddr := p.cfg.ComplianceCheck.ComplianceCheckEndorseServiceFeeAddr
	utxoOutput := p.preResp.GetUtxoOutput()

	checkTxOutput, err := p.generateComplianceCheckTxOutput(complianceCheckFeeAddr, strconv.Itoa(complianceCheckFee))
	if err != nil {
		return nil, err
	}

	complianceCheckFeeBigInt := new(big.Int).SetInt64(int64(complianceCheckFee))
	txInputs, deltaTxOutput, err := p.generateComplianceCheckTxInput(utxoOutput, complianceCheckFeeBigInt)

	if deltaTxOutput != nil {
		checkTxOutput = append(checkTxOutput, deltaTxOutput)
	}

	err = common.SetSeed()
	if err != nil {
		return nil, errors.Wrap(err, "Set seed failed")
	}

	tx := &pb.Transaction{
		Desc:        []byte(""),
		Version:     p.txVersion,
		Coinbase:    false,
		Nonce:       common.GetNonce(),
		Timestamp:   time.Now().UnixNano(),
		TxInputs:    txInputs,
		TxOutputs:   checkTxOutput,
		Initiator:   p.getInitiator(),
		AuthRequire: []string{p.getInitiator()},
	}

	// initiator sign tx and calc tx ID.
	_, err = p.signTx(tx)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (p *Proposal) signTx(tx *pb.Transaction) ([]byte, error) {
	initiator := p.request.initiatorAccount

	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(initiator.PrivateKey)
	if err != nil {
		return nil, err
	}

	digestHash, err := common.MakeTxDigestHash(tx)
	if err != nil {
		return nil, err
	}

	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)

	signatureInfo := &pb.SignatureInfo{
		PublicKey: initiator.PublicKey,
		Sign:      sign,
	}

	var signatureInfos []*pb.SignatureInfo
	signatureInfos = append(signatureInfos, signatureInfo)

	tx.InitiatorSigns = signatureInfos

	if len(tx.GetAuthRequireSigns()) == 0 {
		tx.AuthRequireSigns = signatureInfos
	} else {
		tx.AuthRequireSigns = append(tx.AuthRequireSigns, signatureInfos...)
	}

	// make txid
	tx.Txid, err = common.MakeTransactionID(tx)
	if err != nil {
		return nil, errors.Wrap(err, "Make transaction ID failed.")
	}

	return digestHash, nil
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
	for _, utxo := range utxoOutputs.GetUtxoList() {
		txInput := &pb.TxInput{}
		txInput.RefTxid = utxo.RefTxid
		txInput.RefOffset = utxo.RefOffset
		txInput.FromAddr = utxo.ToAddr
		txInput.Amount = utxo.Amount
		txInputs = append(txInputs, txInput)
	}

	utxoTotal, ok := big.NewInt(0).SetString(utxoOutputs.GetTotalSelected(), 10)
	if !ok {
		return nil, nil, fmt.Errorf("Invalid utxoOutputs.TotalSelected: %s", utxoOutputs.GetTotalSelected())
	}

	// input > output, generate output-input to me
	if utxoTotal.Cmp(totalNeed) > 0 {
		delta := utxoTotal.Sub(utxoTotal, totalNeed)
		txOutput = &pb.TxOutput{
			ToAddr: []byte(p.getInitiator()),
			Amount: delta.Bytes(),
		}
	}

	return txInputs, txOutput, nil
}

func (p *Proposal) calcSelfAmount(totalSelected *big.Int) (string, error) {
	totalNeed := big.NewInt(0)
	amount := big.NewInt(0)
	preResp := p.preResp

	// amount
	if p.request.opt.contractInvokeAmount != "" {
		invokeAmount, ok := big.NewInt(0).SetString(p.request.opt.contractInvokeAmount, 10)
		if !ok {
			return "", common.ErrInvalidAmount
		}
		amount.Add(amount, invokeAmount)
	}

	if p.request.transferAmount != "" {
		transferAmount, ok := big.NewInt(0).SetString(p.request.transferAmount, 10)
		if !ok {
			return "", common.ErrInvalidAmount
		}
		amount.Add(amount, transferAmount)
	}

	// fee
	if !p.request.opt.onlyFeeFromAccount {
		if p.request.opt.fee != "" {
			fee, ok := big.NewInt(0).SetString(p.request.opt.fee, 10)
			if !ok {
				return "", common.ErrInvalidAmount
			}
			amount.Add(amount, fee)
		}

		// gas
		gasUsed := big.NewInt(preResp.GetResponse().GetGasUsed())
		amount.Add(amount, gasUsed)
	}

	// total
	totalNeed.Add(totalNeed, amount)

	selfAmount := totalSelected.Sub(totalSelected, totalNeed)

	return selfAmount.String(), nil
}

func (p *Proposal) genTx() (*pb.Transaction, error) {
	utxoOutput := &pb.UtxoOutput{}
	totalSelected := big.NewInt(0)
	preResp := p.preResp

	utxolist := []*pb.Utxo{}

	if p.complianceCheckTx != nil {
		for index, txOutput := range p.complianceCheckTx.TxOutputs {
			if string(txOutput.ToAddr) == p.getInitiator() {
				utxo := &pb.Utxo{
					Amount:    txOutput.Amount,
					ToAddr:    txOutput.ToAddr,
					RefTxid:   p.complianceCheckTx.Txid,
					RefOffset: int32(index),
				}
				utxolist = append(utxolist, utxo)

				utxoAmount := big.NewInt(0).SetBytes(utxo.Amount)
				totalSelected.Add(totalSelected, utxoAmount)
			}
		}
		utxoOutput = &pb.UtxoOutput{
			UtxoList:      utxolist,
			TotalSelected: totalSelected.String(),
		}

	} else {
		if preResp.UtxoOutput != nil {
			utxoOutput.UtxoList = preResp.GetUtxoOutput().GetUtxoList()
			utxoOutput.TotalSelected = preResp.GetUtxoOutput().GetTotalSelected()

			var ok bool
			totalSelected, ok = big.NewInt(0).SetString(preResp.GetUtxoOutput().GetTotalSelected(), 10)
			if !ok {
				return nil, common.ErrInvalidAmount
			}
		}

		// fee from account
		if p.feePreResp != nil {
			utxoOutput.UtxoList = append(utxoOutput.UtxoList, p.feePreResp.GetUtxoList()...)
		}
	}

	selfAmount, err := p.calcSelfAmount(totalSelected)

	txOutputs, err := p.generateMultiTxOutputs(selfAmount, big.NewInt(preResp.GetResponse().GetGasUsed()))
	if err != nil {
		return nil, err
	}

	txInputs := p.genPureTxInputs(utxoOutput)

	authRequire := make([]string, 0, 1)
	if p.complianceCheckTx != nil {
		authRequire = append(authRequire, p.cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	authRequire = append(authRequire, p.request.initiatorAccount.GetAuthRequire())

	if len(p.request.opt.otherAuthRequire) > 0 {
		authRequire = append(authRequire, p.request.opt.otherAuthRequire...)
	}

	err = common.SetSeed()
	if err != nil {
		return nil, errors.Wrap(err, "Set seed failed")
	}

	tx := &pb.Transaction{
		Desc:             []byte(p.request.opt.desc),
		Version:          p.txVersion,
		Coinbase:         false,
		Nonce:            common.GetNonce(),
		Timestamp:        time.Now().UnixNano(),
		TxInputs:         txInputs,
		TxOutputs:        txOutputs,
		Initiator:        p.getInitiator(),
		AuthRequire:      authRequire,
		TxInputsExt:      preResp.GetResponse().GetInputs(),
		TxOutputsExt:     preResp.GetResponse().GetOutputs(),
		ContractRequests: preResp.GetResponse().GetRequests(),
	}

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

func (p *Proposal) genSelectUtxoRequest(address, amount string) *pb.UtxoInput {
	return &pb.UtxoInput{
		Bcname:    p.getChainName(),
		Address:   address,
		TotalNeed: amount,
	}
}

func (p *Proposal) genPreExecUtxoRequest() (*pb.PreExecWithSelectUTXORequest, error) {
	utxoAddr := p.getInitiator()

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
		txOutput, err := p.makeTxOutput(p.getInitiator(), selfAmount)
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// 4. fee & gasUsed
	feeOutput, err := p.makeFeeTxOutput()
	if err != nil {
		return nil, err
	}
	txOutputs = append(txOutputs, feeOutput...)

	return txOutputs, nil
}

func (p *Proposal) makeFeeTxOutput() ([]*pb.TxOutput, error) {
	txOutputs := make([]*pb.TxOutput, 0, 1)
	fee, err := p.calcAllFee()
	if err != nil {
		return nil, err
	}

	// no gasUsed & fee.
	if fee.Cmp(big.NewInt(0)) <= 0 {
		return txOutputs, nil
	}

	if fee.Cmp(big.NewInt(0)) > 0 {
		txOutput, err := p.makeTxOutput("$", fee.String())
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	// fee from contract account, calc account self output.
	if p.request.opt.onlyFeeFromAccount && p.feePreResp != nil {
		total, ok := big.NewInt(0).SetString(p.feePreResp.GetTotalSelected(), 10)
		if !ok {
			return nil, errors.New("invalid proposal feePreResp totalSelected")
		}
		feeSelf := total.Sub(total, fee)

		txOutput, err := p.makeTxOutput(p.request.initiatorAccount.GetContractAccount(), feeSelf.String())
		if err != nil {
			return nil, err
		}
		txOutputs = append(txOutputs, txOutput)
	}

	return txOutputs, nil
}

func (p *Proposal) calcAllFee() (*big.Int, error) {
	allFee := big.NewInt(0)
	if p.request.opt.fee != "" {
		fee, ok := big.NewInt(0).SetString(p.request.opt.fee, 10)
		if !ok {
			return nil, common.ErrInvalidAmount
		}
		allFee.Add(allFee, fee)
	}

	// gas
	gasUsed := big.NewInt(p.preResp.GetResponse().GetGasUsed())
	allFee.Add(allFee, gasUsed)

	return allFee, nil
}

func (p *Proposal) makeTxOutput(addr, amount string) (*pb.TxOutput, error) {
	txOutput := new(pb.TxOutput)
	txOutput.ToAddr = []byte(addr)
	realToAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
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
	if p.cfg.ComplianceCheck.IsNeedComplianceCheck || p.request.opt.onlyFeeFromAccount {
		return initiator
	}

	if p.request.initiatorAccount.HasContractAccount() {
		initiator = p.request.initiatorAccount.GetContractAccount()
	}

	return initiator
}

func (p *Proposal) genInvokeRequests() ([]*pb.InvokeRequest, error) {
	r := p.request
	if r.contractName == "" && r.opt.contractInvokeAmount != "" {
		return nil, errors.New("can not set contract invoke amount")
	}

	if r.module == "" {
		return nil, nil
	}

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

	authRequires := make([]string, 0, 1)
	if p.cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, p.cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	authRequires = append(authRequires, p.request.initiatorAccount.GetAuthRequire())

	if len(p.request.opt.otherAuthRequire) > 0 {
		authRequires = append(authRequires, p.request.opt.otherAuthRequire...)
	}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      p.getChainName(),
		Requests:    invokeRequests,
		Initiator:   p.getInitiator(),
		AuthRequire: authRequires,
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

	if !req.opt.onlyFeeFromAccount && req.opt.fee != "" {
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
	if p.cfg.ComplianceCheck.IsNeedComplianceCheck && p.cfg.ComplianceCheck.IsNeedComplianceCheckFee {
		totalAmount += int64(p.cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
	}

	return totalAmount, nil
}
