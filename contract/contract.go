// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

const (
	deploy  = "Deploy"
	upgrade = "Upgrade"
)

// WasmContract wasmContract structure
type WasmContract struct {
	ContractName string
	xchain.Xchain
}

// InitWasmContract init a client to deploy/upgrade/invoke/query a wasm contract
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
func (c *WasmContract) DeployWasmContract(args map[string]string, codepath, runtime string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreDeployOrUpgradeWasmContract(args, codepath, runtime, deploy)
	if err != nil {
		log.Printf("DeployWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostWasmContract(preSelectUTXOResponse)
}

// UpgradeWasmContract upgrade a wasm contract
func (c *WasmContract) UpgradeWasmContract(args map[string]string, codepath, runtime string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreDeployOrUpgradeWasmContract(args, codepath, runtime, upgrade)
	if err != nil {
		log.Printf("UpgradeWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostWasmContract(preSelectUTXOResponse)
}

// PreDeployOrUpgradeWasmContract preExe deploy or upgrade wasm contract
func (c *WasmContract) PreDeployOrUpgradeWasmContract(arg map[string]string, codepath, runtime, action string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateDeployOrUpgradeIR(arg, codepath, runtime, c.ContractAccount, c.ContractName, action)
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
	if preExeWithSelRes.Response != nil {
		c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	}
	//	c.Amount = "0"
	c.TotalToAmount = "0"

	return c.GenCompleteTxAndPost(preExeWithSelRes, "")
}

// InvokeWasmContract invoke wasm contract by method name
func (c *WasmContract) InvokeWasmContract(methodName string, args map[string]string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreInvokeWasmContract(methodName, args)
	if err != nil {
		log.Printf("InvokeWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostWasmContract(preSelectUTXOResponse)
}

// PreInvokeWasmContract preExe invoke wasm contract
func (c *WasmContract) PreInvokeWasmContract(methodName string, args map[string]string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeIR(args, methodName, c.ContractName)
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
	invokeRequest := generateInvokeIR(args, methodName, c.ContractName)
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

func generateDeployOrUpgradeIR(arg map[string]string, codepath, runtime, contractAccount, contractName, action string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		Runtime: runtime,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(contractAccount),
		"contract_name": []byte(contractName),
		"contract_code": contractCode,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: action,
		Args:       args,
	}
}

func generateInvokeIR(args map[string]string, methodName, contractName string) *pb.InvokeRequest {
	return &pb.InvokeRequest{
		ModuleName:   "wasm",
		MethodName:   methodName,
		ContractName: contractName,
		Args:         convertToXuperContractArgs(args),
	}
}

func convertToXuperContractArgs(args map[string]string) map[string][]byte {
	argmap := make(map[string][]byte)
	for k, v := range args {
		argmap[k] = []byte(v)
	}
	return argmap
}
