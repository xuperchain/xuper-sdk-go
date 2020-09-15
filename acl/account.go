package acl

import (
	"math/rand"
	"regexp"
	"strconv"
	"time"

	"github.com/xuperchain/xuper-sdk-go/pb"

	"github.com/xuperchain/xuper-sdk-go/common"
)

// CreateContractAccount create contractAccount
//	contractAccount:
//		case 1: XC + 16 numbers + @bcname
//		case 2: 16 numbers
//		case 3: nothing input
func (c *Acl) CreateContractAccount(contractAccount ...string) (string, error) {
	var account string
	if contractAccount != nil {
		account = contractAccount[0]
	}

	switch alen := len(account); {

	//if contractAccount's format is case 1
	case alen > 16:
		// validate contractAccount
		if ok, _ := regexp.MatchString(`^XC\d{16}@`+c.ChainName+`$`, account); !ok {
			return "", common.ErrInvalidContractAccount
		}
		// get contract account representation that xuper chain used
		subRegexp := regexp.MustCompile(`\d{16}`)
		contractAccountByte := subRegexp.Find([]byte(account))
		account = string(contractAccountByte)

	//if contractAccount's format is case 2
	case alen == 16:
		if ok, _ := regexp.MatchString(`\d{16}`, account); !ok {
			return "", common.ErrInvalidContractAccount
		}

	//if contractAccount's format is case 3 then random generate an account
	default:
		//随机数有几率是0开头,fmt会截断开头的0,使用for来判断生成的账户是否合格
		for len(account) != 16 {
			//account = fmt.Sprintf("%08v",
			//	rand.New(rand.NewSource(time.Now().UnixNano())).Int63n(10000000000000000))
			account = (strconv.Itoa(int(time.Now().Unix())) + strconv.Itoa(rand.Int()))[0:16]
		}
	}

	c.contractAccount = account

	//the address weight is 1
	accounts := map[string]float32{
		c.Account.Address: 1,
	}
	//generate invoke request
	invokeRequest := c.GenerateInvokeRequest(newAccount, accounts)

	// pre
	preExeResp, err := c.Pre(invokeRequest)
	if err != nil {
		return "", err
	}
	// post
	_, err = c.Post(preExeResp)
	if err != nil {
		return "", err
	}

	//return the contractAccount
	return "XC" + account + "@" + c.ChainName, nil
}

// SetContractAccountAcl set contractAccount's acl by address list to 1
func (c *Acl) SetContractAccountAcl(contractAccount string, address []string) (string, error) {
	// validate contractAccount
	if ok, _ := regexp.MatchString(`^XC\d{16}@`+c.ChainName+`$`, contractAccount); !ok {
		return "", common.ErrInvalidContractAccount
	}
	c.contractAccount = contractAccount

	//the all address weight is 1
	accounts := make(map[string]float32)
	for _, v := range address {
		accounts[v] = 1
	}

	return c.SetAccountAcl(accounts)
}

// SetAccountAcl set contractAccount's acl by account's acl list
// the method can customize the weight of ever address
func (c *Acl) SetAccountAcl(accounts map[string]float32) (string, error) {
	//set contractAccount's AuthRequireSigns
	c.ContractAccount = c.contractAccount

	//generate invoke request
	invokeRequest := c.GenerateInvokeRequest(setAccountAcl, accounts)

	// pre
	preExeResp, err := c.Pre(invokeRequest)
	if err != nil {
		return "", err
	}
	// post
	return c.Post(preExeResp)
}

// QueryAccountAcl query contractAccount's acl
func (c *Acl) QueryAccountAcl(contractAccount string) (*pb.AclStatus, error) {
	return c.Xchain.QueryAccountAcl(contractAccount)
}
