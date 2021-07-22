package xuper

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common/config"
	"github.com/xuperchain/xuperchain/service/pb"
	"google.golang.org/grpc"
)

func TestBuild(t *testing.T) {
	xc := &XClient{
		xc: &MockXClient{},
		cfg: &config.CommConfig{
			TxVersion: 1,
			ComplianceCheck: config.ComplianceCheckConfig{
				IsNeedComplianceCheck: false,
			},
		},
	}

	acc, _ := account.CreateAccount(1, 1)
	a, e := xc.Transfer(acc, "a", "10")
	if e != nil {
		t.Error(e)
	} else {
		t.Log(a)
	}
}
func TestNewProposal(t *testing.T) {
	type Case struct {
		xclient        *XClient
		request        *Request
		cfg            *config.CommConfig
		expectError    error
		expectProposal *Proposal
	}

	cases := []Case{
		{
			xclient:        nil,
			request:        nil,
			cfg:            nil,
			expectError:    errors.New("new proposal failed, parameters can not be nil"),
			expectProposal: nil,
		},
		{
			xclient:        &XClient{},
			request:        nil,
			cfg:            nil,
			expectError:    errors.New("new proposal failed, parameters can not be nil"),
			expectProposal: nil,
		},
		{
			xclient:        nil,
			request:        &Request{},
			cfg:            nil,
			expectError:    errors.New("new proposal failed, parameters can not be nil"),
			expectProposal: nil,
		},
		{
			xclient:        &XClient{},
			request:        &Request{},
			cfg:            &config.CommConfig{},
			expectError:    nil,
			expectProposal: &Proposal{txVersion: 3},
		},
		{
			xclient: &XClient{},
			request: &Request{},
			cfg: &config.CommConfig{
				TxVersion: 1,
				ComplianceCheck: config.ComplianceCheckConfig{
					IsNeedComplianceCheck: true,
				},
			},
			expectError:    nil,
			expectProposal: &Proposal{txVersion: 1},
		},
	}

	for _, c := range cases {
		p, err := NewProposal(c.xclient, c.request, c.cfg)
		if c.expectError == nil {
			if err != nil {
				t.Error("new proposal assert error failed")
			}
			if p == nil {
				t.Error("new proposal assert proposal failed")
			}
			if p.txVersion != c.expectProposal.txVersion {
				t.Error("new proposal assert proposal tx version failed")
			}

		} else {
			if err.Error() != c.expectError.Error() {
				t.Error("new proposal assert error failed")
			}
			if p != c.expectProposal {
				t.Error("new proposal assert proposal failed")
			}
		}
	}
}

func TestGetInitiator(t *testing.T) {
	type Case struct {
		proposal        *Proposal
		expectInitiator string
	}

	acc := &account.Account{
		Address: "abc",
	}
	acc.SetContractAccount("XC1234567887654321@xuper")

	cases := []Case{
		{
			proposal: &Proposal{
				request: &Request{
					initiatorAccount: &account.Account{
						Address: "abc",
					},
					opt: &requestOptions{
						onlyFeeFromAccount: false,
					},
				},
				cfg: &config.CommConfig{
					TxVersion: 1,
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck: false,
					},
				},
			},
			expectInitiator: "abc",
		},
		{
			proposal: &Proposal{
				request: &Request{
					initiatorAccount: &account.Account{
						Address: "abc",
					},
					opt: &requestOptions{
						onlyFeeFromAccount: false,
					},
				},
				cfg: &config.CommConfig{
					TxVersion: 1,
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck: true,
					},
				},
			},
			expectInitiator: "abc",
		},
		{
			proposal: &Proposal{
				request: &Request{
					initiatorAccount: acc,
					opt: &requestOptions{
						onlyFeeFromAccount: false,
					},
				},
				cfg: &config.CommConfig{
					TxVersion: 1,
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck: false,
					},
				},
			},
			expectInitiator: "XC1234567887654321@xuper",
		},
	}

	for _, c := range cases {
		if c.expectInitiator != c.proposal.getInitiator() {
			t.Error("getInitiator assert failed")
		}
	}
}

