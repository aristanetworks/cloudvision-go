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

// NotificationsForInstantiateChild this is a helper method for
// Providers to use to generate the notifications associated with
// instantiating a child
func NotificationsForInstantiateChild(child types.Entity, attrDef *schema.AttrDef,
	k key.Key) []types.Notification {
	notifs := make([]types.Notification, 2)
	def := child.GetDef().(*schema.TypeDef)
	if def.IsDirectory() {
		// If we just created a directory, just send one notification
		// to delete-all the new directory, instead of sending the
		// directory's attributes, which are internal.
		notifs[0] = types.NewNotificationWithEntity(types.NowInMilliseconds(), child.Path(),
			&[]key.Key{}, nil, child)
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
		notifs[0] = types.NewUpdates(child, initialAttrs)
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
		notifs[1] = types.NewNotificationWithEntity(types.NowInMilliseconds(),
			path, nil, &map[key.Key]interface{}{k: child.Ptr()}, parent)
	}
	return notifs
}
