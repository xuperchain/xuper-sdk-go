package main

import (
	"fmt"
	"log"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/transfer"
)

// define blockchain node and blockchain name
var (
	contractName = "counter"

	// node for test network of XuperOS
	// node = "14.215.179.74:37101"

	// node for official network of XuperOS
	node = "39.156.69.83:37100"

	//	node         = "127.0.0.1:37801"
	bcname = "xuper"
)

func createAccount() (string, error) {
	// create an account for the user,
	// strength 1 means that the number of mnemonics is 12
	// language 1 means that mnemonics is Chinese
	acc, err := account.CreateAccount(1, 1)
	if err != nil {
		return "", fmt.Errorf("create account error: %v\n", err)
	}
	// print the account, mnemonics
	fmt.Println(acc)
	fmt.Println(acc.Mnemonic)

	return acc.Mnemonic, nil
}

func getBalance(mnemonic string) {
	// retrieve the account by mnemonics
	acc, err := account.RetrieveAccount(mnemonic, 1)
	if err != nil {
		fmt.Printf("retrieveAccount err: %v\n", err)
	}
	fmt.Printf("account: %v\n", acc)

	// initialize a client to operate the transaction
	trans := transfer.InitTrans(acc, node, bcname)

	// get balance of the account
	balance, err := trans.GetBalance()
	log.Printf("balance %v, err %v", balance, err)
	return
}

func main() {
	mm, err := createAccount()
	if err != nil {
		os.Exit(-1)
	}
	getBalance(mm)
	return
}
