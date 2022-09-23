// SPDX-License-Identifier: Apache-2.0

// Copyright 2021 PANTHEON.tech
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package vpp2106

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	govppapi "go.fd.io/govpp/api"

	"go.pantheon.tech/stonework/plugins/bfd/vppcalls"
	binapi "go.pantheon.tech/stonework/plugins/binapi/vpp2106/bfd"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/interface_types"
	"go.pantheon.tech/stonework/plugins/binapi/vpp2106/ip_types"
	"go.pantheon.tech/stonework/proto/bfd"
)

var (
	// EventDeliverTimeout defines maximum time to deliver event upstream.
	EventDeliverTimeout = time.Second

	// NotificationChanBufferSize defines size of notification channel buffer.
	NotificationChanBufferSize = 10
)

// AddBfd creates BFD session attached to the defined interface with given configuration ID.
func (h *BfdVppHandler) AddBfd(confID uint32, bfdEntry *bfd.BFD) error {
	// interface
	ifMeta, exists := h.ifIndexes.LookupByName(bfdEntry.Interface)
	if !exists {
		return fmt.Errorf("cannot configure BFD: interface %s is missing", bfdEntry.Interface)
	}

	localAddr, err := ip_types.ParseAddress(bfdEntry.GetLocalIp())
	if err != nil {
		return err
	}
	peerAddr, err := ip_types.ParseAddress(bfdEntry.GetPeerIp())
	if err != nil {
		return err
	}
	if localAddr.Af != peerAddr.Af {
		return fmt.Errorf("both IP addresses must be the same IP version")
	}
	req := &binapi.BfdUDPAdd{
		SwIfIndex:     interface_types.InterfaceIndex(ifMeta.SwIfIndex),
		DesiredMinTx:  bfdEntry.GetMinTxInterval(),
		RequiredMinRx: bfdEntry.GetMinRxInterval(),
		LocalAddr:     localAddr,
		PeerAddr:      peerAddr,
		DetectMult:    uint8(bfdEntry.GetDetectMultiplier()),
		BfdKeyID:      uint8(confID),
		ConfKeyID:     confID,
	}

	resp := &binapi.BfdUDPAddReply{}
	return h.callsChannel.SendRequest(req).ReceiveReply(resp)
}

// DeletebfdEntry removes existing BFD session.
func (h *BfdVppHandler) DeleteBfd(bfdEntry *bfd.BFD) error {
	ifMeta, exists := h.ifIndexes.LookupByName(bfdEntry.Interface)
	if !exists {
		return fmt.Errorf("cannot remove BFD: interface %s is missing", bfdEntry.Interface)
	}

	localAddr, err := ip_types.ParseAddress(bfdEntry.GetLocalIp())
	if err != nil {
		return err
	}
	peerAddr, err := ip_types.ParseAddress(bfdEntry.GetPeerIp())
	if err != nil {
		return err
	}
	if localAddr.Af != peerAddr.Af {
		return fmt.Errorf("both IP addresses must be the same IP version")
	}
	req := &binapi.BfdUDPDel{
		SwIfIndex: interface_types.InterfaceIndex(ifMeta.SwIfIndex),
		LocalAddr: localAddr,
		PeerAddr:  peerAddr,
	}

	resp := &binapi.BfdUDPDelReply{}
	return h.callsChannel.SendRequest(req).ReceiveReply(resp)
}

// DumpBfd returns retrieved BFD data together with BFD state.
func (h *BfdVppHandler) DumpBfd() ([]*vppcalls.BfdDetails, error) {
	var bfdList []*vppcalls.BfdDetails
	reqCtx := h.callsChannel.SendMultiRequest(&binapi.BfdUDPSessionDump{})
	for {
		bfdEntryDetails := &binapi.BfdUDPSessionDetails{}
		if stop, err := reqCtx.ReceiveReply(bfdEntryDetails); err != nil {
			h.log.Error(err)
			return nil, err
		} else if stop {
			break
		}
		ifName, _, exists := h.ifIndexes.LookupBySwIfIndex(uint32(bfdEntryDetails.SwIfIndex))
		if !exists {
			return nil, fmt.Errorf("BFD interface with index %d is missing", bfdEntryDetails.SwIfIndex)
		}
		config := &bfd.BFD{
			Interface:        ifName,
			LocalIp:          bfdEntryDetails.LocalAddr.String(),
			PeerIp:           bfdEntryDetails.PeerAddr.String(),
			MinTxInterval:    bfdEntryDetails.DesiredMinTx,
			MinRxInterval:    bfdEntryDetails.RequiredMinRx,
			DetectMultiplier: uint32(bfdEntryDetails.DetectMult),
		}

		bfdList = append(bfdList, &vppcalls.BfdDetails{
			Config:          config,
			State:           stateToProto(bfdEntryDetails.State),
			ConfKey:         bfdEntryDetails.ConfKeyID,
			BfdKey:          bfdEntryDetails.BfdKeyID,
			IsAuthenticated: bfdEntryDetails.IsAuthenticated,
		})
	}

	return bfdList, nil
}