func TestCalcTotalAmount(t *testing.T) {
	type Case struct {
		proposal     *Proposal
		expectAmount int64
		expectError  string
	}

	cases := []Case{
		{
			proposal: &Proposal{
				request: &Request{
					transferAmount: "a",
				},
			},
			expectAmount: 0,
			expectError:  `strconv.ParseInt: parsing "a": invalid syntax`,
		},
		{
			proposal: &Proposal{
				request: &Request{
					transferAmount: "0",
					opt: &requestOptions{
						onlyFeeFromAccount: false,
						fee:                "a",
					},
				},
			},
			expectAmount: 0,
			expectError:  `strconv.ParseInt: parsing "a": invalid syntax`,
		},
		{
			proposal: &Proposal{
				request: &Request{
					transferAmount: "0",
					opt: &requestOptions{
						onlyFeeFromAccount:   false,
						fee:                  "0",
						contractInvokeAmount: "a",
					},
				},
			},
			expectAmount: 0,
			expectError:  `strconv.ParseInt: parsing "a": invalid syntax`,
		},
		{
			proposal: &Proposal{
				request: &Request{
					transferAmount: "0",
					opt: &requestOptions{
						onlyFeeFromAccount:   false,
						fee:                  "0",
						contractInvokeAmount: "0",
					},
				},
				cfg: &config.CommConfig{
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck:            true,
						IsNeedComplianceCheckFee:         true,
						ComplianceCheckEndorseServiceFee: 1,
					},
				},
			},
			expectAmount: 1,
		},
		{
			proposal: &Proposal{
				request: &Request{
					transferAmount: "3",
					opt: &requestOptions{
						onlyFeeFromAccount:   false,
						fee:                  "10",
						contractInvokeAmount: "8",
					},
				},
				cfg: &config.CommConfig{
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck:            true,
						IsNeedComplianceCheckFee:         true,
						ComplianceCheckEndorseServiceFee: 1,
					},
				},
			},
			expectAmount: 22,
		},
	}

	for _, c := range cases {
		amount, err := c.proposal.calcTotalAmount()
		if err != nil {
			if err.Error() != c.expectError {
				t.Error("calcTotalAmount err assert failed")
			}
		} else {
			if c.expectAmount != amount {
				t.Errorf("calcTotalAmount amount assert failed: expect: %d, acture:%d", c.expectAmount, amount)
			}
		}
	}
}

func TestGenInvokeRPCRequest(t *testing.T) {
	type Case struct {
		proposal            *Proposal
		expectBcname        string
		expectInitiator     string
		expectAuthRequires  []string
		expectInvokeRequest pb.InvokeRequest
		expectError         bool
	}

	cases := []Case{
		{
			proposal: &Proposal{
				request: &Request{
					module: "xkernel",
					initiatorAccount: &account.Account{
						Address: "hello",
					},
					opt: &requestOptions{
						otherAuthRequire: []string{"a", "b"},
					},
				},
				cfg: &config.CommConfig{
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck:             false,
						IsNeedComplianceCheckFee:          true,
						ComplianceCheckEndorseServiceFee:  1,
						ComplianceCheckEndorseServiceAddr: "word",
					},
				},
			},
			expectInvokeRequest: pb.InvokeRequest{},
			expectBcname:        "xuper",
			expectInitiator:     "hello",
			expectAuthRequires:  []string{"hello", "a", "b"},
		},
		{
			proposal: &Proposal{
				request: &Request{
					module: "xkernel",
					initiatorAccount: &account.Account{
						Address: "hello",
					},
					opt: &requestOptions{
						otherAuthRequire: []string{"a", "b"},
					},
				},
				cfg: &config.CommConfig{
					ComplianceCheck: config.ComplianceCheckConfig{
						IsNeedComplianceCheck:             true,
						IsNeedComplianceCheckFee:          true,
						ComplianceCheckEndorseServiceFee:  1,
						ComplianceCheckEndorseServiceAddr: "word",
					},
				},
			},
			expectInvokeRequest: pb.InvokeRequest{},
			expectBcname:        "xuper",
			expectInitiator:     "hello",
			expectAuthRequires:  []string{"word", "hello", "a", "b"},
		},
	}

	for _, c := range cases {
		r, err := c.proposal.genInvokeRPCRequest()
		if err != nil {
			if !c.expectError {
				t.Error(err)
			}
		} else {
			if r == nil {
				continue
			}
			if c.expectBcname != r.Bcname {
				t.Error("genInvokeRPCRequest assert bcname failed")
			}
			if c.expectInitiator != r.Initiator {
				t.Error("genInvokeRPCRequest assert bcname failed")
			}
			for i, v := range r.AuthRequire {
				if c.expectAuthRequires[i] != v {
					t.Error("genInvokeRPCRequest assert AuthRequire failed")
				}
			}
		}
	}
}

