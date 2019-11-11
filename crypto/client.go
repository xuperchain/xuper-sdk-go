// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package common is related to generate crypto client
package crypto

import (
	"github.com/xuperchain/xuperunion/common/log"
	crypto_client "github.com/xuperchain/xuperunion/crypto/client"
	"github.com/xuperchain/xuperunion/crypto/client/base"
)

var CryptoTypeConfig = crypto_client.CryptoTypeDefault

// GetCryptoClient get crypto client
func GetCryptoClient() base.CryptoClient {

	// 实例化方法1
	// 通过引用crypto so文件来获取crypto client
	cryptoClient, err := crypto_client.CreateCryptoClient(CryptoTypeConfig)
	if err != nil {
		log.Error("load crypto client failed, %v", err)
	}
	return cryptoClient

	// 方法2
	// 直接引用代码, client目录中package 需由main 修改为 xchain
	//if cryptoClient == nil {
	//	cryptoClient = xchain.GetInstance().(base.CryptoClient)
	//}
	//cryptoClient = xchain.GetInstance().(base.CryptoClient)
	//return cryptoClient
}
