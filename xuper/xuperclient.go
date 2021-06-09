package xuper

import (
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuperchain/core/pb"
	"google.golang.org/grpc"
)

// XClient xuperchain client.
type XClient struct {
	node  string
	xc    *pb.XchainClient
	xconn *grpc.ClientConn

	ec    *pb.XendorserClient
	econn *grpc.ClientConn
	// cfg *Config
	opt *clientOptions
}

// New new
func New(node string, opts ...ClientOption) *XClient {

	return &XClient{}
}

func (x *XClient) Close() error {
	if x.xc != nil && x.xconn != nil {
		err := x.xconn.Close()
		if err != nil {
			return err
		}
	}

	if x.ec != nil && x.econn != nil {
		err := x.econn.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

// Transfer t
func (x *XClient) Transfer(from *account.Account, to, amount string, opts ...TxOption) (*Transaction, error) {

	return nil, nil
}

// GenerateTx g
func (x *XClient) GenerateTx(req *Request) (*Transaction, error) {
	return nil, nil
}

// PostTx p
func (x *XClient) PostTx(tx *Transaction) (*Transaction, error) {
	return nil, nil
}

func (x *XClient) RegisterBlockEvent(filter *pb.BlockFilter, skipEmptyTx bool) (*Registration, error) {
	return nil, nil
}

func (x *XClient) Unregister(r *Registration) {
	r.Unregister()
}

func (x *XClient) QueryBalance(address string) string {
	return ""
}
