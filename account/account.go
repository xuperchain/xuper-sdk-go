// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package account is related to account operation
package account

import (
	"log"

	"github.com/xuperchain/xuper-sdk-go/common"
	"github.com/xuperchain/xuper-sdk-go/crypto"
)

// Account account structure
type Account struct {
	Address    string
	PrivateKey string
	PublicKey  string
	Mnemonic   string
}

// CreateAccount create an account
func CreateAccount(strength uint8, language int) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.CreateNewAccountWithMnemonic(language, strength)
	if err != nil {
		log.Printf("CreateAccount CreateNewAccountWithMnemonic err: %v", err)
		return nil, err
	}

	account := &Account{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JSONPublicKey,
		PrivateKey: ecdsaAccount.JSONPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// RetrieveAccount retrieve account from mnemonic
func RetrieveAccount(mnemonic string, language int) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaAccount, err := cli.RetrieveAccountByMnemonic(mnemonic, language)
	if err != nil {
		return nil, err
	}

	account := &Account{
		Address:    ecdsaAccount.Address,
		PublicKey:  ecdsaAccount.JSONPublicKey,
		PrivateKey: ecdsaAccount.JSONPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// CreateAndSaveAccountToFile create an account and save to file
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
		PublicKey:  ecdsaAccount.JSONPublicKey,
		PrivateKey: ecdsaAccount.JSONPrivateKey,
		Mnemonic:   ecdsaAccount.Mnemonic,
	}
	return account, nil
}

// GetAccountFromFile get an account from file
func GetAccountFromFile(path, passwd string) (*Account, error) {
	cli := crypto.GetCryptoClient()
	ecdsaPrivateKey, err := cli.GetEcdsaPrivateKeyFromFileByPassword(path, passwd)
	if err != nil {
		return nil, err
	}

	account := &Account{}
	account.PrivateKey, err = cli.GetEcdsaPrivateKeyJSONFormat(ecdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	account.PublicKey, err = cli.GetEcdsaPublicKeyJSONFormat(ecdsaPrivateKey)
	if err != nil {
		return nil, err
	}
	account.Address, err = cli.GetAddressFromPublicKey(&ecdsaPrivateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	return account, err
}
