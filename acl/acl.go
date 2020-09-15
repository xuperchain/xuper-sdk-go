package acl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

const (
	newAccount    = "NewAccount"
	setAccountAcl = "SetAccountAcl"
	setMethodAcl  = "SetMethodAcl"
)

type Acl struct {
	xchain.Xchain
	contractAccount string
	contractName    string
	methodName      string
}

// InitACL init a client to acl
func InitAcl(account *account.Account, node, bcname string) *Acl {
	commConfig := config.GetInstance()
	return &Acl{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			XchainSer: node,
			ChainName: bcname,
		},
	}
}

// Pre pre execute request
func (c *Acl) Pre(invokeRequest *pb.InvokeRequest) (*pb.PreExecWithSelectUTXOResponse, error) {
	// generate preExe request
	invokeRequests := []*pb.InvokeRequest{}
	invokeRequests = append(invokeRequests, invokeRequest)

	extraAmount := int64(0)

	authRequires := []string{}

	//xchain.go #597 此处设置了合约账户需要的签名，因此设置ACL时需要赋值xc.ContractAccount
	//func GenRealTxOnly(){
	//	if xc.ContractAccount != "" {
	//		tx.AuthRequireSigns = signatureInfos
	//	}
	//}

	// if the action is set acl
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	// if ComplianceCheck is needed
	// 是否需要进行合规性背书
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)

		// 是否需要支付合规性背书费用
		if c.Cfg.ComplianceCheck.IsNeedComplianceCheckFee == true {
			extraAmount = int64(c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceFee)
		}
	}

	invokeRPCReq := &pb.InvokeRPCRequest{
		Bcname:      c.ChainName,
		Requests:    invokeRequests,
		Initiator:   c.Account.Address,
		AuthRequire: authRequires,
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

// Post post request
func (c *Acl) Post(preExeResp *pb.PreExecWithSelectUTXOResponse) (string, error) {
	authRequires := []string{}

	// if the action is set acl
	if c.ContractAccount != "" {
		authRequires = append(authRequires, c.ContractAccount+"/"+c.Account.Address)
	}

	// if ComplianceCheck is needed
	if c.Cfg.ComplianceCheck.IsNeedComplianceCheck == true {
		authRequires = append(authRequires, c.Cfg.ComplianceCheck.ComplianceCheckEndorseServiceAddr)
	}

	c.Initiator = c.Account.Address
	c.Fee = strconv.Itoa(int(preExeResp.Response.GasUsed))
	//	c.Amount = "0"
	c.TotalToAmount = "0"
	c.AuthRequire = authRequires
	c.InvokeRPCReq = nil
	c.PreSelUTXOReq = nil

	return c.GenCompleteTxAndPost(preExeResp, "")
}

// GenerateInvokeRequest generate invoke request
func (c *Acl) GenerateInvokeRequest(action string, accounts map[string]float32) *pb.InvokeRequest {
	args := make(map[string][]byte)

	switch action {
	case newAccount:
		args["account_name"] = []byte(c.contractAccount)

	case setAccountAcl:
		args["account_name"] = []byte(c.contractAccount)

	case setMethodAcl:
		args["contract_name"] = []byte(c.contractName)
		args["method_name"] = []byte(c.methodName)
	}

	//权限列表
	var aksWeight string
	format := `"%s": %.1f,`
	//遍历地址及权限
	for k, v := range accounts {
		aksWeight += fmt.Sprintf(format, k, v)
	}
	//去掉最后的逗号
	aksWeight = aksWeight[:strings.LastIndex(aksWeight, ",")]
	//构造json
	acl := `
        {
            "pm": {
                "rule": 1,
                "acceptValue": 1.0
            },
            "aksWeight": {
                %s
            }
        }
        `
	acl = fmt.Sprintf(acl, aksWeight)
	args["acl"] = []byte(acl)

	//for k, v := range args {
	//	fmt.Println(k, string(v))
	//}

	return &pb.InvokeRequest{
		ModuleName: "xkernel",
		MethodName: action,
		Args:       args,
	}
}

// InitAclSet init a client to acl, the client can invoke set method
func InitAclSet(account *account.Account, node, bcname, contractAccount, contractName, methodName string) *Acl {
	commConfig := config.GetInstance()
	return &Acl{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			XchainSer: node,
			ChainName: bcname,
		},
		contractAccount: contractAccount,
		contractName:    contractName,
		methodName:      methodName,
	}
}
