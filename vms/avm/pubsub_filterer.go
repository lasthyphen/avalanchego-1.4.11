package avm

import (
	"github.com/ava-labs/avalanchego/api"
	"github.com/ava-labs/avalanchego/pubsub"
	"github.com/ava-labs/avalanchego/vms/components/djtx"
)

var _ pubsub.Filterer = &filterer{}

type filterer struct {
	tx *Tx
}

func NewPubSubFilterer(tx *Tx) pubsub.Filterer {
	return &filterer{tx: tx}
}

// Apply the filter on the addresses.
func (f *filterer) Filter(filters []pubsub.Filter) ([]bool, interface{}) {
	resp := make([]bool, len(filters))
	for _, utxo := range f.tx.UTXOs() {
		addressable, ok := utxo.Out.(djtx.Addressable)
		if !ok {
			continue
		}

		for _, address := range addressable.Addresses() {
			for i, c := range filters {
				if resp[i] {
					continue
				}
				resp[i] = c.Check(address)
			}
		}
	}
	return resp, api.JSONTxID{
		TxID: f.tx.ID(),
	}
}
