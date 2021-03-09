package contractaccount

import (
	"fmt"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func Example() {
	var (
		node   = "127.0.0.1:37101"
		bcname = "xuper"
	)

	acc, _ := account.CreateAccount(1, 1)

	ca := InitContractAccount(acc, node, bcname)

	txID, err := ca.CreateContractAccount("XC123456789012345@xuper")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Println("TxID:", txID)
}
