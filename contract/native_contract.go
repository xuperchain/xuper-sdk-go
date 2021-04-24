// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"context"
	"log"
	"strconv"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

const nativeType     = "native"
// WasmContract wasmContract structure
type NativeContract struct {
	ContractName string
	xchain.Xchain
}

// InitWasmContract init a client to deploy/invoke/query a wasm contract
func InitNativeContract(account *account.Account, bcName, contractName, contractAccount string,sdkClient *xchain.SDKClient) *NativeContract {
	commConfig := config.GetInstance()

	return &NativeContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			SDKClient:sdkClient,
		},
	}
}

// DeployWasmContract deploy a wasm contract
func (c *NativeContract) DeployNativeContract(args map[string]string, codepath string, runtime string) (string, error) {
	ctx := context.Background()
	// preExec
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResp, err := c.PreDeployNativeContract(args, codepath, runtime)
		if err != nil {
			return "",err
		}
		//todo
		return c.PostNativeContract(preSelectUTXOResp)
	}else{
		deployReq := GenerateDeployInvokeReq(args,codepath,runtime,c.ContractAccount,c.ContractName,nativeType)

		tx,err := c.GenerateTx(deployReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

func (c *NativeContract) UpgradeNativeContract(args map[string]string, codepath string) (string, error) {
	ctx := context.Background()
	// preExec
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResp, err := c.PreUpgradeNativeContract(args, codepath)
		if err != nil {
			return "",err
		}
		//todo
		return c.PostNativeContract(preSelectUTXOResp)
	}else{
		upgradeReq := generateUpgradeInvokReq(args,codepath,c.ContractAccount,c.ContractName,nativeType)

		tx,err := c.GenerateTx(upgradeReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

// PreDeployWasmContract preExe deploy wasm contract
func (c *NativeContract) PreDeployNativeContract(arg map[string]string, codepath string, runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := GenerateDeployInvokeReq(arg, codepath, runtime, c.ContractAccount, c.ContractName,nativeType)
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:  c.ChainName,
		Address: c.Account.Address,
		//		TotalAmount: int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
		TotalAmount: extraAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}


func (c *NativeContract) PreUpgradeNativeContract(arg map[string]string, codepath string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateUpgradeInvokReq(arg, codepath, c.ContractAccount, c.ContractName,nativeType)
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:  c.ChainName,
		Address: c.Account.Address,
		//		TotalAmount: int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
		TotalAmount: extraAmount,
		Request:     invokeRPCReq,
	}
	c.InvokeRPCReq = invokeRPCReq
	c.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return c.PreExecWithSelecUTXO()
}




// PostWasmContract generate complete Tx and post to deploy wasm Contract
func (c *NativeContract) PostNativeContract(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse) (string, error) {
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
	//	c.Amount = "0"
	c.TotalToAmount = "0"

	return c.GenCompleteTxAndPost(preExeWithSelRes, "")
}

// InvokeWasmContract invoke wasm contract by method name
func (c *NativeContract) InvokeNativeContract(methodName string, args map[string]string,fee string) (string, error) {
	ctx := context.Background()
	c.Fee = fee
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResponse, err := c.PreInvokeNativeContract(methodName, args)
		if err != nil {
			log.Printf("InvokeWasmContract preExe failed, err: %v", err)
			return "", err
		}
		// post
		return c.PostNativeContract(preSelectUTXOResponse)
	}else{
		invokeReq := generateInvokeInvokeReq(args,methodName,c.ContractName,nativeType)
		tx,err := c.GenerateTx(invokeReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

// PreInvokeWasmContract preExe invoke wasm contract
func (c *NativeContract) PreInvokeNativeContract(methodName string, args map[string]string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeInvokeReq(args, methodName, c.ContractName,nativeType)
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
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
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

// QueryWasmContract query wasm contract, same as preExe invoke
func (c *NativeContract) QueryNativeContract(methodName string, args map[string]string) (*pb.InvokeRPCResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeInvokeReq(args, methodName, c.ContractName,nativeType)
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

