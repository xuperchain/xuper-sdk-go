package xuper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"

	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuperchain/service/pb"
)

func initQueryOpts(opts ...QueryOption) (*queryOption, error) {
	opt := &queryOption{}
	for _, param := range opts {
		err := param(opt)
		if err != nil {
			return nil, fmt.Errorf("option failed: %v", err)
		}
	}

	return opt, nil
}

func getBCname(opt *queryOption) string {
	chainName := defaultChainName
	if opt.bcname != "" {
		chainName = opt.bcname
	}
	return chainName
}

func (x *XClient) queryTxByID(txID string, opts ...QueryOption) (*pb.Transaction, error) {
	rawTx, err := hex.DecodeString(txID)
	if err != nil {
		return nil, err
	}

	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	txStatus := &pb.TxStatus{
		Bcname: getBCname(opt),
		Txid:   rawTx,
	}
	res, err := x.xc.QueryTx(context.TODO(), txStatus)
	if err != nil {
		return nil, err
	}
	if res.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(res.GetHeader().GetError().String())
	}
	if res.Tx == nil {
		return nil, common.ErrTxNotFound
	}
	return res.Tx, nil
}

func (x *XClient) queryBlockByID(blockID string, opts ...QueryOption) (*pb.Block, error) {
	rawBlockid, err := hex.DecodeString(blockID)
	if err != nil {
		return nil, err
	}

	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	blockIDPB := &pb.BlockID{
		Bcname:      getBCname(opt),
		Blockid:     rawBlockid,
		NeedContent: true,
	}

	block, err := x.xc.GetBlock(context.TODO(), blockIDPB)
	if err != nil {
		return nil, err
	}
	if block.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(block.GetHeader().GetError().String())
	}
	if block.Block == nil {
		return nil, errors.New("block not found")
	}

	return block, nil
}

func (x *XClient) queryBlockByHeight(height int64, opts ...QueryOption) (*pb.Block, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	blockHeightPB := &pb.BlockHeight{
		Bcname: getBCname(opt),
		Height: height,
	}

	block, err := x.xc.GetBlockByHeight(context.TODO(), blockHeightPB)
	if err != nil {
		return nil, err
	}

	if block.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(block.GetHeader().GetError().String())
	}
	if block.Block == nil {
		return nil, errors.New("block not found")
	}

	return block, nil
}

func (x *XClient) queryAccountACL(account string, opts ...QueryOption) (*ACL, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	in := &pb.AclStatus{
		Bcname:      getBCname(opt),
		AccountName: account,
	}
	aclStatus, err := x.xc.QueryACL(context.TODO(), in)
	if err != nil {
		return nil, err
	}

	if aclStatus.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(aclStatus.GetHeader().GetError().String())
	}

	acl := &ACL{}
	pm := PermissionModel{}
	pm.Rule = int32(aclStatus.GetAcl().GetPm().GetRule())
	pm.AcceptValue = aclStatus.GetAcl().GetPm().GetAcceptValue()

	acl.PM = pm
	acl.AksWeight = aclStatus.GetAcl().GetAksWeight()
	return acl, nil

}

func (x *XClient) queryMethodACL(name, method string, opts ...QueryOption) (*ACL, error) { // todo
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	in := &pb.AclStatus{
		Bcname:       getBCname(opt),
		ContractName: name,
		MethodName:   method,
	}

	aclStatus, err := x.xc.QueryACL(context.TODO(), in)
	if err != nil {
		return nil, err
	}

	if aclStatus.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(aclStatus.GetHeader().GetError().String())
	}

	if aclStatus == nil {
		return nil, nil
	}

	acl := &ACL{}
	pm := PermissionModel{}
	pm.Rule = int32(aclStatus.GetAcl().GetPm().GetRule())
	pm.AcceptValue = aclStatus.GetAcl().GetPm().GetAcceptValue()

	acl.PM = pm
	acl.AksWeight = aclStatus.GetAcl().GetAksWeight()
	return acl, nil
}

func (x *XClient) queryAccountContracts(account string, opts ...QueryOption) ([]*pb.ContractStatus, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	req := &pb.GetAccountContractsRequest{
		Bcname:  getBCname(opt),
		Account: account,
	}

	resp, err := x.xc.GetAccountContracts(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	if resp.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(resp.GetHeader().GetError().String())
	}

	return resp.GetContractsStatus(), nil
}

func (x *XClient) queryAddressContracts(address string, opts ...QueryOption) (map[string]*pb.ContractList, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	req := &pb.AddressContractsRequest{
		Address: address,
		Bcname:  getBCname(opt),
	}

	resp, err := x.xc.GetAddressContracts(context.TODO(), req)
	if err != nil {
		return nil, err
	}

	if resp.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(resp.GetHeader().GetError().String())
	}

	return resp.GetContracts(), nil
}

