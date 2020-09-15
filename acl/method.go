package acl

import (
	"regexp"

	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/pb"
)

// SetContractMethodAcl set contractMethod's acl by address list to 1
func (c *Acl) SetContractMethodAcl(contractAccount, contractName, methodName string, address []string) (string, error) {
	// validate contractAccount
	if ok, _ := regexp.MatchString(`^XC\d{16}@`+c.ChainName+`$`, contractAccount); !ok {
		return "", common.ErrInvalidContractAccount
	}
	c.contractAccount = contractAccount
	c.contractName = contractName
	c.methodName = methodName

	//the all address weight is 1
	accounts := make(map[string]float32)
	for _, v := range address {
		accounts[v] = 1
	}

	return c.SetMethodAcl(accounts)
}

// SetMethodAcl set contractMethod's acl By account's acl list
// the method can customize the weight of ever address
func (c *Acl) SetMethodAcl(accounts map[string]float32) (string, error) {
	//set contractAccount's AuthRequireSigns
	c.ContractAccount = c.contractAccount

	//generate invoke request
	invokeRequest := c.GenerateInvokeRequest(setMethodAcl, accounts)

	// pre
	preExeResp, err := c.Pre(invokeRequest)
	if err != nil {
		return "", err
	}
	// post
	return c.Post(preExeResp)
}

// QueryMethodAcl query contractMethod's acl
func (c *Acl) QueryMethodAcl(contractName, methodName string) (*pb.AclStatus, error) {
	return c.Xchain.QueryMethodAcl(contractName, methodName)
}
