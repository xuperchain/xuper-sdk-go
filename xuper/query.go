package xuper

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuperchain/core/pb"
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

	ctx := context.Background()
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	txStatus := &pb.TxStatus{
		Bcname: getBCname(opt),
		Txid:   rawTx,
	}
	res, err := x.xc.QueryTx(ctx, txStatus)
	if err != nil {
		return nil, err
	}
	if res.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(res.Header.Error.String())
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
	ctx := context.Background()
	block, err := x.xc.GetBlock(ctx, blockIDPB)
	if err != nil {
		return nil, err
	}
	if block.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(block.Header.Error.String())
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
	ctx := context.Background()
	block, err := x.xc.GetBlockByHeight(ctx, blockHeightPB)
	if err != nil {
		return nil, err
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
	aclStatus, err := x.xc.QueryACL(context.Background(), in)
	if err != nil {
		return nil, err
	}

	acl := &ACL{}
	pm := PermissionModel{}
	pm.Rule = int32(aclStatus.Acl.Pm.Rule) // 类型别名转换
	pm.AcceptValue = aclStatus.Acl.Pm.AcceptValue

	acl.PM = pm
	acl.AksWeight = aclStatus.Acl.AksWeight
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

	fmt.Printf("in:%+v\n", in)
	aclStatus, err := x.xc.QueryACL(context.Background(), in)
	if err != nil {
		return nil, err
	}

	if aclStatus == nil {
		return nil, nil
	}

	acl := &ACL{}
	pm := PermissionModel{}
	pm.Rule = int32(aclStatus.Acl.Pm.Rule) //类型别名转换
	pm.AcceptValue = aclStatus.Acl.Pm.AcceptValue

	acl.PM = pm
	acl.AksWeight = aclStatus.Acl.AksWeight
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

	ctx := context.Background()
	res, err := x.xc.GetAccountContracts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.GetContractsStatus(), nil
}

func (x *XClient) queryAddressContracts(address string, opts ...QueryOption) (map[string]*pb.ContractList, error) { //todo  return 修改
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	req := &pb.AddressContractsRequest{
		Address: address,
		Bcname:  getBCname(opt),
	}

	ctx := context.Background()
	res, err := x.xc.GetAddressContracts(ctx, req)
	if err != nil {
		return nil, err
	}

	return res.Contracts, nil
}

func (x *XClient) queryBalance(address string, opts ...QueryOption) (*pb.AddressStatus, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	addrstatus := &pb.AddressStatus{
		Address: address,
		Bcs: []*pb.TokenDetail{
			{Bcname: getBCname(opt)},
		},
	}

	ctx := context.Background()
	reply, err := x.xc.GetBalance(ctx, addrstatus)
	if err != nil {
		return nil, err
	}

	if reply.Header.Error != pb.XChainErrorEnum_SUCCESS {
		return nil, errors.New(reply.Header.Error.String())
	}
	return reply, nil
}

func (x *XClient) queryBalanceDetail(address string, opts ...QueryOption) (*pb.AddressBalanceStatus, error) {
	ctx := context.Background()
	addressBalanceStatus := &pb.AddressBalanceStatus{
		Address: address,
	}

	return x.xc.GetBalanceDetail(ctx, addressBalanceStatus)
}

func (x *XClient) querySystemStatus(opts ...QueryOption) (*pb.SystemsStatusReply, error) {
	req := &pb.CommonIn{}
	ctx := context.Background()
	return x.xc.GetSystemStatus(ctx, req)
}

func (x *XClient) queryBlockChains(opts ...QueryOption) (*pb.BCStatus, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	bcStatusPB := &pb.BCStatus{
		Bcname: getBCname(opt),
	}

	return x.xc.GetBlockChainStatus(context.TODO(), bcStatusPB)
}

func (x *XClient) queryBlockChainStatus(chainName string, opts ...QueryOption) (*pb.BCStatus, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	bcStatusPB := &pb.BCStatus{
		Bcname: getBCname(opt),
	}

	return x.xc.GetBlockChainStatus(context.TODO(), bcStatusPB)
}

func (x *XClient) queryNetURL(opts ...QueryOption) (*pb.RawUrl, error) {
	req := &pb.CommonIn{}
	ctx := context.Background()
	return x.xc.GetNetURL(ctx, req)
}

func (x *XClient) queryAccountByAK(address string, opts ...QueryOption) (*pb.AK2AccountResponse, error) {
	opt, err := initQueryOpts(opts...)
	if err != nil {
		return nil, err
	}

	AK2AccountRequest := &pb.AK2AccountRequest{
		Bcname:  getBCname(opt),
		Address: address,
	}

	return x.xc.GetAccountByAK(context.Background(), AK2AccountRequest)
}
