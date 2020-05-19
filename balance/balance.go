// Copyright (c) 2020. Baidu Inc. All Rights Reserved.

// package chain is related to create new blockchain
package balance

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"time"

	"google.golang.org/grpc"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
)

// Chain structure
type Balance struct {
	Cfg       *config.CommConfig
	Account   *account.Account
	XchainSer string
	BcNames   []string
}

// InitBalance init a client to query with chain
func InitBalance(account *account.Account, node string, bcNames []string) *Balance {
	commConfig := config.GetInstance()

	return &Balance{
		Cfg:       commConfig,
		Account:   account,
		XchainSer: node,
		BcNames:   bcNames,
	}
}

// GetBalanceDetail get unfrozen balance and frozen balance
func (bal *Balance) GetBalanceDetails() (string, error) {
	conn, err := grpc.Dial(bal.XchainSer, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
	if err != nil {
		log.Printf("GetBalance connect xchain err: %v", err)
		return "", err
	}
	defer conn.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 15000*time.Millisecond)
	defer cancel()

	tfds := []*pb.TokenFrozenDetails{}
	for _, bcName := range bal.BcNames {
		tfds = append(tfds, &pb.TokenFrozenDetails{Bcname: bcName})
	}

	addrStatus := &pb.AddressBalanceStatus{
		Address: bal.Account.Address,
		Tfds:    tfds,
	}

	c := pb.NewXchainClient(conn)
	res, err := c.GetBalanceDetail(ctx, addrStatus)
	if err != nil {
		return "", err
	}
	if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return "", errors.New(res.Header.Error.String())
	}

	balanceJSON, err := json.Marshal(res.Tfds)
	return string(balanceJSON), err
}
