// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package account is related to account operation.
// You can reate account and mnemonic, or save account private to file.
// You can set contract account for account if you want to deploy contract.
package account

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"

	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/crypto"
)

// Account account structure
type Account struct {
	contractAccount string

	Address    string
	PrivateKey string
	PublicKey  string
	Mnemonic   string
}

// CreateAccount create an account.
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

// RetrieveAccount retrieve account from mnemonic.
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

// CreateAndSaveAccountToFile create an account and save to file.
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

// GetAccountFromPlainFile import account from plain files which are JSON encoded
//
// 指定路径下的结构如下:
//  - keys
//   |-- address
//   |-- private.key
//   |-- public.key
func GetAccountFromPlainFile(path string) (*Account, error) {
	addr, err := ioutil.ReadFile(filepath.Join(path, "address"))
	if err != nil {
		log.Printf("GetAccountFromPlainFile error load address error = %v", err)
		return nil, err
	}
	pubkey, err := ioutil.ReadFile(filepath.Join(path, "public.key"))
	if err != nil {
		log.Printf("GetAccountFromPlainFile error load pubkey error = %v", err)
		return nil, err
	}
	prikey, err := ioutil.ReadFile(filepath.Join(path, "private.key"))
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

// GetAccountFromFile get an account from file and password.
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

// SetContractAccount set contract account.
// If you set contract account, this account represents the contract account.
// In some scenarios, must set contract account, such as deploy contract.
func (a *Account) SetContractAccount(contractAccount string) error {
	if ok, _ := regexp.MatchString(`^XC\d{16}@*`, contractAccount); !ok {
		return common.ErrInvalidContractAccount
	}

	a.contractAccount = contractAccount
	return nil
}

// RemoveContractAccount remove contract account from this account.
func (a *Account) RemoveContractAccount() {
	a.contractAccount = ""
}

// GetAuthRequire get this account's authRequire for transaction.
// If you set contract account, returns $ContractAccount+"/"+$Address, otherwise returns $Address.
func (a *Account) GetAuthRequire() string {
	if a.HasContractAccount() {
		return a.GetContractAccount() + "/" + a.Address
	}
	return a.Address
}

// GetContractAccount get current contract account, returns an empty string if the contract account is not set.
func (a *Account) GetContractAccount() string {
	return a.contractAccount
}

// HasContractAccount reutrn true if you set contract account, otherwise returns false.
func (a *Account) HasContractAccount() bool {
	return a.contractAccount != ""
}
