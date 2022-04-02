// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"

	"github.com/aristanetworks/cloudvision-go/device/gen"
	"github.com/sirupsen/logrus"
)

type inventoryService struct {
	inventory Inventory
}

func (i *inventoryService) Add(ctx context.Context,
	req *gen.AddRequest) (*gen.AddResponse, error) {
	logrus.Infof("InventoryService: Add request for device type: %s",
		req.DeviceConfig.DeviceType)

	info, err := NewDeviceInfo(ctx,
		&Config{
			Device:  req.DeviceConfig.DeviceType,
			Options: req.DeviceConfig.Options,
		})
	if err != nil {
		logrus.Errorf("InventoryService: error creating DeviceInfo: %s", err)
		return nil, err
	}
	if err := i.inventory.Add(info); err != nil {
		return nil, err
	}
	return &gen.AddResponse{
		DeviceInfo: &gen.DeviceInfo{
			DeviceConfig: req.DeviceConfig,
			DeviceID:     info.ID,
		},
	}, nil
}

func (i *inventoryService) Delete(ctx context.Context,
	req *gen.DeleteRequest) (*gen.DeleteResponse, error) {
	logrus.Infof("InventoryService: Delete request: %v", req)
	return &gen.DeleteResponse{}, i.inventory.Delete(req.DeviceID)
}

func (i *inventoryService) Get(ctx context.Context,
	req *gen.GetRequest) (*gen.GetResponse, error) {
	ret := &gen.GetResponse{
		DeviceInfo: &gen.DeviceInfo{},
	}
	info, err := i.inventory.Get(req.DeviceID)
	if err != nil {
		return ret, err
	}
	ret.DeviceInfo = newGenDeviceInfo(info)
	return ret, nil
}

func (i *inventoryService) List(ctx context.Context,
	req *gen.ListRequest) (*gen.ListResponse, error) {
	ret := &gen.ListResponse{}
	infos := i.inventory.List()
	for _, info := range infos {
		ret.DeviceInfos = append(ret.DeviceInfos, newGenDeviceInfo(info))
	}
	return ret, nil
}

func newGenDeviceInfo(info *Info) *gen.DeviceInfo {
	ret := &gen.DeviceInfo{}
	ret.DeviceID = info.ID
	if info.Config == nil {
		return ret
	}
	ret.DeviceConfig = &gen.DeviceConfig{DeviceType: info.Config.Device,
		Options: info.Config.Options}
	return ret
}

// NewInventoryService returns a protobuf DeviceInventoryServer from an Inventory.
func NewInventoryService(inventory Inventory) gen.DeviceInventoryServer {
	return &inventoryService{inventory: inventory}
}
