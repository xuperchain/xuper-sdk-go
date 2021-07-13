// Copyright (c) 2019. Baidu Inc. All Rights Reserved.

// Package xuper xuperchain client and generate tx or post tx.
//
// You can transfer to someone, deploy/invoke/query contract(wasm, native and evm),
// create contract account or set contract account ACL or set contract method ACL, query info from node.
//
// If you need multisign, you can use Transaction sign method add signature, more example in xuper-sdk-go/example/.
package xuper

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"
	"regexp"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuperchain/service/pb"

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
	esc   pb.EventServiceClient
	econn *grpc.ClientConn

	cfg *config.CommConfig
	opt *clientOptions
}

// New new xuper client.
//
// Parameters:
//   - `node`: node GRPC URL.
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
	} else {
		x.cfg = config.GetInstance()
	}

	// init xuper client, endorser client, grpc tls & gzip.
	return x.initConn()
}

func (x *XClient) initConn() error {
	grpcOpts := []grpc.DialOption{}

	if x.opt.grpcTLS != nil && x.opt.grpcTLS.serverName != "" { // TLS enabled
		certificate, err := tls.LoadX509KeyPair(x.opt.grpcTLS.certFile, x.opt.grpcTLS.keyFile)
		if err != nil {
			return err
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
	} else {
		grpcOpts = append(grpcOpts, grpc.WithInsecure())
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
	x.esc = pb.NewEventServiceClient(conn)

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

// DeployNativeGoContract deploy native go contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `code`: Contract code bytes.
//   - `args`: Contract init args.
func (x *XClient) DeployNativeGoContract(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewDeployContractRequest(from, name, nil, code, args, NativeContractModule, GoRuntime, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// DeployNativeJavaContract deploy native java contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `code`: Contract code bytes.
//   - `args`: Contract init args.
func (x *XClient) DeployNativeJavaContract(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewDeployContractRequest(from, name, nil, code, args, NativeContractModule, JavaRuntime, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// DeployWasmContract deploy wasm c++ contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `code`: Contract code bytes.
//   - `args`: Contract init args.
func (x *XClient) DeployWasmContract(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewDeployContractRequest(from, name, nil, code, args, WasmContractModule, CRuntime, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// DeployEVMContract deploy evm contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `abi` : Solidity contract abi.
//   - `bin` : Solidity contract bin.
//   - `args`: Contract init args.
func (x *XClient) DeployEVMContract(from *account.Account, name string, abi, bin []byte, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewDeployContractRequest(from, name, abi, bin, args, EvmContractModule, "", opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// UpgradeWasmContract upgrade wasm contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `code`: Contract code bytes.
//   - `args`: Contract init args.
func (x *XClient) UpgradeWasmContract(from *account.Account, name string, code []byte, opts ...RequestOption) (*Transaction, error) {
	req, err := NewUpgradeContractRequest(from, WasmContractModule, name, code, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// UpgradeNativeContract upgrade native contract.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `name`: Contract name.
//   - `code`: Contract code bytes.
//   - `args`: Contract init args.
func (x *XClient) UpgradeNativeContract(from *account.Account, name string, code []byte, opts ...RequestOption) (*Transaction, error) {
	req, err := NewUpgradeContractRequest(from, NativeContractModule, name, code, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// InvokeWasmContract invoke wasm c++ contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) InvokeWasmContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, WasmContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// InvokeNativeContract invoke native contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) InvokeNativeContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, NativeContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// InvokeEVMContract invoke evm contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) InvokeEVMContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, EvmContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// QueryWasmContract query wasm c++ contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) QueryWasmContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, WasmContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}

	return x.PreExecTx(req)
}

// QueryNativeContract query native contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) QueryNativeContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, NativeContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}
	return x.PreExecTx(req)
}

// QueryEVMContract query evm contract.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `args`  : Contract invoke args.
func (x *XClient) QueryEVMContract(from *account.Account, name, method string, args map[string]string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewInvokeContractRequest(from, EvmContractModule, name, method, args, opts...)
	if err != nil {
		return nil, err
	}
	return x.PreExecTx(req)
}

// Transfer to another address.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `to`    : Transfer receiving address.
//   - `amount`: Transfer amount.
func (x *XClient) Transfer(from *account.Account, to, amount string, opts ...RequestOption) (*Transaction, error) {
	req, err := NewTransferRequest(from, to, amount, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// CreateContractAccount create contract account for initiator.
//
// Parameters:
//   - `from`           : Transaction initiator. NOTE: from must be NOT set contract account, if you set please remove it.
//   - `contractAccount`:The contract account you want to create, such as: XC8888888899999999@xuper.
func (x *XClient) CreateContractAccount(from *account.Account, contractAccount string, opts ...RequestOption) (*Transaction, error) {
	if ok, _ := regexp.MatchString(`^XC\d{16}@*`, contractAccount); !ok {
		return nil, common.ErrInvalidContractAccount
	}

	subRegexp := regexp.MustCompile(`\d{16}`)
	contractAccountByte := subRegexp.Find([]byte(contractAccount))
	contractAccount = string(contractAccountByte)
	req, err := NewCreateContractAccountRequest(from, contractAccount, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// SetAccountACL update contract account acl. NOTE: from account must be set contract account.
//
// Parameters:
//   - `from`: Transaction initiator.
//   - `acl` : The ACL you want to set.
func (x *XClient) SetAccountACL(from *account.Account, acl *ACL, opts ...RequestOption) (*Transaction, error) {
	req, err := NewSetAccountACLRequest(from, acl, opts...)
	if err != nil {
		return nil, err
	}
	return x.Do(req)
}

// SetMethodACL update contract method acl.
//
// Parameters:
//   - `from`  : Transaction initiator.
//   - `name`  : Contract name.
//   - `method`: Contract method.
//   - `acl`   : The ACL you want to set.
func (x *XClient) SetMethodACL(from *account.Account, name, method string, acl *ACL, opts ...RequestOption) (*Transaction, error) {
	req, err := NewSetMethodACLRequest(from, name, method, acl, opts...)
	if err != nil {
		return nil, err
	}

	return x.Do(req)
}

// Do generete tx & post tx.
func (x *XClient) Do(req *Request) (*Transaction, error) {
	transaction, err := x.GenerateTx(req)
	if err != nil {
		return nil, err
	}

	// build transaction only.
	if req.opt.notPost {
		return transaction, nil
	}

	// post tx.
	return x.PostTx(transaction)
}

// GenerateTx generate Transaction.
func (x *XClient) GenerateTx(req *Request) (*Transaction, error) {
	proposal, err := NewProposal(x, req, x.cfg)
	if err != nil {
		return nil, err
	}
	return proposal.Build()
}

// PreExecTx preExec for query.
func (x *XClient) PreExecTx(req *Request) (*Transaction, error) {
	proposal, err := NewProposal(x, req, x.cfg)
	if err != nil {
		return nil, err
	}
	err = proposal.PreExecWithSelectUtxo()
	if err != nil {
		return nil, err
	}

	var cr *pb.ContractResponse
	if len(proposal.preResp.GetResponse().GetResponses()) > 0 {
		cr = proposal.preResp.GetResponse().GetResponses()[len(proposal.preResp.GetResponse().GetResponses())-1]
	}

	return &Transaction{
		ContractResponse: cr,
	}, nil
}

// PostTx post tx to node.
func (x *XClient) PostTx(tx *Transaction) (*Transaction, error) {
	return tx, x.postTx(tx.Tx, tx.Bcname)
}

// WatchBlockEvent new watcher for block event.
func (x *XClient) WatchBlockEvent(opts ...BlockEventOption) (*Watcher, error) {
	watcher, err := x.newWatcher(opts...)
	if err != nil {
		return nil, err
	}
	buf, _ := proto.Marshal(watcher.opt.blockFilter)
	request := &pb.SubscribeRequest{
		Type:   pb.SubscribeType_BLOCK,
		Filter: buf,
	}

	stream, err := x.esc.Subscribe(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	filteredBlockChan := make(chan *FilteredBlock, watcher.opt.blockChanBufferSize)
	exit := make(chan struct{})
	watcher.exit = exit
	watcher.FilteredBlockChan = filteredBlockChan

	go func() {
		defer func() {
			close(filteredBlockChan)
			if err := stream.CloseSend(); err != nil {
				log.Printf("Unregister block event failed, close stream error: %v", err)
			} else {
				log.Printf("Unregister block event success...")
			}
		}()
		for {
			select {
			case <-exit:
				return
			default:
				event, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("Get block event err: %v", err)
					return
				}
				var block pb.FilteredBlock
				err = proto.Unmarshal(event.Payload, &block)
				if err != nil {
					log.Printf("Get block event err: %v", err)
					return
				}
				if len(block.GetTxs()) == 0 && watcher.opt.skipEmptyTx {
					continue
				}
				filteredBlockChan <- fromFilteredBlockPB(&block)
			}
		}
	}()
	return watcher, nil
}

func (x *XClient) newWatcher(opts ...BlockEventOption) (*Watcher, error) {
	opt, err := initEventOpts(opts...)
	if err != nil {
		return nil, err
	}

	watcher := &Watcher{
		opt: opt,
	}
	return watcher, nil
}

func (x *XClient) postTx(tx *pb.Transaction, bcname string) error {
	ctx := context.Background()
	c := x.xc
	txStatus := &pb.TxStatus{
		Bcname: bcname,
		Status: pb.TransactionStatus_UNCONFIRM,
		Tx:     tx,
		Txid:   tx.Txid,
	}
	res, err := c.PostTx(ctx, txStatus)
	if err != nil {
		return errors.Wrap(err, "xuperclient post tx failed")
	}
	if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return fmt.Errorf("Failed to post tx: %s", res.Header.Error.String())
	}
	return nil
}

// QueryTxByID query the tx by txID
//
// Parameters
//  - `txID` : transaction id
func (x *XClient) QueryTxByID(txID string, opts ...QueryOption) (*pb.Transaction, error) {
	return x.queryTxByID(txID, opts...)
}

// QueryBlockByID query the block by blockID
//
// Parameters:
//   - `blockID`  : block id
func (x *XClient) QueryBlockByID(blockID string, opts ...QueryOption) (*pb.Block, error) {
	return x.queryBlockByID(blockID, opts...)
}

// QueryBlockByHeight query the block by block height
//
// Parameters:
//   - `height`  : block height
func (x *XClient) QueryBlockByHeight(height int64, opts ...QueryOption) (*pb.Block, error) {
	return x.queryBlockByHeight(height, opts...)
}

// QueryAccountACL query the ACL by account
//
// Parameters:
//   - `account`  : account, such as XC1111111111111111@xuper
func (x *XClient) QueryAccountACL(account string, opts ...QueryOption) (*ACL, error) {
	return x.queryAccountACL(account, opts...)
}

// QueryMethodACL query the ACL by method
//
// Parameters:
//   - `name`     : contract name
//   - `account`  : account
func (x *XClient) QueryMethodACL(name, method string, opts ...QueryOption) (*ACL, error) {
	return x.queryMethodACL(name, method, opts...)
}

// QueryAccountContracts query all contracts for account
//
// Parameters:
//   - `account`  : account,such as XC1111111111111111@xuper
func (x *XClient) QueryAccountContracts(account string, opts ...QueryOption) ([]*pb.ContractStatus, error) {
	return x.queryAccountContracts(account, opts...)
}

// QueryAddressContracts query all contracts for address
//
// Parameters:
//   - `address`  : address
//
// Returns:
//   - `map`  : contractAccount => contractStatusList
//   - `error`: error
func (x *XClient) QueryAddressContracts(address string, opts ...QueryOption) (map[string]*pb.ContractList, error) {
	return x.queryAddressContracts(address, opts...)
}

// QueryBalance query balance by the address
//
// Parameters:
//   - `address`  : address
func (x *XClient) QueryBalance(address string, opts ...QueryOption) (*big.Int, error) {
	return x.queryBalance(address, opts...)
}

// QueryBalanceDetail query the balance detail by address
//
// Parameters:
//   - `address`  : address
func (x *XClient) QueryBalanceDetail(address string, opts ...QueryOption) ([]*BalanceDetail, error) {
	return x.queryBalanceDetail(address, opts...)
}

// QuerySystemStatus query the system status
func (x *XClient) QuerySystemStatus(opts ...QueryOption) (*pb.SystemsStatusReply, error) {
	return x.querySystemStatus(opts...)
}

// QueryBlockChains query block chains
func (x *XClient) QueryBlockChains(opts ...QueryOption) ([]string, error) {
	return x.queryBlockChains(opts...)
}

// QueryBlockChainStatus query the block chain status
//
// Parameters:
//   - `chainName`  : chainName
func (x *XClient) QueryBlockChainStatus(chainName string, opts ...QueryOption) (*pb.BCStatus, error) {
	return x.queryBlockChainStatus(chainName)
}

// QueryNetURL query the net URL
func (x *XClient) QueryNetURL(opts ...QueryOption) (string, error) {
	return x.queryNetURL(opts...)
}

// QueryAccountByAK query the account  by AK
//
// Parameters:
//   - `address`  : address
func (x *XClient) QueryAccountByAK(address string, opts ...QueryOption) ([]string, error) {
	return x.queryAccountByAK(address, opts...)
}