func (x *XClient) queryBalance(address string, opts ...QueryOption) (*big.Int, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	bcname := getBCname(opt)
	addrstatus := &pb.AddressStatus{
		Address: address,
		Bcs: []*pb.TokenDetail{
			{Bcname: bcname},
		},
	}

	reply, err := x.xc.GetBalance(context.TODO(), addrstatus)
	if err != nil {
		return nil, err
	}

	if reply.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(reply.GetHeader().GetError().String())
	}

	for _, v := range reply.Bcs {
		if v.GetBcname() == bcname {
			if v.GetError() != pb.XChainErrorEnum_SUCCESS {
				return nil, errors.New(v.GetError().String())
			}

			if v.GetBalance() == "" {
				return big.NewInt(0), nil
			}
			bal, ok := big.NewInt(0).SetString(v.GetBalance(), 10)
			if !ok {
				return nil, errors.New("invalid balance query from chain")
			}
			return bal, nil
		}
	}

	return nil, errors.New("invalid bcname:" + bcname)
}

// BalanceDetail address or account balance detailds.
type BalanceDetail struct {
	Balance  string
	IsFrozen bool
}

func (x *XClient) queryBalanceDetail(address string, opts ...QueryOption) ([]*BalanceDetail, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}
	tfds := []*pb.TokenFrozenDetails{{Bcname: getBCname(opt)}}
	ctx := context.Background()
	addressBalanceStatus := &pb.AddressBalanceStatus{
		Address: address,
		Tfds:    tfds,
	}

	bs, err := x.xc.GetBalanceDetail(ctx, addressBalanceStatus)
	if err != nil {
		return nil, err
	}

	if bs.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(bs.GetHeader().GetError().String())
	}

	bcname := getBCname(opt)
	for _, tfd := range bs.Tfds {
		if tfd.Bcname == bcname {
			if tfd.GetError() != pb.XChainErrorEnum_SUCCESS {
				return nil, errors.New(bs.GetHeader().GetError().String())
			}

			result := make([]*BalanceDetail, 0, len(tfd.Tfd))
			for _, v := range tfd.Tfd {
				result = append(result, &BalanceDetail{
					Balance:  v.Balance,
					IsFrozen: v.IsFrozen,
				})
			}

			return result, nil
		}
	}

	return nil, fmt.Errorf("Can not query balance detail for bcname: %s", bcname)
}

func (x *XClient) querySystemStatus(opts ...QueryOption) (*pb.SystemsStatusReply, error) {
	ss, err := x.xc.GetSystemStatus(context.TODO(), &pb.CommonIn{})
	if err != nil {
		return nil, err
	}

	if ss.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(ss.GetHeader().GetError().String())
	}
	return ss, nil
}

func (x *XClient) queryBlockChains(opts ...QueryOption) ([]string, error) {
	bcs, err := x.xc.GetBlockChains(context.TODO(), &pb.CommonIn{})
	if err != nil {
		return nil, err
	}

	if bcs.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(bcs.GetHeader().GetError().String())
	}

	return bcs.GetBlockchains(), nil
}

func (x *XClient) queryBlockChainStatus(chainName string, opts ...QueryOption) (*pb.BCStatus, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	bcStatusPB := &pb.BCStatus{
		Bcname: getBCname(opt),
	}

	bcs, err := x.xc.GetBlockChainStatus(context.TODO(), bcStatusPB)
	if err != nil {
		return nil, err
	}

	if bcs.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(bcs.GetHeader().GetError().String())
	}

	return bcs, err
}

func (x *XClient) queryNetURL(opts ...QueryOption) (string, error) {
	rawURL, err := x.xc.GetNetURL(context.TODO(), &pb.CommonIn{})
	if err != nil {
		return "", err
	}

	if rawURL.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return "", errors.New(rawURL.GetHeader().GetError().String())
	}

	return rawURL.GetRawUrl(), nil
}

func (x *XClient) queryAccountByAK(address string, opts ...QueryOption) ([]string, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	AK2AccountRequest := &pb.AK2AccountRequest{
		Bcname:  getBCname(opt),
		Address: address,
	}

	resp, err := x.xc.GetAccountByAK(context.TODO(), AK2AccountRequest)
	if err != nil {
		return nil, err
	}

	if resp.GetHeader().GetError() != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(resp.GetHeader().GetError().String())
	}
	return resp.GetAccount(), nil
}
