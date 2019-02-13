// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

// Metadata represents all grpc metadata about a device.
type Metadata struct {
	DeviceID   string
	OpenConfig bool
	DeviceType *string
	Alive      *bool
}

const (
	deviceIDMetadata       = "deviceID"
	openConfigMetadata     = "openConfig"
	deviceTypeMetadata     = "deviceType"
	deviceLivenessMetadata = "deviceLiveness"
)

// NewMetadata returns a metadata from an incoming context.
func NewMetadata(ctx context.Context) (Metadata, error) {
	ret := Metadata{}
	var err error
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ret, errors.Errorf("Unable to get metadata from context")
	}

	deviceIDVal := md.Get(deviceIDMetadata)
	if len(deviceIDVal) != 1 {
		return ret, errors.Errorf("Context should have device ID metadata")
	}
	ret.DeviceID = deviceIDVal[0]

	openConfigVal := md.Get(openConfigMetadata)
	if len(openConfigVal) != 1 {
		return ret, errors.Errorf("Context should have openConfig metadata")
	}
	ret.OpenConfig, err = strconv.ParseBool(openConfigVal[0])
	if err != nil {
		return ret, errors.Errorf("Error parsing openConfig value: %v", err)
	}

	deviceTypeVal := md.Get(deviceTypeMetadata)
	if len(deviceTypeVal) != 0 {
		ret.DeviceType = &deviceTypeVal[0]
	}

	deviceLivenessVal := md.Get(deviceLivenessMetadata)
	if len(deviceLivenessVal) != 0 {
		t := true
		ret.Alive = &t
	}

	return ret, nil
}