type MockXchainClient interface {
	pb.XchainClient
}

var _ MockXchainClient = new(MockXClient)

// 实现 MockXchainClient 接口，
type MockXClient struct {
}

// SelectUTXOBySize merge many utxos into a few of utxos
func (mcx *MockXClient) SelectUTXOBySize(ctx context.Context, in *pb.UtxoInput, opts ...grpc.CallOption) (*pb.UtxoOutput, error) {
	return nil, nil
}

// PostTx post Transaction to a node
func (mcx *MockXClient) PostTx(ctx context.Context, in *pb.TxStatus, opts ...grpc.CallOption) (*pb.CommonReply, error) {
	return &pb.CommonReply{
		Header: newHeader(),
	}, nil
}
func (mcx *MockXClient) QueryACL(ctx context.Context, in *pb.AclStatus, opts ...grpc.CallOption) (*pb.AclStatus, error) {
	return &pb.AclStatus{
		Header:       newHeader(),
		AccountName:  in.GetAccountName(),
		ContractName: in.GetContractName(),
		MethodName:   in.GetMethodName(),
		Acl: &pb.Acl{
			Pm:        &pb.PermissionModel{Rule: 1, AcceptValue: 1.0},
			AksWeight: map[string]float64{"a": 1.0},
		},
	}, nil
}
func (mcx *MockXClient) QueryUtxoRecord(ctx context.Context, in *pb.UtxoRecordDetail, opts ...grpc.CallOption) (*pb.UtxoRecordDetail, error) {
	return nil, nil
}
func (mcx *MockXClient) QueryContractStatData(ctx context.Context, in *pb.ContractStatDataRequest, opts ...grpc.CallOption) (*pb.ContractStatDataResponse, error) {
	return nil, nil
}
func (mcx *MockXClient) GetAccountContracts(ctx context.Context, in *pb.GetAccountContractsRequest, opts ...grpc.CallOption) (*pb.GetAccountContractsResponse, error) {
	return &pb.GetAccountContractsResponse{
		Header: newHeader(),
	}, nil
}

// QueryTx query Transaction by TxStatus,
// Bcname and Txid are required for this
func (mcx *MockXClient) QueryTx(ctx context.Context, in *pb.TxStatus, opts ...grpc.CallOption) (*pb.TxStatus, error) {
	return &pb.TxStatus{
		Header: newHeader(),
		Txid:   in.GetTxid(),
		Tx: &pb.Transaction{
			Txid: in.GetTxid(),
		},
	}, nil
}

// GetBalance get balance of an address,
// Address is required for this
func (mcx *MockXClient) GetBalance(ctx context.Context, in *pb.AddressStatus, opts ...grpc.CallOption) (*pb.AddressStatus, error) {
	return &pb.AddressStatus{
		Header:  newHeader(),
		Address: in.GetAddress(),
		Bcs:     []*pb.TokenDetail{{Bcname: "xuper", Balance: "100"}},
	}, nil
}

