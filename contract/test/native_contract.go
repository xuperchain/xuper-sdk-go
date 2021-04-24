// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

package test

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
	"log"
	"strconv"
)

const (
	nativeType     = "native"
	txVersion	=	1
)

type NativeContract struct {
	ContractName string
	xchain.Xchain
}

func InitNativeContract(account *account.Account, bcName, contractName, contractAccount string,sdkClient  *xchain.SDKClient) *NativeContract {
	return &NativeContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             config.GetInstance(),
			Account:         account,
			//XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			SDKClient:		 sdkClient,
		},
	}
}

// Deploy deploy Native contract. args: constructor parameters.
func (c *NativeContract) Deploy(args map[string]string, bin []byte,runtime string) (string, error) {
	ctx := context.Background()
	// preExec
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResp, err := c.PreDeployNativeContract(args, bin, runtime)
		if err != nil {
			return "",err
		}
		//todo
		return c.PostNativeContract(preSelectUTXOResp,"0")
	}else{
		fmt.Println("未check")
		tx,err := c.generateTx(args,bin,runtime)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}



// Deploy deploy Native contract. args: constructor parameters.
func (c *NativeContract) Invoke(args map[string]string, bin []byte,runtime string) (string, error) {
	ctx := context.Background()
	// preExec
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResp, err := c.PreDeployNativeContract(args, bin, runtime)
		if err != nil {
			return "",err
		}
		//todo
		return c.PostNativeContract(preSelectUTXOResp,"0")
	}else{
		fmt.Println("未check")
		tx,err := c.generateTx(args,bin,runtime)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

func (c *NativeContract) generateTx(args map[string]string,code []byte,runtime string)(*pb.Transaction,error){
	ctx := context.Background()

	var invokeRequests []*pb.InvokeRequest
	invokeReq,err := c.generateDeployNativeIR(args,code,runtime)
	if err != nil {
		return nil,err
	}
	invokeRequests = append(invokeRequests, invokeReq)
	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:c.ChainName,
		Requests:invokeRequests,
		Initiator:c.GenInitiator(),
		AuthRequire:authRequires,
	}
	c.InvokeRPCReq = invokeRPCReq
	preInvokeRPCResp,err := c.GenPreExecResp()
	if err != nil {
		return nil,err
	}
	return c.GenRawTx(ctx,preInvokeRPCResp.GetResponse(),preInvokeRPCResp.Response.Requests)
}


//func (c *NativeContract) GenPreExecResp(ctx context.Context,arg map[string]string,code []byte,runtime string)(*pb.InvokeRPCResponse, []*pb.InvokeRequest, error){
//	argsMap := make(map[string]interface{},len(arg))
//	for k,v := range arg{
//		argsMap[k] = v
//	}
//	x3args,err := convertToXuper3Args(argsMap)
//	if err != nil {
//		return nil,nil,err
//	}
//	initArgs,_ := json.Marshal(x3args)
//	desc := &pb.WasmCodeDesc{
//		Runtime:runtime,
//		ContractType:nativeType,
//	}
//	descBuf,err := proto.Marshal(desc)
//	if err != nil {
//		return nil,nil,err
//	}
//	args := map[string][]byte{
//		"account_name":[]byte(c.Account.Address),
//		"contract_name":[]byte(c.ContractName),
//		"contract_code":code,
//		"contract_desc":descBuf,
//		"init_args":initArgs,
//	}
//
//	preExeReqs := []*pb.InvokeRequest{}
//
//	preExecReq := &pb.InvokeRequest{
//		ModuleName:"xkernel",
//		MethodName:"Deploy",
//		Args:       args,
//	}
//	preExeReqs = append(preExeReqs, preExecReq)
//
//	preExeRPCReq := &pb.InvokeRPCRequest{
//		Bcname:bcname,
//		Requests:preExeReqs,
//		Initiator:contractAccount, 		// todo?
//	}
//	c.InvokeRPCReq = preExeRPCReq
//	resp,err := c.PreExec()
//	if err != nil {
//		return nil,nil,err
//	}
//	return resp,preExeReqs,nil
//}






