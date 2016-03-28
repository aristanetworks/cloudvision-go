// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"arista/schema"
	"arista/types"
	"fmt"

	"github.com/aristanetworks/goarista/key"
)

// NotificationsForInstantiateChild this is a helper method for
// Providers to use to generate the notifications associated with
// instantiating a child
func NotificationsForInstantiateChild(child types.Entity, attrDef *schema.AttrDef,
	k key.Key) []types.Notification {
	notifs := make([]types.Notification, 2)
	def := child.GetDef().(*schema.TypeDef)
	t := types.NowInMilliseconds()
	if def.IsDirectory() {
		// If we just created a directory, just send one notification
		// to delete-all the new directory, instead of sending the
		// directory's attributes, which are internal.
		notifs[0] = types.NewNotificationWithEntity(t, child.Path(), &[]key.Key{}, nil, child)
	} else {
		path := child.Path()
		initialAttrs := make(map[key.Key]interface{}, len(def.Attrs))
		for attrName := range def.Attrs {
			v, _ := child.GetAttribute(attrName)
			attrKey := key.New(attrName)
			if _, ok := v.(types.Collection); ok {
				// Transform any collection into a pointer.
				initialAttrs[attrKey] = types.Pointer{Pointer: path + "/" + attrName}
			} else {
				initialAttrs[attrKey] = v
			}
		}
		notifs[0] = types.NewNotificationWithEntity(t, child.Path(), nil, &initialAttrs, child)
	}
	parent := child.Parent()
	attrName := attrDef.Name
	if k == nil { // Regular attribute
		notifs[1] = types.NewUpdates(parent,
			map[key.Key]interface{}{key.New(attrName): child.Ptr()})
	} else { // Collection
		// The path to notify on is the path of the entity + "/" + the
		// collection name, *except* if we're adding an entry to a directory.
		path := parent.Path()
		if !parent.GetDef().IsDirectory() {
			path += "/" + attrName
		}
		notifs[1] = types.NewNotificationWithEntity(t, path, nil,
			&map[key.Key]interface{}{k: child.Ptr()}, parent)
	}
	return notifs
}

// NotificationsForDeleteChild is a helper for Providers. It returns
// the notifs that should be sent when an entity is deleted.
// TODO: This should be recursive, sending a delete notification for
// every entity under the one getting deleted.
func NotificationsForDeleteChild(child types.Entity, attrDef *schema.AttrDef,
	k key.Key) ([]types.Notification, error) {
	parent := child.Parent()
	if parent == nil {
		return nil, fmt.Errorf("Can't generate notifications. Entity %q has nil parent",
			child.Path())
	}

	notifs := make([]types.Notification, 2)
	t := types.NowInMilliseconds()
	path := parent.Path()
	if attrDef.IsCollection() {
		// Use path to collection
		path += "/" + attrDef.Name
	} else {
		// Key is attribute name
		k = key.New(attrDef.Name)
	}
	// First notif removes this entity from its parent's singleton attribute or collection
	notifs[0] = types.NewNotificationWithEntity(t, path, &[]key.Key{k}, nil, parent)
	// Second notif zeroes out the child's attributes.
	// TODO: We should send one of these notifs recursively for every
	// entity under this entity.
	notifs[1] = types.NewNotificationWithEntity(t, child.Path(), &[]key.Key{}, nil, child)
	return notifs, nil
}
