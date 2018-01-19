// Copyright (c) 2016 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package provider

import (
	"fmt"
	"time"

	"arista/types"

	"github.com/aristanetworks/goarista/key"
	"github.com/aristanetworks/goarista/path"
)

// NotificationsForInstantiateChild this is a helper method for
// Providers to use to generate the notifications associated with
// instantiating a child
func NotificationsForInstantiateChild(ts time.Time, child types.Entity, attrDef *types.AttrDef,
	k key.Key) []types.Notification {
	notifs := make([]types.Notification, 2)
	def := child.GetDef()
	if def.IsDirectory() {
		// If we just created a directory, just send one notification
		// to delete-all the new directory, instead of sending the
		// directory's attributes, which are internal.
		notifs[0] = types.NewNotifWithEntity(ts, child.PathComponents(),
			[]key.Key{}, nil, child)
	} else {
		p := child.PathComponents()
		initialAttrs := make(map[key.Key]interface{}, len(def.Attrs))
		for attrName := range def.Attrs {
			v, _ := child.GetAttribute(attrName)
			attrKey := key.New(attrName)
			if _, ok := v.(types.Collection); ok {
				// Transform any collection into a pointer.
				initialAttrs[attrKey] = types.Pointer{
					Pointer: path.Append(p, attrName).String(),
				}
			} else {
				initialAttrs[attrKey] = v
			}
		}
		notifs[0] = types.NewNotifWithEntity(ts, child.PathComponents(),
			nil, initialAttrs, child)
	}
	parent := child.Parent()
	attrName := attrDef.Name
	if k == nil { // Regular attribute
		notifs[1] = types.NewUpdates(parent,
			map[key.Key]interface{}{key.New(attrName): child.Ptr()})
	} else { // Collection
		// The path to notify on is the path of the entity + "/" + the
		// collection name, *except* if we're adding an entry to a directory.
		p := parent.PathComponents()
		if !parent.GetDef().IsDirectory() {
			p = path.Append(p, attrName)
		}
		notifs[1] = types.NewNotifWithEntity(ts, p, nil,
			map[key.Key]interface{}{k: child.Ptr()}, parent)
	}
	// In "AgentMode" the ordering of the two notifications should be switched
	if GetMode() == AgentMode {
		notifs[0], notifs[1] = notifs[1], notifs[0]
	}
	return notifs
}

// NotificationsForDeleteChild is a helper for Providers. It returns
// the notifs that should be sent when an entity is deleted.
func NotificationsForDeleteChild(ts time.Time, child types.Entity, attrDef *types.AttrDef,
	k key.Key) ([]types.Notification, error) {
	parent := child.Parent()
	if parent == nil {
		return nil, fmt.Errorf("Can't generate notifications: %s",
			NewErrParentIsNil(child.PathComponents()))
	}

	p := parent.PathComponents()
	if attrDef.IsColl {
		// Use path to collection
		p = path.Append(p, attrDef.Name)
	} else {
		// Key is attribute name
		k = key.New(attrDef.Name)
	}

	notifs, err := recursiveEntityDeleteNotification(nil, child, child.GetDef(), ts)
	if err != nil {
		return notifs, fmt.Errorf("Error recursively deleting entities with"+
			" notifications under %q: %s",
			child.PathComponents(), err)
	}

	// Zero out the child's attributes.
	notifs = append(notifs, types.NewNotifWithEntity(ts, child.PathComponents(),
		[]key.Key{}, nil, child))

	// Finally remove this entity from its parent's attribute or collection
	notifs = append(notifs, types.NewNotifWithEntity(ts, p, []key.Key{k}, nil, parent))

	return notifs, nil
}

// recursiveEntityDeleteNotification recursively walks down a deleted entity,
// looking for and deleting any child instantiating attributes that hold entities.
// notifs is appended to and returned.
func recursiveEntityDeleteNotification(notifs []types.Notification, e types.Entity,
	def *types.TypeDef, ts time.Time) ([]types.Notification, error) {
	if !def.TypeFlags.IsEntity {
		// Should be impossible, as it would imply something wrong with the schema
		panic(fmt.Sprintf("Found an entity %#v at path %s with isEntity=false in typeDef: %#v",
			e, e.PathComponents(), def))
	}

	var childEntities []types.Entity
	// afterRecurseNotifs are notifs that should be added after the
	// recursive call
	var afterRecurseNotifs []types.Notification
	for _, attr := range def.Attrs {
		if !attr.IsInstantiating {
			if attr.IsColl {
				afterRecurseNotifs = append(afterRecurseNotifs, types.NewNotifWithEntity(ts,
					path.Append(e.PathComponents(), attr.Name), []key.Key{}, nil, e))
			}
			continue
		}
		if attr.IsColl {
			afterRecurseNotifs = append(afterRecurseNotifs, types.NewNotifWithEntity(ts,
				path.Append(e.PathComponents(), attr.Name), []key.Key{}, nil, e))
			children := e.GetCollection(attr.Name)
			children.ForEach(func(k key.Key, child interface{}) error {
				childEntities = append(childEntities, child.(types.Entity))
				return nil
			})
		} else if child, ok := e.GetEntity(attr.Name); ok {
			childEntities = append(childEntities, child)
		}
	}
	// For every child entity we found, we recursively call ourselves to look
	// for more child entities that need to be deleted, and then call
	// types.NewNotifWithEntity to send notification regarding those deleted entities
	for _, childEntity := range childEntities {
		var err error
		notifs, err = recursiveEntityDeleteNotification(notifs, childEntity,
			childEntity.GetDef(), ts)
		if err != nil {
			return notifs, fmt.Errorf("Error recursively deleting entities with"+
				"notifications under %q: %s",
				childEntity.PathComponents(), err)
		}
		notifs = append(notifs, types.NewNotifWithEntity(ts,
			childEntity.PathComponents(), []key.Key{}, nil, childEntity))
	}
	return append(notifs, afterRecurseNotifs...), nil
}

// NotificationsForCollectionCount is a helper method for Providers to use to
// generate the notifications associated with collection counts.
func NotificationsForCollectionCount(ts time.Time, collName string, count uint32,
	parent types.Entity) types.Notification {

	if GetMode() != StreamingMode || parent.GetDef().IsDirectory() {
		return nil
	}

	return types.NewNotifWithEntity(ts, path.Append(parent.PathComponents(), "_counts"),
		nil, map[key.Key]interface{}{key.New(collName): count}, parent)
}
