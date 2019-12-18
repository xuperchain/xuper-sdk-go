// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to generate crypto client
package crypto

import (
	"github.com/xuperchain/xuperunion/crypto/client/base"
	"github.com/xuperchain/xuperunion/crypto/client/xchain"
)

func getInstance() interface{} {
	return &eccdefault.XchainCryptoClient{}
}

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {
	cryptoClient := getInstance().(base.CryptoClient)
	return cryptoClient
}
