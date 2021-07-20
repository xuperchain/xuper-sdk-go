package xuper

import (
	"encoding/json"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"

	"github.com/xuperchain/xuper-sdk-go/v2/account"
	"github.com/xuperchain/xuper-sdk-go/v2/common"
	"github.com/xuperchain/xuperchain/service/pb"
)

// Request xuperchain transaction request.
type Request struct {
	initiatorAccount *account.Account

	// contract parameters, kernel or user contract.
	module       string
	contractName string
	methodName   string

	// 重要：这个 args 是 pb.InvokeRpcRequest 里面的 args，xuperclient 和 request 参数中的 args 为合约调用的 args。
	args map[string][]byte

	// transfer parameters.
	transferTo     string
	transferAmount string

	opt *requestOptions
}

const (
	// NativeContractModule native contract module.
	NativeContractModule = "native"
	// WasmContractModule wasm contract module.
	WasmContractModule = "wasm"
	// EvmContractModule evm contract module.
	EvmContractModule = "evm"

	// GoRuntime go contract runtime.
	GoRuntime = "go"
	// CRuntime c++ contract runtime.
	CRuntime = "c"
	// JavaRuntime java contract runtime.
	JavaRuntime = "java"

	// EvmJSONEncoded evm contract invoke abi encoded.
	EvmJSONEncoded = "jsonEncoded"
	// EvmJSONEncodedTrue evm contract invoke abi encoded.
	EvmJSONEncodedTrue = "true"

	// XkernelModule xkernel contract module
	XkernelModule = "kernel"
	// Xkernel3Module xkernel contract module
	Xkernel3Module = "xkernel"
	// XkernelDeployMethod xkernel contract deploy contract method.
	XkernelDeployMethod = "Deploy"
	// XkernelUpgradeMethod xkernel contract upgrade contract method.
	XkernelUpgradeMethod = "Upgrade"
	// XkernelNewAccountMethod xkernel contract create contract account method.
	XkernelNewAccountMethod = "NewAccount"
	// XkernelSetAccountACLMethod xkernel contract set account ACL method.
	XkernelSetAccountACLMethod = "SetAccountAcl"
	// XkernelSetMethodACLMethod xkernel contract set method ACL method.
	XkernelSetMethodACLMethod = "SetMethodAcl"

	// ArgAccountName account name field.
	ArgAccountName = "account_name"
	// ArgContractName contract name field.
	ArgContractName = "contract_name"
	// ArgContractCode contract code field.
	ArgContractCode = "contract_code"
	// ArgContractDesc contract desc field.
	ArgContractDesc = "contract_desc"
	// ArgInitArgs contract init args field.
	ArgInitArgs = "init_args"
	// ArgContractAbi evm abi field.
	ArgContractAbi = "contract_abi"
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
	args map[string][]byte,
	transferTo, transferAmount string,
	opts ...RequestOption,
) (*Request, error) {

	if initiator == nil {
		return nil, errors.New("initiator can not be nil")
	}

	opt, err := initOpts(opts...)
	if err != nil {
		return nil, err
	}

	if opt.onlyFeeFromAccount && !initiator.HasContractAccount() {
		return nil, errors.Wrap(common.ErrInvalidAccount,
			"initiator contract account can not be nil when set fee from account.")
	}

	return &Request{
		initiatorAccount: initiator,
		module:           module,
		contractName:     contractName,
		methodName:       methodName,
		args:             args,
		transferTo:       transferTo,
		transferAmount:   transferAmount,
		opt:              opt,
	}, nil
}

// SetInitiatorAccount set request initiator.
func (r *Request) SetInitiatorAccount(account *account.Account) error {
	r.initiatorAccount = account
	return nil
}

// SetArgs set request args. NOTE: this is pb.InvokeRPCRequest args, not contract invoke args.
func (r *Request) SetArgs(args map[string][]byte) error {
	r.args = args
	return nil
}

// SetModule set request contract module.
func (r *Request) SetModule(module string) error {
	r.module = module
	return nil
}

// SetContractName set
func (r *Request) SetContractName(contractName string) error {
	r.contractName = contractName
	return nil
}

// SetTransferTo set
func (r *Request) SetTransferTo(to string) error {
	r.transferTo = to
	return nil
}

// SetTransferAmount set
func (r *Request) SetTransferAmount(amount string) error {
	r.transferAmount = amount
	return nil
}

// NewTransferRequest set
func NewTransferRequest(from *account.Account, to, amount string, opts ...RequestOption) (*Request, error) {
	if from == nil {
		return nil, common.ErrInvalidInitiator
	}

	if to == "" {
		return nil, common.ErrInvalidParam
	}

	amount, ok := common.IsValidAmount(amount)
	if !ok {
		return nil, common.ErrInvalidAmount
	}

	return NewRequest(from, "", "", "", nil, to, amount, opts...)
}

