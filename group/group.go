package group

import (
	"strings"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/contract"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

const (
	contractName = "group_chain"
	addChain     = "addChain"
	delChain     = "delChain"
	listChain    = "listChain"
	addNode      = "addNode"
	delNode      = "delNode"
	listNode     = "listNode"
)

type Group struct {
	xchain.Xchain
	wasm *contract.WasmContract
}

// InitGroup init a client to group
func InitGroup(account *account.Account, node, bcName string) *Group {
	commConfig := config.GetInstance()
	return &Group{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			XchainSer: node,
			ChainName: bcName,
		},
		wasm: contract.InitWasmContract(account, node, bcName, contractName, ""),
	}
}

func (c *Group) AddChain(bcname string) (string, error) {
	return c.wasm.InvokeWasmContract(addChain, map[string]string{"bcname": bcname})
}

func (c *Group) DelChain(bcname string) (string, error) {
	return c.wasm.InvokeWasmContract(delChain, map[string]string{"bcname": bcname})
}

func (c *Group) ListChain() ([]string, error) {
	resp, err := c.wasm.QueryWasmContract(listChain, nil)
	if err != nil {
		return nil, err
	}
	var result string
	for _, data := range resp.GetResponse().GetResponse() {
		result += strings.ReplaceAll(string(data), "\u0001", ",")
	}
	result = strings.Trim(result, ",")
	return strings.Split(result, ","), nil
}

func (c *Group) AddNode(bcname, neturl, address string) (string, error) {
	return c.wasm.InvokeWasmContract(addNode, map[string]string{
		"bcname":  bcname,
		"ip":      neturl,
		"address": address,
	})
}

func (c *Group) DelNode(bcname, neturl string) (string, error) {
	return c.wasm.InvokeWasmContract(delNode, map[string]string{
		"bcname": bcname,
		"ip":     neturl,
	})
}

func (c *Group) ListNode(bcname string) ([]string, error) {
	resp, err := c.wasm.QueryWasmContract(listNode, map[string]string{
		"bcname": bcname,
	})
	if err != nil {
		return nil, err
	}
	var result string
	for _, data := range resp.GetResponse().GetResponse() {
		result += strings.ReplaceAll(string(data), "\u0001", ",")
	}
	result = strings.Trim(result, ",")
	return strings.Split(result, ","), nil
}
