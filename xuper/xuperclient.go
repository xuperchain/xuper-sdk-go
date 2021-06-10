package xuper

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuperchain/core/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/encoding/gzip"
)

// XClient xuperchain client.
type XClient struct {
	node  string
	xc    pb.XchainClient
	xconn *grpc.ClientConn

	ec    pb.XendorserClient
	econn *grpc.ClientConn

	cfg *config.CommConfig
	opt *clientOptions
}

// New new xuper client.
func New(node string, opts ...ClientOption) (*XClient, error) {
	opt := &clientOptions{}
	for _, param := range opts {
		err := param(opt)
		if err != nil {
			return nil, fmt.Errorf("option failed: %v", err)
		}
	}

	xclient := &XClient{
		node: node,
		opt:  opt,
	}

	err := xclient.init()
	if err != nil {
		return nil, err
	}

	return xclient, nil
}

func (x *XClient) init() error {
	var err error

	if x.opt.configFile != "" {
		x.cfg, err = config.GetConfig(x.opt.configFile)
		if err != nil {
			return err
		}
	}

	// init xuper client, endorser client, grpc tls & gzip.
	return x.initConn()
}

func (x *XClient) initConn() error {
	grpcOpts := []grpc.DialOption{}

	if x.opt.grpcTLS.serverName != "" { // TLS enabled
		certificate, err := tls.LoadX509KeyPair(x.opt.grpcTLS.certFile, x.opt.grpcTLS.keyFile)
		if err != nil {
			log.Fatal(err)
		}

		certPool := x509.NewCertPool()
		ca, err := ioutil.ReadFile(x.opt.grpcTLS.cacertFile)
		if err != nil {
			return err
		}
		if ok := certPool.AppendCertsFromPEM(ca); !ok {
			return errors.New("certPool add ca cert failed")
		}

		creds := credentials.NewTLS(&tls.Config{
			Certificates: []tls.Certificate{certificate},
			ServerName:   x.opt.grpcTLS.serverName,
			RootCAs:      certPool,
		})

		grpcOpts = append(grpcOpts, grpc.WithTransportCredentials(creds))
	}

	if x.opt.useGrpcGZIP { // gzip enabled
		grpcOpts = append(grpcOpts, grpc.WithDefaultCallOptions(grpc.UseCompressor(gzip.Name)))
	}

	grpcOpts = append(grpcOpts, grpc.WithMaxMsgSize(64<<20-1))

	conn, err := grpc.Dial(
		x.node,
		grpcOpts...,
	)
	if err != nil {
		return err
	}

	x.xconn = conn
	x.xc = pb.NewXchainClient(conn)

	if x.cfg.ComplianceCheck.IsNeedComplianceCheck { // endorser no TLS, mayble future.
		econn, err := grpc.Dial(x.cfg.EndorseServiceHost, grpc.WithInsecure(), grpc.WithMaxMsgSize(64<<20-1))
		if err != nil {
			return err
		}
		x.econn = econn
		x.ec = pb.NewXendorserClient(econn)
	}

	return nil
}

// Close close xuper client all connections.
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
