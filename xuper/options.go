package xuper

import "github.com/pkg/errors"

type clientOptions struct {
	configFile  string
	useGrpcGZIP bool
	grpcTLS     *grpcTLSConfig
}

type grpcTLSConfig struct {
	serverName string
	cacertFile string
	certFile   string
	keyFile    string
}

type requestOptions struct {
	onlyFeeFromAccount   bool
	fee                  string
	bcname               string
	contractInvokeAmount string
	desc                 string
	otherAuthRequire     []string
	notPost              bool
	gasFromAddress       bool
}

type queryOption struct {
	bcname string
	blockNums int64
	waitTime int64
}

// RequestOption tx opt.
type RequestOption func(opt *requestOptions) error

// ClientOption xuperclient opt.
type ClientOption func(opt *clientOptions) error

// QueryOption query opt.
type QueryOption func(opt *queryOption) error

// WithQueryBcname query method bcname option.
func WithQueryBcname(bcname string) QueryOption {
	return func(opts *queryOption) error {
		opts.bcname = bcname
		return nil
	}
}

// WithBlockNums query tx status
func WithBlockNums(blockNums int64) QueryOption  {
	return func(opts *queryOption) error {
		if 0 < blockNums && blockNums < 100 {
			opts.blockNums = blockNums
			return nil
		}
		return errors.New("block nums must be between 0 and 100")
	}
}

// WithMaxWaitTime query tx status max wait time
func WithMaxWaitTime(waitTime int64) QueryOption {
	return func(opts *queryOption) error {
		opts.waitTime = waitTime
		return nil
	}
}


// WithConfigFile set xuperclient config file.
func WithConfigFile(configFile string) ClientOption {
	return func(opts *clientOptions) error {
		opts.configFile = configFile
		return nil
	}
}

// WithGrpcGZIP use gzip.
func WithGrpcGZIP() ClientOption {
	return func(opts *clientOptions) error {
		opts.useGrpcGZIP = true
		return nil
	}
}

// WithGrpcTLS grpc TLS cert config.
func WithGrpcTLS(serverName, cacertFile, certFile, keyFile string) ClientOption {
	return func(opts *clientOptions) error {
		if opts.grpcTLS == nil {
			opts.grpcTLS = new(grpcTLSConfig)
		}
		opts.grpcTLS.serverName = serverName
		opts.grpcTLS.cacertFile = cacertFile
		opts.grpcTLS.certFile = certFile
		opts.grpcTLS.keyFile = keyFile
		return nil
	}
}

// WithFeeFromAccount fee & gas from contract account.
func WithFeeFromAccount() RequestOption {
	return func(opts *requestOptions) error {
		opts.onlyFeeFromAccount = true
		return nil
	}
}

func WithGasFromAddress() RequestOption {
	return func(opts *requestOptions) error {
		opts.gasFromAddress = true
		return nil
	}
}

// WithFee set fee.
func WithFee(fee string) RequestOption {
	return func(opts *requestOptions) error {
		opts.fee = fee
		return nil
	}
}

// WithBcname set blockchain name.
func WithBcname(name string) RequestOption {
	return func(opts *requestOptions) error {
		if name == "" {
			return errors.New("invalid bcname")
		}
		opts.bcname = name
		return nil
	}
}

// WithContractInvokeAmount set transfer to contract when invoke contract.
func WithContractInvokeAmount(amount string) RequestOption {
	return func(opts *requestOptions) error {
		opts.contractInvokeAmount = amount
		return nil
	}
}

// WithDesc set tx desc.
func WithDesc(desc string) RequestOption {
	return func(opts *requestOptions) error {
		opts.desc = desc
		return nil
	}
}

// WithNotPost generate transaction only, won't post to server.
func WithNotPost() RequestOption {
	return func(opts *requestOptions) error {
		opts.notPost = true
		return nil
	}
}

// WithOtherAuthRequires for multisign, other address need sign, exclude initiator.
func WithOtherAuthRequires(authRequires []string) RequestOption {
	return func(opts *requestOptions) error {
		opts.otherAuthRequire = authRequires
		return nil
	}
}


