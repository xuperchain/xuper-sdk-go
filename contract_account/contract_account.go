// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package contractaccount is related to contract account operation
package contractaccount

import (
	"regexp"
	"strconv"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

// ContractAccount contract account
type ContractAccount struct {
	xchain.Xchain
}

// InitContractAccount init a client to create contractAccount
func InitContractAccount(account *account.Account, node, bcname string) *ContractAccount {
	commConfig := config.GetInstance()
	return &ContractAccount{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			XchainSer: node,
			ChainName: bcname,
		},
	}
}

// CreateContractAccount create contractAccount
func (ca *ContractAccount) CreateContractAccount(contractAccount string) (string, error) {
	// preExe
	preExeResp, err := ca.PreCreateContractAccount(contractAccount)
	if err != nil {
		return "", err
	}
	// post
	return ca.PostCreateContractAccount(preExeResp)
}

// PreCreateContractAccount preExe create contract account
func (ca *ContractAccount) PreCreateContractAccount(contractAccount string) (*pb.PreExecWithSelectUTXOResponse, error) {
	// validate contractAccount
	if ok, _ := regexp.MatchString(`^XC\d{16}@`+ca.ChainName+`$`, contractAccount); !ok {
		return nil, common.ErrInvalidContractAccount
	}

	// get contract account representation that xuper chain used
	subRegexp := regexp.MustCompile(`\d{16}`)
	contractAccountByte := subRegexp.Find([]byte(contractAccount))
	contractAccount = string(contractAccountByte)

	// generate preExe request
	invokeRequests := []*pb.InvokeRequest{}
	invokeRequest := generateInvokeRequest(contractAccount, ca.Account.Address)
	invokeRequests = append(invokeRequests, invokeRequest)

	authRequires := []string{}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:    ca.ChainName,
		Requests:  invokeRequests,
		Initiator: ca.Account.Address,
		//		AuthRequire: authRequires,
	}

	extraAmount := int64(0)

	// if ComplianceCheck is needed
	// 是否需要进行合规性背书
	if ca.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, ca.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
		invokeRPCReq.AuthRequire = authRequires

		// 是否需要支付合规性背书费用
		if ca.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(ca.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}

	}

	preSelUTXOReq := &pb.PreExecWithSelectUTXORequest{
		Bcname:  ca.ChainName,
		Address: ca.Account.Address,
		//		TotalAmount: int64(ca.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee),
		TotalAmount: extraAmount,
		Request:     invokeRPCReq,
	}
	ca.InvokeRPCReq = invokeRPCReq
	ca.PreSelUTXOReq = preSelUTXOReq

	// preExe
	return ca.PreExecWithSelecUTXO()
}

// PostCreateContractAccount generate complete Tx and post to create contract account
func (ca *ContractAccount) PostCreateContractAccount(preExeResp *pb.PreExecWithSelectUTXOResponse) (string, error) {
	authRequires := []string{}
	// if ComplianceCheck is needed
	if ca.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, ca.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	ca.Initiator = ca.Account.Address
	ca.Fee = strconv.Itoa(int(preExeResp.Response.GasUsed))
	//	ca.Amount = "0"
	ca.TotalToAmount = "0"
	ca.AuthRequire = authRequires
	ca.InvokeRPCReq = nil
	ca.PreSelUTXOReq = nil

	return ca.GenCompleteTxAndPost(preExeResp, "")
}

func generateInvokeRequest(contractAccount, address string) *pb.InvokeRequest {
	args := make(map[string][]byte)
	args["account_name"] = []byte(contractAccount)

	defaultACL := `
        {
            "pm": {
                "rule": 1,
                "acceptValue": 1.0
            },
            "aksWeight": {
                "` + address + `": 1.0
            }
        }
        `
	args["acl"] = []byte(defaultACL)

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: "NewAccount",
		Args:       args,
	}
}
