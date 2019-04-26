// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"context"
	"testing"

	"github.com/aristanetworks/goarista/test"
	"google.golang.org/grpc/metadata"
)

func TestNewMetadata(t *testing.T) {

	boolTrue := true
	deviceType := "Target"
	deviceAddress := "1.1.1.1"

	for _, tc := range []struct {
		desc       string
		shouldFail bool
		md         Metadata
		ctx        context.Context
	}{
		{
			desc: "required fields only",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				deviceIDMetadata, "id",
				typeCheckMetadata, "true",
				deviceAddressMetadata, deviceAddress,
				openConfigMetadata, "true"),
			md: Metadata{
				DeviceID:      "id",
				OpenConfig:    true,
				TypeCheck:     true,
				DeviceAddress: deviceAddress,
			},
		},
		{
			desc: "complete metadata",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				deviceIDMetadata, "id",
				typeCheckMetadata, "true",
				deviceAddressMetadata, deviceAddress,
				openConfigMetadata, "true",
				deviceTypeMetadata, deviceType,
				deviceLivenessMetadata, "true"),
			md: Metadata{
				DeviceID:      "id",
				OpenConfig:    true,
				DeviceType:    &deviceType,
				Alive:         &boolTrue,
				TypeCheck:     true,
				DeviceAddress: deviceAddress,
			},
		},
		{
			desc: "missing device ID",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				typeCheckMetadata, "true",
				deviceAddressMetadata, deviceAddress,
				openConfigMetadata, "true",
				deviceTypeMetadata, deviceType,
				deviceLivenessMetadata, "true"),
			shouldFail: true,
		},
		{
			desc: "missing openConfig",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				typeCheckMetadata, "true",
				deviceAddressMetadata, deviceAddress,
				deviceIDMetadata, "id",
				deviceTypeMetadata, deviceType,
				deviceLivenessMetadata, "true"),
			shouldFail: true,
		},
		{
			desc: "missing typeCheck",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				deviceIDMetadata, "id",
				deviceAddressMetadata, deviceAddress,
				openConfigMetadata, "true",
				deviceTypeMetadata, deviceType,
				deviceLivenessMetadata, "true"),
			shouldFail: true,
		},
		{
			desc: "missing deviceAddress",
			ctx: metadata.AppendToOutgoingContext(
				context.Background(),
				deviceIDMetadata, "id",
				typeCheckMetadata, "true",
				openConfigMetadata, "true",
				deviceTypeMetadata, deviceType,
				deviceLivenessMetadata, "true"),
			shouldFail: true,
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			data, err := NewMetadataFromOutgoing(tc.ctx)
			if err != nil && !tc.shouldFail {
				t.Fatal(err)
			}
			if err == nil && tc.shouldFail {
				t.Fatalf("Test should have error but doesn't")
			}
			if !tc.shouldFail {
				if diff := test.Diff(tc.md, data); diff != "" {
					t.Fatalf("Unexpected metadata: Diff is %v", diff)
				}
			}
		})
	}
}
