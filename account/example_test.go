// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// package account is related to account operation
package account

import (
	"fmt"
	"os"
)

func ExampleCreateAccount() {
	account, err := CreateAccount(1, 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(account)
}

func ExampleRetrieveAccount() {
	account, err := RetrieveAccount("玉 脸 驱 协 介 跨 尔 籍 杆 伏 愈 即", 1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(account)
}

func ExampleXchainToEVMAddress() {
	evmAddr, addrType, err := XchainToEVMAddress("dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if addrType == XchainAddrType {
		fmt.Println("地址为 xchain 普通账户地址类型")
	}

	fmt.Println(evmAddr)

	evmAddr, addrType, err = XchainToEVMAddress("XC1111111111111113@xuper")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if addrType == ContractAccountType {
		fmt.Println("地址为合约账户类型")
	}

	fmt.Println(evmAddr)
}
