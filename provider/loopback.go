// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"context"
	"fmt"
	"time"

	"arista/schema"
	"arista/types"

	"github.com/aristanetworks/goarista/key"
)

type loopback struct {
	ch     chan<- types.Notification
	isInit bool
}

// NewLoopback returns a new loopback provider that accepts updates and simply
// reflects them back into the given channel.  If this provider is to be used
// with an Agent (i.e. pass it to an Agent's WithProvider() option) then just
// pass nil as the channel instead.
func NewLoopback(notif chan<- types.Notification) EOSProvider {
	return &loopback{ch: notif}
}

func (l *loopback) Init(s *schema.Schema, root types.Entity, notif chan<- types.Notification) {
	if l.ch == nil {
		l.ch = notif
	}
	l.isInit = true
}

func (l *loopback) Run(ctx context.Context) error {
	if !l.isInit {
		return fmt.Errorf("provider is uninitialized")
	}
	return nil
}

func (l *loopback) WaitForNotification() {}

func (l *loopback) Write(notif types.Notification) error {
	if l.ch != nil {
		l.ch <- notif
	}
	return nil
}

func (l *loopback) InstantiateChild(ts time.Time, child types.Entity, attrDef *types.AttrDef,
	k key.Key, ctorArgs map[string]interface{}) error {
	if l.ch == nil {
		return nil
	}
	notifs := NotificationsForInstantiateChild(ts, child, attrDef, k)
	for _, n := range notifs {
		l.ch <- n
	}
	return nil
}

func (l *loopback) DeleteChild(ts time.Time, child types.Entity, attrDef *types.AttrDef,
	k key.Key) error {
	notifs, err := NotificationsForDeleteChild(ts, child, attrDef, k)
	if err != nil {
		return err
	}
	for _, n := range notifs {
		l.ch <- n
	}
	return nil
}
