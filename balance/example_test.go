package balance

import (
	"fmt"
	"os"

	"github.com/xuperchain/xuper-sdk-go/account"
)

func Example() {
	var (
		node    = "127.0.0.1:37101"
		baNames = []string{"xuper"}
	)

	acc, _ := account.CreateAccount(1, 1)

	b := InitBalance(acc, node, baNames)

	details, err := b.GetBalanceDetails()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println(details)
}