// WatchBfdEvents starts BFD event watcher.
func (h *BfdVppHandler) WatchBfdEvents(ctx context.Context, eventChan chan<- *bfd.BFDEvent) error {
	notificationChan := make(chan govppapi.Message, NotificationChanBufferSize)

	// subscribe to BFD notifications
	sub, err := h.callsChannel.SubscribeNotification(notificationChan, &binapi.BfdUDPSessionDetails{})
	if err != nil {
		return fmt.Errorf("subscribing to VPP notification (bfd_session_event) failed: %v", err)
	}
	unsubscribe := func() {
		if err := sub.Unsubscribe(); err != nil {
			h.log.Warnf("unsubscribing VPP notification (bfd_session_event) failed: %v", err)
		}
	}

	go func() {
		h.log.Debugf("start watching BFD events")
		defer h.log.Debugf("done watching BFD events (%v)", ctx.Err())

		for {
			select {
			case e, open := <-notificationChan:
				if !open {
					h.log.Debugf("BFD events channel was closed")
					unsubscribe()
					return
				}

				bfdEvent, ok := e.(*binapi.BfdUDPSessionDetails)
				if !ok {
					h.log.Debugf("unexpected notification type: %#v", bfdEvent)
					continue
				}

				event, err := h.toBfdEvent(bfdEvent)
				if err != nil {
					h.log.Warn(err)
					continue
				}

				select {
				case eventChan <- event:
					// ok
				case <-ctx.Done():
					unsubscribe()
					return
				default:
					// in case the channel is full
					go func() {
						select {
						case eventChan <- event:
							// sent ok
						case <-time.After(EventDeliverTimeout):
							h.log.Warnf("BFD (conf-ID: %d) event dropped, cannot deliver", bfdEvent.ConfKeyID)
						}
					}()
				}
			case <-ctx.Done():
				unsubscribe()
				return
			}
		}
	}()

	// enable BFD events from VPP
	req := &binapi.WantBfdEvents{
		PID:           uint32(os.Getpid()),
		EnableDisable: true,
	}
	resp := &binapi.WantBfdEventsReply{}
	err = h.callsChannel.SendRequest(req).ReceiveReply(resp)
	// do not return error on repeated subscribe attempt
	if errors.Is(err, govppapi.VPPApiError(govppapi.INVALID_REGISTRATION)) {
		h.log.Debugf("already subscribed to BFD events: %v", err)
		return nil
	}
	return err
}

func (h *BfdVppHandler) toBfdEvent(bfdEvent *binapi.BfdUDPSessionDetails) (*bfd.BFDEvent, error) {
	ifName, _, exists := h.ifIndexes.LookupBySwIfIndex(uint32(bfdEvent.SwIfIndex))
	if !exists {
		return nil, fmt.Errorf("BFD event for unknown interface (sw_if_index: %d)",
			bfdEvent.SwIfIndex)
	}
	event := &bfd.BFDEvent{
		Interface:    ifName,
		LocalIp:      bfdEvent.LocalAddr.String(),
		PeerIp:       bfdEvent.PeerAddr.String(),
		SessionState: stateToProto(bfdEvent.State),
	}
	return event, nil
}

func stateToProto(state binapi.BfdState) bfd.BFDEvent_SessionState {
	switch state {
	case 1:
		return bfd.BFDEvent_Down
	case 2:
		return bfd.BFDEvent_Init
	case 3:
		return bfd.BFDEvent_Up
	}
	return bfd.BFDEvent_Unknown
}
