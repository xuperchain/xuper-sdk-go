package xuper

import (
	"errors"
	"fmt"

	"github.com/xuperchain/xuperchain/service/pb"
)

// Watcher event watcher.
type Watcher struct {
	FilteredBlockChan <-chan *FilteredBlock
	exit              chan<- struct{}

	opt *blockEventOption
}

func initEventOpts(opts ...BlockEventOption) (*blockEventOption, error) {
	opt := &blockEventOption{
		blockChanBufferSize: 100, // default 100.
		blockFilter: &pb.BlockFilter{
			Bcname: "xuper", // default xuper.
		},
	}

	for _, param := range opts {
		err := param(opt)
		if err != nil {
			return nil, fmt.Errorf("event option failed: %v", err)
		}
	}

	return opt, nil
}

// Close close watcher.
func (w *Watcher) Close() {
	close(w.exit)
}

// FilteredBlock pb.FilteredBlock
type FilteredBlock struct {
	Bcname      string                 `json:"bcname,omitempty"`
	Blockid     string                 `json:"blockid,omitempty"`
	BlockHeight int64                  `json:"block_height,omitempty"`
	Txs         []*FilteredTransaction `json:"txs,omitempty"`
}

// FilteredTransaction pb.FilteredTransaction
type FilteredTransaction struct {
	Txid   string           `json:"txid,omitempty"`
	Events []*ContractEvent `json:"events,omitempty"`
}

// ContractEvent pb.ContractEvent
type ContractEvent struct {
	Contract string `json:"contract,omitempty"`
	Name     string `json:"name,omitempty"`
	Body     string `json:"body,omitempty"`
}

func fromFilteredBlockPB(pbblock *pb.FilteredBlock) *FilteredBlock {
	block := &FilteredBlock{
		Bcname:      pbblock.Bcname,
		Blockid:     pbblock.Blockid,
		BlockHeight: pbblock.BlockHeight,
		Txs:         make([]*FilteredTransaction, 0, len(pbblock.Txs)),
	}

	for _, pbtx := range pbblock.Txs {
		tx := &FilteredTransaction{
			Txid:   pbtx.Txid,
			Events: make([]*ContractEvent, 0, len(pbtx.Events)),
		}
		for _, pbevent := range pbtx.Events {
			tx.Events = append(tx.Events, &ContractEvent{
				Contract: pbevent.Contract,
				Name:     pbevent.Name,
				Body:     string(pbevent.Body),
			})
		}
		block.Txs = append(block.Txs, tx)
	}
	return block
}

// BlockEventOption event opt.
type BlockEventOption func(*blockEventOption) error

type blockEventOption struct {
	blockFilter *pb.BlockFilter

	blockChanBufferSize uint
	skipEmptyTx         bool
}

// WithBlockChanBufferSize block event block channel size, default 100.
func WithBlockChanBufferSize(size uint) BlockEventOption {
	return func(f *blockEventOption) error {
		if size < 0 {
			return errors.New("Invalid size for watcher blockChanBufferSize chan")
		}
		f.blockChanBufferSize = size
		return nil
	}
}

// WithSkipEmplyTx block event skip emply tx block.
func WithSkipEmplyTx() BlockEventOption {
	return func(f *blockEventOption) error {
		f.skipEmptyTx = true
		return nil
	}
}

// WithBlockEventBcname blockchain name.
func WithBlockEventBcname(name string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.Bcname = name
		return nil
	}
}

// WithContract indicates the contract name from which tx are to be received.
func WithContract(contract string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.Contract = contract
		return nil
	}
}

// WithEventName indicates the event name from which events are to be received.
func WithEventName(eventName string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.EventName = eventName
		return nil
	}
}

// WithInitiator indicates the contract initiator from which tx are to be received.
func WithInitiator(initiator string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.Initiator = initiator
		return nil
	}
}

// WithAuthRequire indicates the auth require from which tx are to be received.
func WithAuthRequire(authRequire string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.AuthRequire = authRequire
		return nil
	}
}

// WithFromAddr indicates the transfer address from which tx are to be received.
func WithFromAddr(fromAddr string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.FromAddr = fromAddr
		return nil
	}
}

// WithToAddr indicates the receiver address from which tx are to be received.
func WithToAddr(toAddr string) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.ToAddr = toAddr
		return nil
	}
}

// WithBlockRange indicates the block range.
func WithBlockRange(startBlock, endBlock string) BlockEventOption {
	return func(f *blockEventOption) error {
		if f.blockFilter.Range == nil {
			f.blockFilter.Range = &pb.BlockRange{}
		}
		f.blockFilter.Range.Start = startBlock
		f.blockFilter.Range.End = endBlock
		return nil
	}
}

// WithExcludeTx indicates if exclude tx.
func WithExcludeTx(excludeTx bool) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.ExcludeTx = excludeTx
		return nil
	}
}

// WithExcludeTxEvent indicates if exclude tx event.
func WithExcludeTxEvent(excludeTxEvent bool) BlockEventOption {
	return func(f *blockEventOption) error {
		f.blockFilter.ExcludeTxEvent = excludeTxEvent
		return nil
	}
}
