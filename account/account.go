// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package account 账户相关操作，包括创建账户、获取账户、xchain 与 evm 地址转换等。
package account

import (
	"io/ioutil"
	"log"

	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/crypto"
)

// Account 账户结构
type Account struct {
	// Address 账户地址
	Address string

	// PrivateKey 账户私钥
	PrivateKey string

	// PublicKey 账户公钥
	PublicKey string

	// Mnemonic 账户的助记词
	Mnemonic string
}

// CreateAccount 根据助记词强度以及助记词语言创建账户
//
//Parameters:
//   - `strength`：1弱（12个助记词），2中（18个助记词），3强（24个助记词）。
//   - `language`：1中文，2英文。
func CreateAccount(strength uint8, language int) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.CreateNewAccountWithMnemonic(language, strength)
	if err != nil {
		log.Printf("CreateAccount CreateNewAccountWithMnemonic err: %v", err)
		return nil, err
	}

	account := &Account{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JsonPublicKey,
		PrivateKey: ecdsaAccount.JsonPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// RetrieveAccount 根据助记词恢复账户
//
// Parameters:
//   - `mnemonic`： 助记词，例如："玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即"。
//   - `language`： 1中文，2英文。
func RetrieveAccount(mnemonic string, language int) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.RetrieveAccountByMnemonic(mnemonic, language)
	if err != nil {
		return nil, err
	}
	account := &Account{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JsonPublicKey,
		PrivateKey: ecdsaAccount.JsonPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// CreateAndSaveAccountToFile 创建账户并保存到文件
//
// Parameters:
//   - `path`：保存路径。
//   - `passwd`： 密码。
//   - `strength`：助记词强度。
//   - `language`：助记词语言。
func CreateAndSaveAccountToFile(path, passwd string, strength uint8, language int) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.CreateNewAccountWithMnemonic(language, strength)
	if err != nil {
		return nil, err
	}

	err = common.PathExistsAndMkdir(path)
	if err != nil {
		return nil, err
	}

	_, err = cli.RetrieveAccountByMnemonicAndSavePrivKey(path, language, ecdsaAccount.Mnemonic, passwd)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JsonPublicKey,
		PrivateKey: ecdsaAccount.JsonPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// GetAccountFromPlainFile 通过账户私钥纯文件获取账户，包括账户的地址、私钥、公钥，恢复的账户不包括助记词。
//
// 指定路径下的结构如下:
//  - keys
//   |-- address
//   |-- private.key
//   |-- public.key
func GetAccountFromPlainFile(path string) (*Account, error) {
	addr, err := ioutil.ReadFile(path + "/address")
	if err != nil {
		log.Printf("GetAccountFromPlainFile error load address error = %v", err)
		return nil, err
	}
	pubkey, err := ioutil.ReadFile(path + "/public.key")
	if err != nil {
		log.Printf("GetAccountFromPlainFile error load pubkey error = %v", err)
		return nil, err
	}
	prikey, err := ioutil.ReadFile(path + "/private.key")
	if err != nil {
		log.Printf("GetAccountFromPlainFile error load prikey error = %v", err)
		return nil, err
	}

	account := &Account{
		Address:    string(addr),
		PublicKey:  string(pubkey),
		PrivateKey: string(prikey),
	}
	return account, nil
}

// GetAccountFromFile 通过密码以及账户私钥文件获取账户，不包括账户的助记词。
func GetAccountFromFile(path, passwd string) (*Account, error) {
	cryptoClient := crypto.GetCryptoClient()
	ecdsaPrivateKey, err := cryptoClient.GetEcdsaPrivateKeyFromFileByPassword(path, passwd)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	account.PrivateKey, err = cryptoClient.GetEcdsaPrivateKeyJsonFormatStr(ecdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	account.PublicKey, err = cryptoClient.GetEcdsaPublicKeyJsonFormatStr(ecdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	account.Address, err = cryptoClient.GetAddressFromPublicKey(&ecdsaPrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return account, err
}
