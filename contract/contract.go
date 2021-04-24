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

const wasmType  = "wasm"


// WasmContract wasmContract structure
type WasmContract struct {
	ContractName string
	xchain.Xchain
}

// InitWasmContract init a client to deploy/invoke/query a wasm contract
func InitWasmContract(account *account.Account, node, bcName, contractName, contractAccount string) *WasmContract {
	commConfig := config.GetInstance()

	return &WasmContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
		},
	}
}

// DeployWasmContract deploy a wasm contract
func (c *WasmContract) DeployWasmContract(args map[string]string, codepath string, runtime string) (string, error) {
	// preExe
	ctx := context.Background()
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResponse, err := c.PreDeployWasmContract(args, codepath, runtime)
		if err != nil {
			log.Printf("DeployWasmContract preExe failed, err: %v", err)
			return "", err
		}
		// post
		return c.PostWasmContract(preSelectUTXOResponse)
	}else{
		deployReq := GenerateDeployInvokeReq(args,codepath,runtime,c.ContractAccount,c.ContractName,wasmType)
		tx,err := c.GenerateTx(deployReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

func (c *WasmContract) UpgradeWasmContract(args map[string]string, codepath string) (string, error) {
	// preExe
	ctx := context.Background()

	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResponse, err := c.PreUpgradeWasmContract(args, codepath)
		if err != nil {
			log.Printf("DeployWasmContract preExe failed, err: %v", err)
			return "", err
		}
		// post
		return c.PostWasmContract(preSelectUTXOResponse)
	}else{
		upgradeReq := generateUpgradeInvokReq(args,codepath,c.ContractAccount,c.ContractName,wasmType)

		tx,err := c.GenerateTx(upgradeReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}
}

// PreDeployWasmContract preExe deploy wasm contract
func (c *WasmContract) PreDeployWasmContract(arg map[string]string, codepath string, runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := GenerateDeployInvokeReq(arg, codepath, runtime, c.ContractAccount, c.ContractName,wasmType)
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

func (c *WasmContract) PreUpgradeWasmContract(arg map[string]string, codepath string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateUpgradeInvokReq(arg, codepath, c.ContractAccount, c.ContractName,wasmType)
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
func (c *WasmContract) PostWasmContract(preExeWithSelRes *pb.PreExecWithSelectUTXOResponse) (string, error) {
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
func (c *WasmContract) InvokeWasmContract(methodName string, args map[string]string) (string, error) {
	// preExe
	ctx := context.Background()
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		preSelectUTXOResponse, err := c.PreInvokeWasmContract(methodName, args)
		if err != nil {
			log.Printf("InvokeWasmContract preExe failed, err: %v", err)
			return "", err
		}
		// post
		return c.PostWasmContract(preSelectUTXOResponse)
	}else{
		invokeReq := generateInvokeInvokeReq(args,methodName,c.ContractName,wasmType)
		tx,err := c.GenerateTx(invokeReq)
		if err != nil {
			return "",err
		}
		return c.SendTx(ctx,tx)
	}

}

// PreInvokeWasmContract preExe invoke wasm contract
func (c *WasmContract) PreInvokeWasmContract(methodName string, args map[string]string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeInvokeReq(args, methodName, c.ContractName,wasmType)
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
func (c *WasmContract) QueryWasmContract(methodName string, args map[string]string) (*pb.InvokeRPCResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeInvokeReq(args, methodName, c.ContractName,wasmType)
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

	return c.PreExec()
}

func convertToXuperContractArgs(args map[string]string) map[string][]byte {
	argmap := make(map[string][]byte)
	for k, v := range args {
		argmap[k] = []byte(v)
	}
	return argmap
}
