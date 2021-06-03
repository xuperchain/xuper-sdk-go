package xuper

import (
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuperchain/core/pb"
)

type Transaction struct {
	Tx               *pb.Transaction
	ContractResponse *pb.ContractResponse

	Fee     string
	GasUsed string

	DigestHash []byte
}

func (t *Transaction) Sign(account *account.Account) error {

	return nil
}
