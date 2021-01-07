package account

import (
	"fmt"
	"testing"
)

func TestX2E(t *testing.T) {
	aa := []string{
		"XC1111111111111113@xuper",
		"dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN",
		"storagedata11",
	}

	for _, v := range aa {
		fmt.Println("addr:", v)
		a, ty, e := XchainToEVMAddress(v)
		fmt.Println("aa1:", a)
		fmt.Println("ty:", ty)
		fmt.Println("e:", e)
	}
}

func TestE2X(t *testing.T) {
	aa := []string{
		"3131313231313131313131313131313131313133",
		"93F86A462A3174C7AD1281BCF400A9F18D244E06",
		"313131312D2D2D73746F72616765646174613131",
	}

	for _, v := range aa {
		fmt.Println("addr:", v)
		a, ty, e := EVMToXchainAddress(v)
		fmt.Println("aa1:", a)
		fmt.Println("ty:", ty)
		fmt.Println("e:", e)
	}
}
