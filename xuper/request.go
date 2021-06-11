package xuper

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuperchain/core/pb"
)

type Request struct {
	initiatorAccount *account.Account

	// contract parameters, kernel or user contract.
	module       string
	code         []byte
	contractName string
	methodName   string
	args         map[string][]byte

	// transfer parameters.
	transferTo     string
	transferAmount string

	opt *requestOptions
}

const (
	NativeContractType = "naitve"
	WasmContractType   = "wasm"
	EvmContractType    = "evm"

	GoRuntime   = "go"
	CRuntime    = "c"
	JavaRuntime = "java"

	EvmJSONEncoded     = "jsonEncoded"
	EvmJSONEncodedTrue = "true"

	XkernelModule           = "xkernel"
	XkernelDeployMethod     = "Deploy"
	XkernelUpgradeMethod    = "Upgrade"
	XkernelNewAccountMethod = "NewAccount"

	ArgAccountName  = "account_name"
	ArgContractName = "contract_name"
	ArgContractCode = "contract_code"
	ArgContractDesc = "contract_desc"
	ArgInitArgs     = "init_args"
	ArgContractAbi  = "contract_abi"
)

func initOpts(opts ...RequestOption) (*requestOptions, error) {
	opt := &requestOptions{}
	for _, param := range opts {
		err := param(opt)
		if err != nil {
			return nil, fmt.Errorf("option failed: %v", err)
		}
	}
	return opt, nil
}

// NewRequest new custom request.
func NewRequest(
	initiator *account.Account,
	module, contractName, methodName string,
	code []byte,
	args map[string][]byte,
	transferTo, transferAmount string,
	opts ...RequestOption,
) (*Request, error) {

	if initiator == nil {
		return nil, errors.New("initiator can not be nil")
	}

	return nil, nil
}

// NewDeployContractRequest new request for deploy contract, wasm, evm and native.
func NewDeployContractRequest(from *account.Account, name string, code []byte, args map[string]string, contractType, runtime string, opts ...RequestOption) (*Request, error) {
	if from == nil || !from.HasContractAccount() {
		return nil, common.ErrInvalidAccount
	}

	reqArgs := generateDeployArgs(args, code, contractType, runtime, from.GetContractAccount(), name)

	return NewRequest(from, XkernelModule, "", XkernelDeployMethod, code, reqArgs, "", "", opts...)
}

func makeRequestArgs() {

}

// 参考
func generateDeployArgs(arg map[string]string, code []byte, contractType, runtime, contractAccount, contractName string) map[string][]byte {
	argstmp := convertToXuperContractArgs(arg)
	initArgs, _ := json.Marshal(argstmp)

	desc := &pb.WasmCodeDesc{
		ContractType: contractType,
		Runtime:      runtime,
	}
	contractDesc, _ := proto.Marshal(desc)

	args := map[string][]byte{
		"account_name":  []byte(contractAccount),
		"contract_name": []byte(contractName),
		"contract_code": code,
		"contract_desc": contractDesc,
		"init_args":     initArgs,
	}
	return args
}

func convertToXuperContractArgs(args map[string]string) map[string][]byte {
	argmap := make(map[string][]byte)
	for k, v := range args {
		argmap[k] = []byte(v)
	}
	return argmap
}

func convertToXuper3EvmArgs(args map[string]interface{}) (map[string][]byte, error) {
	input, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	// 此处与 server 端结构相同，如果 jsonEncoded 字段修改，server 端也要修改（core/contract/evm/creator.go）。
	ret := map[string][]byte{
		"input":        input,
		EvmJSONEncoded: []byte(EvmJSONEncodedTrue),
	}
	return ret, nil
}

func NewInvokeContractRequest(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Request, error) {

	return nil, nil
}

func NewTransferRequest(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Request, error) {

	return nil, nil
}

func (r *Request) SetInitiatorAccount(account *account.Account) error {
	r.initiatorAccount = account
	return nil
}

func (r *Request) SetArgs(args map[string][]byte) error {
	r.args = args
	return nil
}

func (r *Request) SetModule(module string) error {
	r.module = module
	return nil
}

// func (r *Request) SetRuntime(runtime string) error {
// 	// r.runtime = runtime
// 	return nil
// }

func (r *Request) SetContractName(contractName string) error {
	r.contractName = contractName
	return nil
}

func (r *Request) SetCode(code []byte) error {
	r.code = code
	return nil
}

func (r *Request) SetTransferTo(to string) error {
	r.transferTo = to
	return nil
}
func (r *Request) SetTransferAmount(amount string) error {
	r.transferAmount = amount
	return nil
}