func (c *NativeContract) generateDeployNativeIR(arg map[string]string, bin []byte, runtime string) (*pb.InvokeRequest, error) {
	argsMap := make(map[string]interface{}, len(arg))
	for k, v := range arg {
		argsMap[k] = v
	}

	x3args, err := convertToXuper3Args(argsMap)			// 此处需要修改
	if err != nil {
		return nil, err
	}

	initArgs, _ := json.Marshal(x3args)

	desc := &pb.WasmCodeDesc{
		Runtime:runtime,
		ContractType: nativeType,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(c.ContractAccount),
		"contract_name": []byte(c.ContractName),
		"contract_code": bin,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}

	fmt.Printf("args:")
	fmt.Println(args)

	return &pb.InvokeRequest{				// todo 强制返回xkernel ?
		ModuleName: "xkernel",
		MethodName: "Deploy",
		Args:       args,
	}, nil
}

// PreDeployNativeContract preExecAndSelectUTXO
func (c *NativeContract) PreDeployNativeContract(arg map[string]string, bin []byte,runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	var invokeRequests []*pb.InvokeRequest
	invokeRequest, err := c.generateDeployNativeIR(arg, bin,runtime)
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)

	// 以下代码和部署 wasm 合约时一样，但是没有抽抽来。
	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: extraAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}

// PostNativeContract post and generate complete tx for deploy Native contract.    ???为什么要多这一步？不直接从GenPreExecWithSelectUTXO结束就调用GenompleteTxAndPost
func (c *NativeContract) PostNativeContract(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse, amount string) (string, error) {
	amount, ok := common.IsValidAmount(amount)
	if !ok {
		return "", common.ErrInvalidAmount
	}

	// populates fields
	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	c.Initiator = c.Account.Address
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil
	c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	c.TotalToAmount = "0"

	// EVM 合约调用时可以转账，因此这部分需要增加。     // todo
	if amount != "0" {
		toAddressAndAmount := make(map[string]string)
		toAddressAndAmount[c.ContractName] = amount
		c.ToAddressAndAmount = toAddressAndAmount
		c.TotalToAmount = amount
	}

	return c.GenCompleteTxAndPost(preExeWithSelRes, "")
}

// PreInvokeNativeContract preExe invoker Native contract.
func (c *NativeContract) PreInvokeNativeContract(methodName string, args map[string]string, amount string) (*pb.PreExecWithSelectUTXOResponse, error) {
	amountInt64, err := strconv.ParseInt(amount, 10, 64)
	if err != nil {
		log.Printf("Transfer amount to int64 err: %v", err)
		return nil, err
	}

	var invokeRequests []*pb.InvokeRequest

	invokeRequest, err := c.generateInvokeNativeIR(methodName, args, c.ContractAccount,amount)
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}
	needTotalAmount := amountInt64 + extraAmount

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: needTotalAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}

func (c *NativeContract) generateInvokeNativeIR(methodName string, args map[string]string, contractAccount string, amount string) (*pb.InvokeRequest, error) {
	argsMap := make(map[string]interface{}, len(args))
	for k, v := range args {
		argsMap[k] = v
	}

	irArgs, err := convertToXuper3Args(argsMap)
	if err != nil {
		return nil, err
	}

	ir := &pb.InvokeRequest{
		ModuleName:   nativeType,
		MethodName:   methodName,
		ContractName: c.ContractName,
		Args:         irArgs,
	}

	if amount != "0" {
		ir.Amount = amount
	}

	return ir, nil
}

// Query call Native view function.
func (c *NativeContract) Query(methodName string, args map[string]string) (*pb.InvokeRPCResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest, err := c.generateInvokeNativeIR(methodName, args, c.ContractAccount, "")
	if err != nil {
		return nil, err
	}
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}
	c.InvokeRPCReq = invokeRPCReq

	// preExe
	return c.PreExec()
}


func convertToXuper3Args(args map[string]interface{}) (map[string][]byte, error) {
	argmap := make(map[string][]byte)
	for k, v := range args {
		s, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("bad key %s, expect string value, got %v", k, v)
		}
		argmap[k] = []byte(s)
	}
	return argmap, nil
}

