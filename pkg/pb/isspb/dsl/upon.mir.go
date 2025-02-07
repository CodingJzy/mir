// Code generated by Mir codegen. DO NOT EDIT.

package isspbdsl

import (
	dsl "github.com/filecoin-project/mir/pkg/dsl"
	types4 "github.com/filecoin-project/mir/pkg/pb/availabilitypb/types"
	types1 "github.com/filecoin-project/mir/pkg/pb/eventpb/types"
	types "github.com/filecoin-project/mir/pkg/pb/isspb/types"
	types5 "github.com/filecoin-project/mir/pkg/pb/trantorpb/types"
	types2 "github.com/filecoin-project/mir/pkg/trantor/types"
	types3 "github.com/filecoin-project/mir/pkg/types"
)

// Module-specific dsl functions for processing events.

func UponEvent[W types.Event_TypeWrapper[Ev], Ev any](m dsl.Module, handler func(ev *Ev) error) {
	dsl.UponMirEvent[*types1.Event_Iss](m, func(ev *types.Event) error {
		w, ok := ev.Type.(W)
		if !ok {
			return nil
		}

		return handler(w.Unwrap())
	})
}

func UponPushCheckpoint(m dsl.Module, handler func() error) {
	UponEvent[*types.Event_PushCheckpoint](m, func(ev *types.PushCheckpoint) error {
		return handler()
	})
}

func UponSBDeliver(m dsl.Module, handler func(sn types2.SeqNr, data []uint8, aborted bool, leader types3.NodeID, instanceId types3.ModuleID) error) {
	UponEvent[*types.Event_SbDeliver](m, func(ev *types.SBDeliver) error {
		return handler(ev.Sn, ev.Data, ev.Aborted, ev.Leader, ev.InstanceId)
	})
}

func UponDeliverCert(m dsl.Module, handler func(sn types2.SeqNr, cert *types4.Cert, empty bool) error) {
	UponEvent[*types.Event_DeliverCert](m, func(ev *types.DeliverCert) error {
		return handler(ev.Sn, ev.Cert, ev.Empty)
	})
}

func UponNewConfig(m dsl.Module, handler func(epochNr types2.EpochNr, membership *types5.Membership) error) {
	UponEvent[*types.Event_NewConfig](m, func(ev *types.NewConfig) error {
		return handler(ev.EpochNr, ev.Membership)
	})
}
