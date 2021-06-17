package xuper

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/xuperchain/xuperchain/core/pb"

	"github.com/golang/protobuf/proto"
)

// Watcher event watcher.
type Watcher struct {
	*XClient
	eventConsumerBufferSize uint
	SkipEmptyTx             bool
}

// InitWatcher new watcher instance.
func InitWatcher(xuperClient *XClient, eventConsumerBufferSize uint, skipEmptyTx bool) *Watcher {
	return &Watcher{
		xuperClient,
		eventConsumerBufferSize,
		skipEmptyTx,
	}
}

// Registration registration
type Registration struct {
	FilteredBlockChan <-chan *FilteredBlock
	exit              chan<- struct{}
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

// FromFilteredBlockPB convert pb.FilteredBlock to FilteredBlock
func FromFilteredBlockPB(pbblock *pb.FilteredBlock) *FilteredBlock {
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

// NewBlockFilter create a block filter
func NewBlockFilter(bcname string, opts ...BlockFilterOption) (blockFilter *pb.BlockFilter, err error) {
	blockFilter = &pb.BlockFilter{
		Bcname: bcname,
		Range:  &pb.BlockRange{},
	}

	for _, param := range opts {
		err := param(blockFilter)
		if err != nil {
			return nil, fmt.Errorf("option failed: %v", err)
		}
	}

	return blockFilter, nil
}

// BlockFilterOption describes a functional parameter for the New constructor
type BlockFilterOption func(*pb.BlockFilter) error

// WithContract indicates the contract name from which tx are to be received.
func WithContract(contract string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.Contract = contract
		return nil
	}
}

// WithEventName indicates the event name from which events are to be received.
func WithEventName(eventName string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.EventName = eventName
		return nil
	}
}

// WithInitiator indicates the contract initiator from which tx are to be received.
func WithInitiator(initiator string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.Initiator = initiator
		return nil
	}
}

// WithAuthRequire indicates the auth require from which tx are to be received.
func WithAuthRequire(authRequire string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.AuthRequire = authRequire
		return nil
	}
}

// WithFromAddr indicates the transfer address from which tx are to be received.
func WithFromAddr(fromAddr string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.FromAddr = fromAddr
		return nil
	}
}

// WithToAddr indicates the receiver address from which tx are to be received.
func WithToAddr(toAddr string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.ToAddr = toAddr
		return nil
	}
}

// WithBlockRange indicates the block range.
func WithBlockRange(startBlock, endBlock string) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.Range.Start = startBlock
		f.Range.End = endBlock
		return nil
	}
}

// WithExcludeTx indicates if exclude tx.
func WithExcludeTx(excludeTx bool) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.ExcludeTx = excludeTx
		return nil
	}
}

// WithExcludeTxEvent indicates if exclude tx event.
func WithExcludeTxEvent(excludeTxEvent bool) BlockFilterOption {
	return func(f *pb.BlockFilter) error {
		f.ExcludeTxEvent = excludeTxEvent
		return nil
	}
}

// RegisterBlockEvent registers for block events.
// Registration.Unregister must be called when the registration is no longer needed.
//  Parameters:
//  filter is an optional filter that filters out unwanted events.
//
//  Returns:
//  the registration is used to receive events. The channel is closed when Unregister is called.
func (w *Watcher) RegisterBlockEvent(filter *pb.BlockFilter, skipEmptyTx bool) (*Registration, error) {
	buf, _ := proto.Marshal(filter)
	request := &pb.SubscribeRequest{
		Type:   pb.SubscribeType_BLOCK,
		Filter: buf,
	}

	xclient := w.XClient.esc
	stream, err := xclient.Subscribe(context.TODO(), request)
	if err != nil {
		return nil, err
	}

	filteredBlockChan := make(chan *FilteredBlock, w.eventConsumerBufferSize)
	exit := make(chan struct{})
	reg := &Registration{
		FilteredBlockChan: filteredBlockChan,
		exit:              exit,
	}

	go func() {
		defer func() {
			close(filteredBlockChan)
			if err := stream.CloseSend(); err != nil {
				log.Printf("Unregister block event failed, close stream error: %v", err)
			} else {
				log.Printf("Unregister block event success...")
			}
		}()
		for {
			select {
			case <-exit:
				return
			default:
				event, err := stream.Recv()
				if err == io.EOF {
					return
				}
				if err != nil {
					log.Printf("Get block event err: %v", err)
					return
				}
				var block pb.FilteredBlock
				err = proto.Unmarshal(event.Payload, &block)
				if err != nil {
					log.Printf("Get block event err: %v", err)
					return
				}
				if len(block.GetTxs()) == 0 && skipEmptyTx {
					continue
				}
				filteredBlockChan <- FromFilteredBlockPB(&block)
			}
		}
	}()
	return reg, nil
}

// Unregister removes the given registration and closes the event channel.
//  Parameters:
//  reg is the registration handle that was returned from RegisterBlockEvent
func (r *Registration) Unregister() {
	r.exit <- struct{}{}
}