// NewDeployContractRequest new request for deploy contract, wasm, evm and native.
func NewDeployContractRequest(from *account.Account, name string, abi, code []byte, args map[string]string, contractType, runtime string, opts ...RequestOption) (*Request, error) {
	if from == nil || !from.HasContractAccount() {
		return nil, common.ErrInvalidAccount
	}

	if name == "" || contractType == "" || len(code) == 0 {
		return nil, common.ErrInvalidParam
	}

	reqArgs := generateDeployArgs(args, abi, code, contractType, runtime, from.GetContractAccount(), name)

	return NewRequest(from, Xkernel3Module, "", XkernelDeployMethod, reqArgs, "", "", opts...)
}

// NewInvokeContractRequest new request for invoke contract, wasm, evm and native.
func NewInvokeContractRequest(from *account.Account, module, name, method string, args map[string]string, opts ...RequestOption) (*Request, error) {
	if from == nil {
		return nil, errors.New("invalid initiator")
	}

	if module == "" && name == "" && method == "" {
		return nil, common.ErrInvalidParam
	}

	reqArgs, err := generateInvokeArgs(args, module)
	if err != nil {
		return nil, err
	}

	return NewRequest(from, module, name, method, reqArgs, "", "", opts...)
}

// NewUpgradeContractRequest new upgrade contract request. NOTE: evm contract upgrade disabled!
func NewUpgradeContractRequest(from *account.Account, module, name string, code []byte, opts ...RequestOption) (*Request, error) {
	if from == nil || !from.HasContractAccount() {
		return nil, common.ErrInvalidAccount
	}

	if module == "" || name == "" || len(code) == 0 {
		return nil, common.ErrInvalidParam
	}

	reqArgs := generateDeployArgs(nil, nil, code, module, "", from.GetContractAccount(), name)
	return NewRequest(from, Xkernel3Module, "", XkernelUpgradeMethod, reqArgs, "", "", opts...)
}

// NewCreateContractAccountRequest new request for create contract account.
func NewCreateContractAccountRequest(from *account.Account, contractAccount string, opts ...RequestOption) (*Request, error) {
	if from == nil || from.HasContractAccount() {
		return nil, common.ErrInvalidAccount
	}

	if contractAccount == "" {
		return nil, common.ErrInvalidAccount
	}

	args, err := genAccountACLArgs(getDefaultACL(from.Address), contractAccount)
	if err != nil {
		return nil, err
	}
	return NewRequest(from, Xkernel3Module, "", XkernelNewAccountMethod, args, "", "", opts...)
}

// NewSetMethodACLRequest new request for set method ACL.
func NewSetMethodACLRequest(from *account.Account, name, method string, acl *ACL, opts ...RequestOption) (*Request, error) {
	if from == nil {
		return nil, common.ErrInvalidAccount
	}

	if acl == nil {
		return nil, errors.New("invalid ACL")
	}

	if method == "" || name == "" {
		return nil, common.ErrInvalidParam
	}

	args, err := genMethodACLArgs(acl, name, method)
	if err != nil {
		return nil, err
	}
	return NewRequest(from, Xkernel3Module, "", XkernelSetMethodACLMethod, args, "", "", opts...)
}

// NewSetAccountACLRequest new request for set contract account acl.
func NewSetAccountACLRequest(from *account.Account, acl *ACL, opts ...RequestOption) (*Request, error) {
	if from == nil || !from.HasContractAccount() {
		return nil, common.ErrInvalidAccount
	}

	args, err := genAccountACLArgs(acl, from.GetContractAccount())
	if err != nil {
		return nil, err
	}
	return NewRequest(from, Xkernel3Module, "", XkernelSetAccountACLMethod, args, "", "", opts...)
}

func generateDeployArgs(arg map[string]string, abi, code []byte, module, runtime, contractAccount, contractName string) map[string][]byte {
	argstmp := map[string][]byte{}
	if module == EvmContractModule {
		argsTmp := make(map[string]interface{}, len(arg))
		for k, v := range arg {
			argsTmp[k] = v
		}
		argstmp, _ = convertToXuper3EvmArgs(argsTmp)
	} else {
		argstmp = convertToXuperContractArgs(arg)
	}

	initArgs, _ := json.Marshal(argstmp)

	desc := &pb.WasmCodeDesc{
		ContractType: module,
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

	if module == EvmContractModule {
		args["contract_abi"] = abi
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

func generateInvokeArgs(arg map[string]string, module string) (map[string][]byte, error) {
	if module == EvmContractModule {
		// todo
		argsTmp := make(map[string]interface{}, len(arg))
		for k, v := range arg {
			argsTmp[k] = v
		}
		return convertToXuper3EvmArgs(argsTmp)

	}
	return convertToXuperContractArgs(arg), nil
}

func genAccountACLArgs(acl *ACL, contractAccount string) (map[string][]byte, error) {
	ACLBytes, err := json.Marshal(acl)
	if err != nil {
		return nil, errors.New("invalid ACL")
	}

	args := map[string][]byte{
		"account_name": []byte(contractAccount),
		"acl":          ACLBytes,
	}
	return args, nil
}

func genMethodACLArgs(acl *ACL, name, method string) (map[string][]byte, error) {
	ACLStr, err := json.Marshal(acl)
	if err != nil {
		return nil, errors.New("invalid ACL")
	}

	args := map[string][]byte{
		"contract_name": []byte(name),
		"method_name":   []byte(method),
		"acl":           []byte(ACLStr),
	}
	return args, nil
}
