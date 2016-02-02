// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.

package provider

import (
	"arista/schema"
	"arista/types"
	"arista/util"

	"github.com/aristanetworks/goarista/key"
)

// NotificationsForInstantiateChild this is a helper method for
// Providers to use to generate the notifications associated with
// instantiating a child
func NotificationsForInstantiateChild(child types.Entity, attrDef *schema.AttrDef,
	k key.Key, ctorArgs map[string]interface{}) ([]types.Notification, error) {
	notifs := make([]types.Notification, 2)
	if child.GetDef().IsDirectory() {
		// If we just created a directory, just send one notification
		// to delete-all the new directory, instead of sending the
		// directory's attributes, which are internal.
		notifs[0] = types.NewNotificationWithEntity(types.NowInMilliseconds(), child.Path(),
			&[]key.Key{}, nil, child)
	} else {
		initialAttrs := util.CopyStringMapToKeyMap(ctorArgs)
		// Transform any collection into a pointer.
		var path string
		for k, v := range initialAttrs {
			if _, ok := v.(types.Collection); ok {
				if path == "" { // If we don't yet know our child's path..
					path = child.Path() // .. compute it only once.
				}
				initialAttrs[k] = types.Pointer{Pointer: path + "/" + k.Key().(string)}
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
	return notifs, nil
}
