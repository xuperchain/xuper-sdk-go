// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package contract 智能合约相关操作，包括 wasm 以及 evm 合约
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

// WasmContract wasm 合约
type WasmContract struct {
	ContractName string
	xchain.Xchain
}

// InitWasmContract 创建 wasm 合约实例，用来调用合约。
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

// DeployWasmContract 部署 wasm 合约。没有错误返回交易 ID。
func (c *WasmContract) DeployWasmContract(args map[string]string, codepath string, runtime string) (string, error) {
	// preExe
	preSelectUTXOResponse, err := c.PreDeployWasmContract(args, codepath, runtime)
	if err != nil {
		log.Printf("DeployWasmContract preExe failed, err: %v", err)
		return "", err
	}
	// post
	return c.PostWasmContract(preSelectUTXOResponse)
}

// PreDeployWasmContract 部署合约的预执行接口，返回预执行结果。
func (c *WasmContract) PreDeployWasmContract(arg map[string]string, codepath string, runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateDeployIR(arg, codepath, runtime, c.ContractAccount, c.ContractName)
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

// PostWasmContract 合约部署与调用交易的 post 接口，没有错误返回交易 ID。
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

// InvokeWasmContract 根据方法名字与参数调用合约。
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

// PreInvokeWasmContract 调用合约交易的预执行接口。
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

// QueryWasmContract 查询 wasm 接口，不消耗手续费，类似于调用合约的预执行阶段。
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

func generateDeployIR(arg map[string]string, codepath string, runtime, contractAccount, contractName string) *pb.InvokeRequest {
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
		MethodName: "Deploy",
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
