package fakebatchdb

import (
	msctypes "github.com/filecoin-project/mir/pkg/availability/multisigcollector/types"
	"github.com/filecoin-project/mir/pkg/dsl"
	"github.com/filecoin-project/mir/pkg/modules"
	batchdbpbdsl "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/dsl"
	batchdbpbtypes "github.com/filecoin-project/mir/pkg/pb/availabilitypb/batchdbpb/types"
	trantorpbtypes "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	tt "github.com/filecoin-project/mir/pkg/trantor/types"
	t "github.com/filecoin-project/mir/pkg/types"
)

// ModuleConfig sets the module ids. All replicas are expected to use identical module configurations.
type ModuleConfig struct {
	Self t.ModuleID // id of this module
}

type moduleState struct {
	BatchStore       map[msctypes.BatchID]batchInfo
	TransactionStore map[tt.TxID]*trantorpbtypes.Transaction
}

type batchInfo struct {
	txIDs    []tt.TxID
	metadata []byte
}

// NewModule returns a new module for a fake batch database.
// It stores all the data in memory in plain go maps.
func NewModule(mc ModuleConfig) modules.Module {
	m := dsl.NewModule(mc.Self)

	state := moduleState{
		BatchStore:       make(map[msctypes.BatchID]batchInfo),
		TransactionStore: make(map[tt.TxID]*trantorpbtypes.Transaction),
	}

	// On StoreBatch request, just store the data in the local memory.
	batchdbpbdsl.UponStoreBatch(m, func(batchID msctypes.BatchID, txIDs []tt.TxID, txs []*trantorpbtypes.Transaction, metadata []byte, origin *batchdbpbtypes.StoreBatchOrigin) error {
		state.BatchStore[batchID] = batchInfo{
			txIDs:    txIDs,
			metadata: metadata,
		}

		for i, txID := range txIDs {
			state.TransactionStore[txID] = txs[i]
		}

		batchdbpbdsl.BatchStored(m, origin.Module, origin)
		return nil
	})

	// On LookupBatch request, just check the local map.
	batchdbpbdsl.UponLookupBatch(m, func(batchID msctypes.BatchID, origin *batchdbpbtypes.LookupBatchOrigin) error {

		info, found := state.BatchStore[batchID]
		if !found {
			batchdbpbdsl.LookupBatchResponse(m, origin.Module, false, nil, origin)
			return nil
		}

		txs := make([]*trantorpbtypes.Transaction, len(info.txIDs))
		for i, txID := range info.txIDs {
			txs[i] = state.TransactionStore[txID]
		}

		batchdbpbdsl.LookupBatchResponse(m, origin.Module, true, txs, origin)
		return nil
	})

	return m
}
