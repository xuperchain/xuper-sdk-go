package common

import (
	"github.com/xuperchain/xuperchain/core/pb"
	"github.com/xuperchain/xuperchain/core/utxo/txhash"
)

// MakeTransactionID 计算交易ID，包括签名。
func MakeTransactionID(tx *pb.Transaction) ([]byte, error) {
	return txhash.MakeTransactionID(tx)
}

// MakeTxDigestHash 计算交易哈希，不包括签名。
func MakeTxDigestHash(tx *pb.Transaction) ([]byte, error) {
	return txhash.MakeTxDigestHash(tx)
}