// GetFrozenBalance get two kinds of balance
// 1. Still be frozen of an address
// 2. Available now of an address
// Address is required for this
func (mcx *MockXClient) GetBalanceDetail(ctx context.Context, in *pb.AddressBalanceStatus, opts ...grpc.CallOption) (*pb.AddressBalanceStatus, error) {
	return &pb.AddressBalanceStatus{
		Header:  newHeader(),
		Address: in.GetAddress(),
		Tfds:    []*pb.TokenFrozenDetails{{Bcname: "xuper", Tfd: []*pb.TokenFrozenDetail{{Balance: "100"}}}},
	}, nil
}

// GetFrozenBalance get balance that still be frozen of an address,
// Address is required for this
func (mcx *MockXClient) GetFrozenBalance(ctx context.Context, in *pb.AddressStatus, opts ...grpc.CallOption) (*pb.AddressStatus, error) {
	return &pb.AddressStatus{
		Header:  newHeader(),
		Address: in.GetAddress(),
		Bcs:     []*pb.TokenDetail{{Bcname: "xuper", Balance: "100"}},
	}, nil
}

// GetBlock get block by blockid and return if the block in trunk or in branch
func (mcx *MockXClient) GetBlock(ctx context.Context, in *pb.BlockID, opts ...grpc.CallOption) (*pb.Block, error) {
	return &pb.Block{
		Header:  newHeader(),
		Blockid: in.GetBlockid(),
		Block: &pb.InternalBlock{
			Blockid: in.GetBlockid(),
		},
	}, nil
}

// GetBlockByHeight get block by height and return if the block in trunk or in
// branch
func (mcx *MockXClient) GetBlockByHeight(ctx context.Context, in *pb.BlockHeight, opts ...grpc.CallOption) (*pb.Block, error) {
	return &pb.Block{
		Header:  newHeader(),
		Blockid: []byte("aa"),
		Block: &pb.InternalBlock{
			Blockid: []byte("aa"),
			Height:  in.GetHeight(),
		},
	}, nil
}
func (mcx *MockXClient) GetBlockChainStatus(ctx context.Context, in *pb.BCStatus, opts ...grpc.CallOption) (*pb.BCStatus, error) {
	return &pb.BCStatus{
		Header: newHeader(),
		Bcname: in.GetBcname(),
		Block: &pb.InternalBlock{
			Blockid: []byte("aa"),
			Height:  188,
		},
	}, nil
}

// Get blockchains query blockchains
func (mcx *MockXClient) GetBlockChains(ctx context.Context, in *pb.CommonIn, opts ...grpc.CallOption) (*pb.BlockChains, error) {
	return &pb.BlockChains{
		Header:      newHeader(),
		Blockchains: []string{"xuper"},
	}, nil
}

// GetSystemStatus query system status
func (mcx *MockXClient) GetSystemStatus(ctx context.Context, in *pb.CommonIn, opts ...grpc.CallOption) (*pb.SystemsStatusReply, error) {
	return &pb.SystemsStatusReply{
		Header:        newHeader(),
		SystemsStatus: &pb.SystemsStatus{},
	}, nil
}

func (mcx *MockXClient) GetConsensusStatus(ctx context.Context, in *pb.ConsensusStatRequest, opts ...grpc.CallOption) (*pb.ConsensusStatus, error) {
	return nil, nil
}

// GetNetURL return net url
func (mcx *MockXClient) GetNetURL(ctx context.Context, in *pb.CommonIn, opts ...grpc.CallOption) (*pb.RawUrl, error) {
	return nil, nil
}

// 新的Select utxos接口, 不需要签名，可以支持选择账户的utxo
func (mcx *MockXClient) SelectUTXO(ctx context.Context, in *pb.UtxoInput, opts ...grpc.CallOption) (*pb.UtxoOutput, error) {
	need, ok := big.NewInt(0).SetString(in.GetTotalNeed(), 10)
	if !ok {
		return nil, errors.New("invalid totalNeed")
	}
	need.Add(need, big.NewInt(100))

	a := big.NewInt(0).SetBytes(need.Bytes()).SetInt64(10)

	utxoList := []*pb.Utxo{
		{
			Amount:  big.NewInt(10).Bytes(),
			ToAddr:  []byte(in.Address),
			RefTxid: []byte("a"),
		}, {
			Amount:  a.Bytes(),
			ToAddr:  []byte(in.Address),
			RefTxid: []byte("a"),
		},
	}

	response := &pb.UtxoOutput{
		Header:        newHeader(),
		TotalSelected: need.String(),
		UtxoList:      utxoList,
	}
	return response, nil
}

