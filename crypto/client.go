// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to generate crypto client
package crypto

import (
	//	"github.com/xuperchain/crypto/client/service/gm"
	"github.com/xuperchain/crypto/client/service/xchain"
	//	"github.com/xuperchain/xuper-sdk-go/config"
	//	"github.com/xuperchain/xuperchain/core/crypto/client/base"
)

//func getInstance() interface{} {
//	// @todo 更新xchain crypto的依赖到crypto库
//	switch config.GetInstance().Crypto {
//	case config.CRYPTO_XCHAIN:
//		return &eccdefault.XchainCryptoClient{}
//	case config.CRYPTO_GM:
//		return &gm.GmCryptoClient{}
//	default:
//		return &eccdefault.XchainCryptoClient{}
//	}
//}
//
//// GetCryptoClient get crypto client
//func GetCryptoClient() base.CryptoClient {
//	cryptoClient := getInstance().(base.CryptoClient)
//	return cryptoClient
//}

// GetCryptoClient get crypto client
func GetCryptoClient() *xchain.XchainCryptoClient {
	return &xchain.XchainCryptoClient{}
}

// GetCryptoClient get crypto client
func GetXchainCryptoClient() *xchain.XchainCryptoClient {
	return &xchain.XchainCryptoClient{}
}
