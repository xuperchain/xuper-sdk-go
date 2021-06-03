package xuper

import "github.com/xuperchain/xuper-sdk-go/v2/account"

type Request struct {
	initiatorAccount *account.Account

	// contract parameters.
	module       string
	runtime      string
	code         []byte
	contractName string
	args         map[string]string

	// transfer parameters.
	transferTo     string
	transferAmount string

	opt *requestOptions
}

type requestOptions struct {
	feeFromAccount       bool
	fee                  string
	bcname               string
	contractInvokeAmount string
	desc                 string
	notPost              bool
}

type RequestOption func(opt *requestOptions) error

func NewDeployContractRequest(from *account.Account, name string, code []byte, args map[string]string, opts ...RequestOption) (*Request, error) {

	return nil, nil
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

func (r *Request) SetArgs(args map[string]string) error {
	r.args = args
	return nil
}

func (r *Request) SetModule(module string) error {
	r.module = module
	return nil
}

func (r *Request) SetRuntime(runtime string) error {
	r.runtime = runtime
	return nil
}

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
