package xuper

import (
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuperchain/core/pb"
)

type Transaction struct {
	Tx *pb.Transaction
	// ComplianceCheckTx *pb.Transaction // 这个字段应该不用
	ContractResponse *pb.ContractResponse

	Fee     string
	GasUsed int64

	DigestHash []byte
}

func (t *Transaction) Sign(account *account.Account) error {
	// 把 account 添加到交易的 AuthRequire 中然后重新计算 digestHash，然后再签名。
	return nil
}

func (t *Transaction) CalcDigestHash() {

}
