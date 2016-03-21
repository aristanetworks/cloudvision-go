// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"arista/schema"
	"arista/types"

	"github.com/aristanetworks/goarista/key"
)

type loopback chan<- types.Notification

// NewLoopback returns a new loopback provider that accepts updates and simply
// reflects them back into the given channel.
func NewLoopback(notif chan<- types.Notification) Provider {
	return loopback(notif)
}

func (l loopback) Run(s *schema.Schema, root types.Entity, notif chan<- types.Notification) {
	return
}

func (l loopback) WaitForNotification() {
	return
}

func (l loopback) Stop() {
	return
}

func (l loopback) Write(notif types.Notification) error {
	l <- notif
	return nil
}

func (l loopback) InstantiateChild(child types.Entity, attrDef *schema.AttrDef,
	k key.Key, ctorArgs map[string]interface{}) error {
	notifs := NotificationsForInstantiateChild(child, attrDef, k)
	for _, n := range notifs {
		l <- n
	}
	return nil
}

func (l loopback) DeleteChild(child types.Entity, attrDef *schema.AttrDef, k key.Key) error {
	notifs, err := NotificationsForDeleteChild(child, attrDef, k)
	if err != nil {
		return err
	}
	for _, n := range notifs {
		l <- n
	}
	return nil
}
