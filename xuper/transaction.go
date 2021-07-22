package xuper

import (
	"errors"
	"strings"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/crypto"

	"github.com/xuperchain/xuperchain/service/pb"
)

// Transaction xuperchain transaction.
type Transaction struct {
	Tx               *pb.Transaction
	ContractResponse *pb.ContractResponse
	Bcname           string

	Fee     string
	GasUsed int64

	DigestHash []byte
}

// Sign account sign for tx, for multisign.multisign
func (t *Transaction) Sign(account *account.Account) error {
	if account == nil {
		return errors.New("Transaction sign account can not be nil")
	}
	// 对于多签，在交易预执行时就需要写好所有的需要签名的地址到 AuthRequire 字段，其他地址再进行签名时，需要检查是否已经在 AuthRequire 字段中。
	// 同时签名的顺序也要保持一致，不然上链时会失败。
	if !inSlice(t.Tx.AuthRequire, account.GetAuthRequire()) {
		return errors.New("this account not in transaction's AuthRequire list")
	}

	if t.DigestHash == nil {
		digestHash, err := common.MakeTxDigestHash(t.Tx)
		if err != nil {
			return err
		}
		t.DigestHash = digestHash
	}

	cryptoClient := crypto.GetCryptoClient()
	privateKey, err := cryptoClient.GetEcdsaPrivateKeyFromJsonStr(account.PrivateKey)
	if err != nil {
		return err
	}

	sign, err := cryptoClient.SignECDSA(privateKey, t.DigestHash)

	signatureInfo := &pb.SignatureInfo{
		PublicKey: account.PublicKey,
		Sign:      sign,
	}

	t.Tx.AuthRequireSigns = append(t.Tx.AuthRequireSigns, signatureInfo)
	t.Tx.InitiatorSigns = append(t.Tx.InitiatorSigns, signatureInfo)

	// make txid
	t.Tx.Txid, err = common.MakeTransactionID(t.Tx)

	return err
}

func inSlice(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}

		splitRes := strings.Split(v, "/")
		addr := splitRes[len(splitRes)-1]
		if addr == str {
			return true
		}
	}
	return false
}
