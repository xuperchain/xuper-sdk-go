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

const nativeType = "native"

// WasmContract wasmContract structure
type NativeContract struct {
	ContractName string
	xchain.Xchain
}

// InitWasmContractWithClient init a client to deploy/invoke/query a wasm contract
func InitNativeContractWithClient(account *account.Account, bcName, contractName, contractAccount string, xuperClient *xchain.XuperClient) *NativeContract {
	commConfig := config.GetInstance()
	return &NativeContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:     commConfig,
			Account: account,
			//XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			XuperClient:     xuperClient,
		},
	}
}

func InitNativeContract(account *account.Account, node, bcName, contractName, contractAccount string) *NativeContract {
	commConfig := config.GetInstance()
	client, err := xchain.NewXuperClient(node)
	if err != nil {
		return nil
	}
	return &NativeContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
			XuperClient:     client,
		},
	}
}

// DeployWasmContract deploy a wasm contract
func (c *NativeContract) DeployNativeContract(args map[string]string, codepath string, runtime string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreDeployNativeContract(args, codepath, runtime)
	if err != nil {
		log.Printf("DeployWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostNativeContract(preSelectUTXOResponse)
}

// PreDeployWasmContract preExe deploy wasm contract
func (c *NativeContract) PreDeployNativeContract(arg map[string]string, codepath string, runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateDeployNativeIR(arg, codepath, runtime, c.ContractAccount, c.ContractName)
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
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
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
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
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
func (c *NativeContract) InvokeNativeContract(methodName string, args map[string]string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreInvokeNativeContract(methodName, args)
	if err != nil {
		log.Printf("InvokeWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostNativeContract(preSelectUTXOResponse)
}

func (c *NativeContract) UpgradeNativeContract(args map[string]string, codepath string) (string, error) {
	// preExec
	preSelectUTXOResp, err := c.PreUpgradeNativeContract(args, codepath)
	if err != nil {
		return "", err
	}
	return c.PostNativeContract(preSelectUTXOResp)
}

func (c *NativeContract) PreUpgradeNativeContract(arg map[string]string, codepath string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateUpgradeInvokReq(arg, codepath, c.ContractAccount, c.ContractName, nativeType)
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
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee {
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

func generateUpgradeInvokReq(arg map[string]string, codepath string, contractAccount, contractName, contractType string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		ContractType: contractType,
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
		MethodName: "Upgrade",
		Args:       args,
	}
}

// PreInvokeWasmContract preExe invoke wasm contract
func (c *NativeContract) PreInvokeNativeContract(methodName string, args map[string]string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeNativeIR(args, methodName, c.ContractName)
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

		//		extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
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

// QueryWasmContract query wasm contract, same as preExe invoke
func (c *NativeContract) QueryNativeContract(methodName string, args map[string]string) (*pb.InvokeRPCResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateInvokeNativeIR(args, methodName, c.ContractName)
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

func generateDeployNativeIR(arg map[string]string, codepath string, runtime, contractAccount, contractName string) *pb.InvokeRequest {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	contractCode, err := ioutil.ReadFile(codepath)
	if err != nil {
		log.Printf("get wasm contract code error: %v", err)
		return nil
	}
	desc := &pb.WasmCodeDesc{
		ContractType: nativeType,
		Runtime:      runtime,
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
		MethodName: "Deploy",
		Args:       args,
	}
}

func generateInvokeNativeIR(args map[string]string, methodName, contractName string) *pb.InvokeRequest {
	return &pb.InvokeRequest{
		ModuleName:   "native",
		MethodName:   methodName,
		ContractName: contractName,
		Args:         convertToXuperContractArgs(args),
	}
}

//func convertToXuperContractArgs(args map[string]string) map[string][]byte {
//	argmap := make(map[string][]byte)
//	for k, v := range args {
//		argmap[k] = []byte(v)
//	}
//	return argmap
//}
