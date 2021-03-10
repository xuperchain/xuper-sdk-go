package transfer

import (
	"fmt"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func Example_transfer() {
	var (
		node   = ""
		bcname = "xuper"
	)

	from, _ := account.CreateAccount(1, 1)
	to, _ := account.CreateAccount(1, 1)

	trans := InitTrans(from, node, bcname)
	txID, err := trans.Transfer(to.Address, "888", "100", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("txID:", txID)
}

func Example_batchTrasfer() {
	var (
		node   = ""
		bcname = "xuper"
	)

	from, _ := account.CreateAccount(1, 1)
	trans := InitTrans(from, node, bcname)

	toAddrAmountMap := map[string]string{
		"tom": "111",
		"bob": "222",
	}
	txID, err := trans.BatchTransfer(toAddrAmountMap, "888", "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("txID:", txID)
}
