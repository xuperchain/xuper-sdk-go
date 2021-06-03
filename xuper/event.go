package xuper

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

func (r *Registration) Unregister() {
	r.exit <- struct{}{}
}
