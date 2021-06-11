package xuper

import (
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/crypto"

	"github.com/xuperchain/xuperchain/core/pb"
)

type Transaction struct {
	Tx               *pb.Transaction
	ContractResponse *pb.ContractResponse

	Fee     string
	GasUsed int64

	DigestHash []byte
}

// Sign account sign for tx, for mulitysign.
func (t *Transaction) Sign(account *account.Account) error {
	// 把 account 添加到交易的 AuthRequire 中然后重新计算 digestHash，然后再签名。
	if !inSlice(t.Tx.AuthRequire, account.GetAuthRequire()) {
		t.Tx.AuthRequire = append(t.Tx.AuthRequire, account.GetAuthRequire())
	}

	digestHash, err := common.MakeTxDigestHash(t.Tx)
	if err != nil {
		return err
	}

	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(account.PrivateKey)
	if err != nil {
		return err
	}

	sign, err := cryptoClient.SignECDSA(privateKey, digestHash)

	signatureInfo := &pb.SignatureInfo{
		PublicKey: account.PublicKey,
		Sign:      sign,
	}

	t.Tx.AuthRequireSigns = append(t.Tx.AuthRequireSigns, signatureInfo)

	// make txid
	t.Tx.Txid, err = common.MakeTransactionID(t.Tx)

	return err
}

func inSlice(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
