// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to generate crypto client
package crypto

import (
	"github.com/xuperchain/xuperchain/core/crypto/client/base"
	gmclient "github.com/xuperchain/xuperchain/core/crypto/client/gm/gmclient"
	xchain "github.com/xuperchain/xuperchain/core/crypto/client/xchain"
	"github.com/xuperchain/xuper-sdk-go/config"
)

func getInstance() interface{} {
	if config.GetInstance().Crypto == "gm" {
		return &gmclient.GmCryptoClient{}
	} else {
		return &xchain.XchainCryptoClient{}
	}
}

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {
	cryptoClient := getInstance().(base.CryptoClient)
	return cryptoClient
}
