// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package crypto is related to generate crypto client
package crypto

import (
	"github.com/xuperchain/crypto/client/service/base"
	"github.com/xuperchain/crypto/client/service/gm"
	"github.com/xuperchain/crypto/client/service/xchain"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
)

func getInstance() interface{} {
	switch config.GetInstance().Crypto {
	case config.CRYPTO_XCHAIN:
		return &xchain.XchainCryptoClient{}
	case config.CRYPTO_GM:
		return &gm.GmCryptoClient{}
	default:
		return &xchain.XchainCryptoClient{}
	}
}

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {
	cryptoClient := getInstance().(base.CryptoClient)
	return cryptoClient
}

// GetXchainCryptoClient get xchain crypto client
func GetXchainCryptoClient() *xchain.XchainCryptoClient {
	return &xchain.XchainCryptoClient{}
}

// GetGmCryptoClient get gm crypto client
func GetGmCryptoClient() *gm.GmCryptoClient {
	return &gm.GmCryptoClient{}
}
