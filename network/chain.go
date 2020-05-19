// Copyright (c) 2020. Baidu Inc. All Rights Reserved.

// package chain is related to create new blockchain
package network

import (
	"encoding/json"
	"errors"
	"log"
	"math/big"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/transfer"
	"github.com/xuperchain/xuper-sdk-go/txhash"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

var (
	ErrNegativeAmount = errors.New("amount in transaction can not be negative number")
)

const (
	RootTxVersion = 0
)

// Chain structure
type Chain struct {
	//	ContractName string
	transfer.Trans
}

// InitChain init a client to operate with chain
func InitChain(account *account.Account, node, bcname string) *Chain {
	commConfig := config.GetInstance()

	return &Chain{
		Trans: transfer.Trans{
			Xchain: xchain.Xchain{
				Cfg:       commConfig,
				Account:   account,
				XchainSer: node,
				ChainName: bcname,
			},
		},
	}
}

type DescJSON struct {
	Module string `json:"Module"`
	Method string `json:"Method"`
	Args   struct {
		Name string `json:"name"`
		Data string `json:"data"`
	} `json:"Args"`
}

type ParamJSON struct {
	Version   string `json:"version"`
	Consensus struct {
		Miner string `json:"miner"`
	} `json:"consensus"`
	Predistribution []struct {
		Address string `json:"address"`
		Quota   string `json:"quota"`
	} `json:"predistribution"`
}

// Create a brand new blockchain
func (c *Chain) CreateChain(desc string) (string, error) {
	// step 1: get the default config json content
	//	filePath := "./conf/xuper.json"
	//	configJson, err := ioutil.ReadFile(filePath)
	//	if err != nil {
	//		fmt.Println("read file " + filePath + " error")
	//		return err
	//	}

	descJson := []byte(desc)
	descJsonObj := &DescJSON{}
	jsErr := json.Unmarshal(descJson, descJsonObj)
	if jsErr != nil {
		return "", jsErr
	}

	to := descJsonObj.Args.Name
	amount := c.Trans.Xchain.Cfg.MinNewChainAmount
	fee := "0"
	txid, err := c.Trans.Transfer(to, amount, fee, desc)
	if err != nil {
		log.Printf("create chain err: %v\n", err)
		return "", err
	}

	log.Printf("Real txid: %v\n", txid)
	return txid, nil
}

// generateCreateChainTx - genereate a tx which will be used for creating a new blockchain
func (c *Chain) generateCreateChainTx(desc string) (*pb.Transaction, error) {
	descJson := []byte(desc)
	descJsonObj := &DescJSON{}
	jsErr := json.Unmarshal(descJson, descJsonObj)
	if jsErr != nil {
		return nil, jsErr
	}

	paramJson := []byte(descJsonObj.Args.Data)

	jsObj := &ParamJSON{}
	jsErr = json.Unmarshal(paramJson, jsObj)
	if jsErr != nil {
		return nil, jsErr
	}

	utxoTx := &pb.Transaction{Version: RootTxVersion}

	// 暂时只支持配置给一个初始化账户打钱
	for _, pd := range jsObj.Predistribution {
		amount := big.NewInt(0)
		amount.SetString(pd.Quota, 10) // 10进制转换大整数
		if amount.Cmp(big.NewInt(0)) < 0 {
			return nil, ErrNegativeAmount
		}
		txOutput := &pb.TxOutput{}
		txOutput.ToAddr = []byte(pd.Address)
		txOutput.Amount = amount.Bytes()
		utxoTx.TxOutputs = append(utxoTx.TxOutputs, txOutput)

		// 暂时只支持配置给一个初始化账户打钱
		break
	}

	utxoTx.Desc = []byte(desc)
	utxoTx.Coinbase = true
	utxoTx.Txid, _ = txhash.MakeTransactionID(utxoTx)

	return utxoTx, nil
}
