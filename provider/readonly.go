// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"arista/types"
	"fmt"

	"github.com/aristanetworks/goarista/key"
)

// ReadOnly helps implement the update methods of a read-only provider.
type ReadOnly struct{}

// Write always fails.
func (ro ReadOnly) Write(notif types.Notification) error {
	return fmt.Errorf("cannot write to %s: path is read-only", notif.Path())
}

// InstantiateChild always fails.
func (ro ReadOnly) InstantiateChild(child types.Entity, attrDef *types.AttrDef,
	k key.Key, ctorArgs map[string]interface{}) error {
	return fmt.Errorf("cannot instantiate %s: parent entity is read-only", child.Path())
}

// DeleteChild always fails.
func (ro ReadOnly) DeleteChild(child types.Entity, attrDef *types.AttrDef, k key.Key) error {
	return fmt.Errorf("cannot delete %s: parent entity is read-only", child.Path())
}
