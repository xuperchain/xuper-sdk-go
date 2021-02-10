package account

import (
	"testing"
)

func TestTrans(t *testing.T) {
	type Case struct {
		XchainAddress string
		Type          string
		EVMAddress    string
	}

	cases := []Case{
		{
			XchainAddress: "XC1111111111111113@xuper",
			Type:          "contract-account",
			EVMAddress:    "3131313231313131313131313131313131313133",
		},
		{
			XchainAddress: "dpzuVdosQrF2kmzumhVeFQZa1aYcdgFpN",
			Type:          "xchain",
			EVMAddress:    "93F86A462A3174C7AD1281BCF400A9F18D244E06",
		},
		{
			XchainAddress: "storagedata11",
			Type:          "contract-name",
			EVMAddress:    "313131312D2D2D73746F72616765646174613131",
		},
	}

	for _, c := range cases {
		evmAddr, addrType, err := XchainToEVMAddress(c.XchainAddress)
		if err != nil {
			t.Error(err)
		}
		Assert(c.EVMAddress, evmAddr, t)
		Assert(c.Type, addrType, t)

		xAddr, addrType, err := EVMToXchainAddress(c.EVMAddress)
		if err != nil {
			t.Error(err)
		}
		Assert(c.XchainAddress, xAddr, t)
		Assert(c.Type, addrType, t)
	}
}

func Assert(expect, acture string, t *testing.T) {
	if expect != acture {
		t.Errorf("expect: %s, acture: %s.", expect, acture)
	}
}
