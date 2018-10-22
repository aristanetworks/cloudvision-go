// Copyright (c) 2015 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"fmt"
	"time"

	"arista/schema"
	"arista/types"

	"github.com/aristanetworks/goarista/key"
)

// Mode is a enum that determines what mode the provider is currently in,
// and that in turn determines the ordering of Notifications sent out
// during entity creation (typically consists of a Notification about child's
// initial attributes and a Notification linking the child to its parent).
type Mode int

const (
	// AgentMode is used for all Go agents except TerminAttr.
	// The ordering of Notifications sent out during entity creation is
	// first Notification links child to its parent,
	// second Notification informs the child's initial attributes.
	AgentMode Mode = iota
	// StreamingMode is used for TerminAttr only.
	// The ordering of Notifications sent out during entity creation is
	// first Notification informs the child's initial attributes,
	// second Notification links child to its parent.
	StreamingMode
)

// Default mode for providerMode is AgentMode.
// Setter for providerMode is func SetMode(mode Mode) below.
// Getter for providerMode is func GetMode() below.
var providerMode = AgentMode

// A EOSProvider "owns" certain entities.  There are providers for entities
// coming from different sources (from Sysdb, from Smash, from /proc, etc.).
// Providers typically run in their own Goroutine(s), e.g. to read from the
// socket from Sysdb or from the shared memory files for Smash.  Providers can
// be asked to stop.  They also have a method used to write an update back to
// the source (e.g. send a message to Sysdb or update a shared-memory file for
// Smash).  Some providers can be read-only (e.g. the provider exposing data
// from /proc).
type EOSProvider interface {
	Provider

	// Init does necessary initialization on the provider.
	Init(s *schema.Schema, root types.Entity, notification chan<- types.Notification)

	// WaitForNotification() waits for a provider to be able to stream out notifications.
	WaitForNotification()

	// Write asks the provider to apply the updates carried by the given
	// Notification to its data source (e.g. by sending an update to Sysdb
	// or updating a Smash table, etc.).
	Write(notif types.Notification) error

	// InstantiateChild asks the provider to instantiate the new child
	// entity in the provider's data source.  k is the key in the
	// parent's collection that this entity is being instantiated
	// in. If the entity is not part of a collection k should be nil.
	// The given timestamp is used in the notification(s) emitted.
	// Can return ErrParentNotFound.
	InstantiateChild(ts time.Time, child types.Entity, attrDef *types.AttrDef,
		k key.Key, ctorArgs map[string]interface{}) error

	// DeleteChild asks the provider to drop the child entity.
	// attrDef is the attribute under which this entity was instantiated.
	// If the attribute is a collection, k should be set to the key in
	// that collection corresponding to this child. If the attribute
	// is a singleton k should be nil.
	// The given timestamp is used in the notification(s) emitted.
	DeleteChild(ts time.Time, child types.Entity, attrDef *types.AttrDef, k key.Key) error
}

// EntityExistor can be optionally implemented by Providers. It
// provides a way to check if an entity is supposed to exist. This is
// used in testing to look for leaks of entities.
type EntityExistor interface {
	EntityExists(e types.Entity) bool
}

// ErrParentNotFound comes from InstantiateChild when the child's
// parent is unknown.
type ErrParentNotFound struct {
	childPath  key.Path
	parentPath key.Path
}

// ErrParentIsNil comes from InstantiateChild when the child's
// parent is nil.
type ErrParentIsNil struct {
	childPath key.Path
}

// NewErrParentNotFound creates a new ErrParentNotFound.
func NewErrParentNotFound(childPath, parentPath key.Path) error {
	return &ErrParentNotFound{
		childPath:  childPath,
		parentPath: parentPath,
	}
}

// NewErrParentIsNil creates a new ErrParentIsNil.
func NewErrParentIsNil(childPath key.Path) error {
	return &ErrParentIsNil{childPath: childPath}
}

func (e *ErrParentNotFound) Error() string {
	return fmt.Sprintf("Parent of %s (%s) not found", e.childPath, e.parentPath)
}

func (e *ErrParentIsNil) Error() string {
	return fmt.Sprintf("Parent of %s is nil", e.childPath)
}

// IsParentNotFound tells you if an error is a ErrParentNotFound.
func IsParentNotFound(err error) bool {
	_, ok := err.(*ErrParentNotFound)
	return ok
}

// IsParentNil tells you if an error is an ErrParentIsNil.
func IsParentNil(err error) bool {
	_, ok := err.(*ErrParentIsNil)
	return ok
}

// SetMode takes in a Mode enum and sets the global variable providerMode to the
// input enum.
func SetMode(mode Mode) {
	providerMode = mode
}

// GetMode returns the Mode of the global variable providerMode.
func GetMode() Mode {
	return providerMode
}
