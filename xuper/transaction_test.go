package xuper

import (
	"testing"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuperchain/service/pb"
)

func TestTransaction(t *testing.T) {
	// acc, _ := account.CreateAccount(1, 1)
	acc1, _ := account.CreateAccount(1, 1)

	type Case struct {
		hasHash       bool
		signAcc       *account.Account
		hasErr        bool
		inAuthRequire bool
	}

	cases := []*Case{
		{
			hasHash:       true,
			signAcc:       nil,
			hasErr:        true,
			inAuthRequire: true,
		},
		{
			hasHash:       false,
			signAcc:       acc1,
			hasErr:        false,
			inAuthRequire: true,
		},
		{
			hasHash:       true,
			signAcc:       acc1,
			hasErr:        true,
			inAuthRequire: false,
		},
	}

	for _, c := range cases {
		tx := &Transaction{
			Tx: &pb.Transaction{},
		}
		if c.hasHash {
			tx.DigestHash = []byte("haha")
		}

		if c.inAuthRequire {
			tx.Tx.AuthRequire = append(tx.Tx.AuthRequire, acc1.GetAuthRequire())
		}

		err := tx.Sign(c.signAcc)
		if c.hasErr {
			if err == nil {
				t.Error("Transactions assert failed1")
			}
		} else {
			if err != nil {
				t.Error("Transactions assert failed2")
			}
		}
	}

}
