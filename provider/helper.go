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

	t := types.NowInMilliseconds()
	path := parent.Path()
	if attrDef.IsCollection() {
		// Use path to collection
		path += "/" + attrDef.Name
	} else {
		// Key is attribute name
		k = key.New(attrDef.Name)
	}

	var notifs []types.Notification
	// First notif removes this entity from its parent's singleton attribute or collection
	notifs = append(notifs, types.NewNotificationWithEntity(t, path, &[]key.Key{k}, nil, parent))
	// Second notif zeroes out the child's attributes.
	notifs = append(notifs, types.NewNotificationWithEntity(t, child.Path(),
		&[]key.Key{}, nil, child))

	var err error
	notifs, err = recursiveEntityDeleteNotification(notifs, child,
		child.GetDef().(*schema.TypeDef), t)
	if err != nil {
		return notifs, fmt.Errorf("Error recursively deleting entities with"+
			" notifications under %q: %s",
			child.Path(), err)
	}

	return notifs, nil
}

// recursiveEntityDeleteNotification recursively walks down a deleted entity,
// looking for and deleting any child instantiating attributes that hold entities.
// notifs is appended to and returned.
func recursiveEntityDeleteNotification(notifs []types.Notification, e types.Entity,
	def *schema.TypeDef, t types.Milliseconds) ([]types.Notification, error) {
	if !def.TypeFlags.IsEntity {
		// Should be impossible, as it would imply something wrong with the schema
		panic(fmt.Sprintf("Found an entity %#v at path %s with isEntity=false in typeDef: %#v",
			e, e.Path(), def))
	}

	for _, attr := range def.Attrs {
		if !attr.IsInstantiating {
			if attr.IsColl {
				notifs = append(notifs, types.NewNotificationWithEntity(t,
					e.Path()+"/"+attr.Name, &[]key.Key{}, nil, e))
			}
			continue
		}
		var childEntities []types.Entity
		child, found := e.GetEntity(attr.Name)
		if !found {
			// The attribute is not an entity, so it's a collection.
			// We have to do more work to get the entity(s) out
			notifs = append(notifs, types.NewNotificationWithEntity(t,
				e.Path()+"/"+attr.Name, &[]key.Key{}, nil, e))
			children := e.GetCollection(attr.Name)
			for _, key := range children.Keys() {
				child, found := children.Get(key)
				if !found {
					continue
				}
				childEntities = append(childEntities, child.(types.Entity))
			}
		} else if child != nil {
			childEntities = append(childEntities, child)
		}
		// For every child entity we found, we recursively call ourselves to look
		// for more child entities that need to be deleted, and then call
		// types.NewNotificationWithEntity to send notification regarding those deleted entities
		for _, childEntity := range childEntities {
			var err error
			notifs, err = recursiveEntityDeleteNotification(notifs, childEntity,
				childEntity.GetDef().(*schema.TypeDef), t)
			if err != nil {
				return notifs, fmt.Errorf("Error recursively deleting entities with"+
					"notifications under %q: %s",
					childEntity.Path(), err)
			}
			notifs = append(notifs, types.NewNotificationWithEntity(t,
				childEntity.Path(), &[]key.Key{}, nil, childEntity))
		}
	}
	return notifs, nil
}
