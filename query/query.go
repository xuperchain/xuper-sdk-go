package query

import (
	"github.com/xuperchain/xuper-sdk-go/account"
	"github.com/xuperchain/xuper-sdk-go/config"
	"github.com/xuperchain/xuper-sdk-go/pb"
	"github.com/xuperchain/xuper-sdk-go/xchain"
)

type QueryClient struct {
	xchain.Xchain
}

func InitClient(account *account.Account, bcName string, sdkClient *xchain.SDKClient) *QueryClient {
	commConfig := config.GetInstance()
	return &QueryClient{
		Xchain: xchain.Xchain{
			Cfg:       commConfig,
			Account:   account,
			ChainName: bcName,
			SDKClient: sdkClient,
		},
	}
}

func (qc *QueryClient) QueryBlockByHeight(height int64) (*pb.Block, error) {
	return qc.Xchain.QueryBlockByHeight(height)
}

func (qc *QueryClient) GetAccountByAk(address string) (*pb.AK2AccountResponse, error) {
	return qc.Xchain.GetAccountByAk(address)
}

//
//
func (qc *QueryClient) GetAccountContracts(address string) (*pb.GetAccountContractsResponse, error) {
	return qc.Xchain.GetAccountContracts(address)
}

//
func (qc *QueryClient) QueryUTXORecord(addr string, utxoItemNum int64) (*pb.UtxoRecordDetail, error) {
	return qc.Xchain.QueryUTXORecord(addr, utxoItemNum)
}

//
func (qc *QueryClient) QueryContractMethondAcl(contract string, method string) (*pb.AclStatus, error) {
	return qc.Xchain.QueryContractMethondAcl(contract, method)

}