// PreExecWithSelectUTXO preExec & selectUtxo
func (mcx *MockXClient) PreExecWithSelectUTXO(ctx context.Context, in *pb.PreExecWithSelectUTXORequest, opts ...grpc.CallOption) (*pb.PreExecWithSelectUTXOResponse, error) {

	output, _ := mcx.SelectUTXO(ctx, &pb.UtxoInput{
		Header:    newHeader(),
		Bcname:    in.GetBcname(),
		Address:   in.Address,
		TotalNeed: strconv.Itoa(int(in.GetTotalAmount())),
	})

	cr, _ := mcx.PreExec(ctx, in.Request)
	response := &pb.PreExecWithSelectUTXOResponse{
		Header:     newHeader(),
		Bcname:     in.GetBcname(),
		Response:   cr.GetResponse(),
		UtxoOutput: output,
	}
	return response, nil
}

//  DposCandidates get all candidates of the tdpos consensus
func (mcx *MockXClient) DposCandidates(ctx context.Context, in *pb.DposCandidatesRequest, opts ...grpc.CallOption) (*pb.DposCandidatesResponse, error) {
	return nil, nil
}

//  DposNominateRecords get all records nominated by an user
func (mcx *MockXClient) DposNominateRecords(ctx context.Context, in *pb.DposNominateRecordsRequest, opts ...grpc.CallOption) (*pb.DposNominateRecordsResponse, error) {
	return nil, nil
}

//  DposNomineeRecords get nominated record of a candidate
func (mcx *MockXClient) DposNomineeRecords(ctx context.Context, in *pb.DposNomineeRecordsRequest, opts ...grpc.CallOption) (*pb.DposNomineeRecordsResponse, error) {
	return nil, nil
}

//  DposVoteRecords get all vote records voted by an user
func (mcx *MockXClient) DposVoteRecords(ctx context.Context, in *pb.DposVoteRecordsRequest, opts ...grpc.CallOption) (*pb.DposVoteRecordsResponse, error) {
	return nil, nil
}

//  DposVotedRecords get all vote records of a candidate
func (mcx *MockXClient) DposVotedRecords(ctx context.Context, in *pb.DposVotedRecordsRequest, opts ...grpc.CallOption) (*pb.DposVotedRecordsResponse, error) {
	return nil, nil
}

//  DposCheckResults get check results of a specific term
func (mcx *MockXClient) DposCheckResults(ctx context.Context, in *pb.DposCheckResultsRequest, opts ...grpc.CallOption) (*pb.DposCheckResultsResponse, error) {
	return nil, nil
}

// DposStatus get dpos status
func (mcx *MockXClient) DposStatus(ctx context.Context, in *pb.DposStatusRequest, opts ...grpc.CallOption) (*pb.DposStatusResponse, error) {
	return nil, nil
}

// GetAccountByAK get account sets contain a specific address
func (mcx *MockXClient) GetAccountByAK(ctx context.Context, in *pb.AK2AccountRequest, opts ...grpc.CallOption) (*pb.AK2AccountResponse, error) {
	return &pb.AK2AccountResponse{
		Header:  newHeader(),
		Bcname:  in.GetBcname(),
		Account: []string{"XC1111@xuper"},
	}, nil
}

// GetAddressContracts get contracts of accounts contain a specific address
func (mcx *MockXClient) GetAddressContracts(ctx context.Context, in *pb.AddressContractsRequest, opts ...grpc.CallOption) (*pb.AddressContractsResponse, error) {
	return &pb.AddressContractsResponse{
		Header: newHeader(),
	}, nil
}

