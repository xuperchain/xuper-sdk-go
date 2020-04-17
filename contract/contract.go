// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contract is related to contract operation
package contract

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuperchain/core/pb"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/xchain"

	"github.com/xuperdata/teesdk"
)

// WasmContract wasmContract structure
type WasmContract struct {
	ContractName string
	tfc          *teesdk.TEEClient
	xchain.Xchain
}

// InitWasmContract init a client to deploy/invoke/query a wasm contract
func InitWasmContract(account *account.Account, node, bcName, contractName, contractAccount string) *WasmContract {
	commConfig := config.GetInstance()

	wc := &WasmContract{
		ContractName: contractName,
		Xchain: xchain.Xchain{
			Cfg:             commConfig,
			Account:         account,
			XchainSer:       node,
			ChainName:       bcName,
			ContractAccount: contractAccount,
		},
	}

	if commConfig.TC.Enable {
		wc.tfc = teesdk.NewTEEClient(
			commConfig.TC.Uid,
			commConfig.TC.Token,
			commConfig.TC.Auditors[0].PublicDer,
			commConfig.TC.Auditors[0].Sign,
			commConfig.TC.Auditors[0].EnclaveInfoConfig,
			commConfig.TC.TMSPort)
	}
	return wc
}

// DeployWasmContract deploy a wasm contract
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

// PreDeployWasmContract preExe deploy wasm contract
func (c *WasmContract) PreDeployWasmContract(arg map[string]string, codepath string, runtime string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	var invokeRequests []*pb.InvokeRequest
	invokeRequest := generateDeployIR(arg, codepath, runtime, c.ContractAccount, c.ContractName)
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}
	authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}
	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
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
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	c.Initiator = c.Account.Address
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil
	c.Fee = strconv.Itoa(int(preExeWithSelRes.Response.GasUsed))
	c.Amount = "0"

	return c.GenCompleteTxAndPost(preExeWithSelRes)
}

// InvokeWasmContract invoke wasm contract by method name
func (c *WasmContract) InvokeWasmContract(methodName string, args map[string]string) (string, error) {
	// preExe
	commConfig := config.GetInstance()
	// TODO fix bug
	if commConfig.TC.Enable && methodName == "store" {
		encryptedArgs, err := c.EncryptArgs(commConfig.TC.Svn, args)
		if err != nil {
			log.Println("EncryptArgs error,", err)
			return "", err
		}
		args = map[string]string{}
		err = json.Unmarshal([]byte(encryptedArgs), &args)
		if err != nil {
			return "", err
		}
	}
	preSelectUTXOResponse, err := c.PreInvokeWasmContract(methodName, args)
	if err != nil {
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
	authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
	}
	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:      c.ChainName,
		Address:     c.Account.Address,
		TotalAmount: int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
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
