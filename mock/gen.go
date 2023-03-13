// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package mock

//go:generate go run github.com/golang/mock/mockgen -package mock -destination mock_gnmiclient.gen.go github.com/openconfig/gnmi/proto/gnmi GNMIClient
