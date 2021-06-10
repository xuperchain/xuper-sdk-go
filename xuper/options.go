package xuper

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

type txOptions struct {
	feeFromAccount       bool
	fee                  string
	bcname               string
	contractInvokeAmount string
	desc                 string
	notPost              bool
}

// TxOption tx
type TxOption func(opt *txOptions) error

// ClientOption func
type ClientOption func(opt *clientOptions) error

// func WithFeeFromAccount() Option {
// 	return func(opts *options) error {
// 		opts.feeFromAccount = true
// 		return nil
// 	}
// }

// func WithFee(fee string) Option {
// 	return func(opts *options) error {
// 		// todo check fee valid
// 		opts.fee = fee
// 		return nil
// 	}
// }

// func WithBcname(name string) Option {
// 	return func(opts *options) error {
// 		if name == "" {
// 			return errors.New("invalid bcname")
// 		}
// 		opts.bcname = name
// 		return nil
// 	}
// }

// func WithContractInvokeAmount(amount string) Option {
// 	return func(opts *options) error {
// 		opts.contractInvokeAmount = amount
// 		return nil
// 	}
// }

// func WithDesc(desc string) Option {
// 	return func(opts *options) error {
// 		opts.desc = desc
// 		return nil
// 	}
// }