//预执行合约
func (mcx *MockXClient) PreExec(ctx context.Context, in *pb.InvokeRPCRequest, opts ...grpc.CallOption) (*pb.InvokeRPCResponse, error) {
	cr := &pb.ContractResponse{
		Status:  200,
		Message: "a",
		Body:    []byte("ok"),
	}
	crb, _ := json.Marshal(cr)
	ir := &pb.InvokeResponse{
		GasUsed:  10,
		Requests: in.GetRequests(),
		Inputs:   []*pb.TxInputExt{{Bucket: "a", Key: []byte("a"), RefTxid: []byte("b"), RefOffset: 0}},
		Outputs:  []*pb.TxOutputExt{{Bucket: "a", Key: []byte("a"), Value: []byte("b")}},
		Response: [][]byte{crb, crb},
	}

	response := &pb.InvokeRPCResponse{
		Header:   newHeader(),
		Bcname:   in.GetBcname(),
		Response: ir,
	}
	return response, nil
}

func newHeader() *pb.Header {
	//XChainErrorEnum_SUCCESS
	return &pb.Header{
		Error: pb.XChainErrorEnum_SUCCESS,
	}
}

type MockEndorserClient interface {
	pb.XendorserClient
}

type MockEClient struct{}

func (mec *MockEClient) EndorserCall(ctx context.Context, in *pb.EndorserRequest, opts ...grpc.CallOption) (*pb.EndorserResponse, error) {
	peur := new(pb.PreExecWithSelectUTXORequest)
	json.Unmarshal(in.GetRequestData(), peur)
	ta := strconv.Itoa(int(peur.GetTotalAmount()))

	data := []byte{}
	if in.RequestName == "PreExecWithFee" {
		mxc := &MockXClient{}
		mxc.SelectUTXO(ctx, &pb.UtxoInput{
			Header:    in.GetHeader(),
			Bcname:    "xuper",
			Address:   peur.Address,
			TotalNeed: ta,
		})
		cr := &pb.ContractResponse{
			Status:  200,
			Message: "a",
			Body:    []byte("ok"),
		}
		crb, _ := json.Marshal(cr)
		ir := &pb.InvokeResponse{
			GasUsed:  10,
			Inputs:   []*pb.TxInputExt{{Bucket: "a", Key: []byte("a"), RefTxid: []byte("b"), RefOffset: 0}},
			Outputs:  []*pb.TxOutputExt{{Bucket: "a", Key: []byte("a"), Value: []byte("b")}},
			Response: [][]byte{crb},
			Requests: peur.GetRequest().GetRequests(),
		}
		per := &pb.PreExecWithSelectUTXOResponse{
			Header:   in.GetHeader(),
			Bcname:   in.GetBcName(),
			Response: ir,
			UtxoOutput: &pb.UtxoOutput{
				UtxoList: []*pb.Utxo{
					{
						ToAddr: []byte(peur.Address),
						Amount: big.NewInt(100).Bytes(),
					},
				},
				TotalSelected: "1000",
			},
		}
		data, _ = json.Marshal(per)
	}
	resp := &pb.EndorserResponse{
		Header:       in.Header,
		ResponseName: in.GetRequestName(),
		EndorserSign: &pb.SignatureInfo{PublicKey: "pubkey", Sign: []byte("endirserSign")},
		ResponseData: data,
	}
	return resp, nil
}

type MockEventClient interface {
	pb.EventServiceClient
}

type MockESClient struct {
	grpc.ClientStream
}

type eventServiceSubscribeClient struct {
	grpc.ClientStream
}

func (x *eventServiceSubscribeClient) Recv() (*pb.Event, error) {
	m := new(pb.Event)
	// if err := x.ClientStream.RecvMsg(m); err != nil {
	// 	return nil, err
	// }
	time.Sleep(time.Millisecond * 100)
	return m, nil
}

func (mesc *MockESClient) Subscribe(ctx context.Context, in *pb.SubscribeRequest, opts ...grpc.CallOption) (pb.EventService_SubscribeClient, error) {
	return &eventServiceSubscribeClient{}, nil
}
