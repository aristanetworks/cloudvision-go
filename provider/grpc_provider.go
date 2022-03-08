// Copyright (c) 2021 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package provider

import (
	"google.golang.org/grpc"
)

// A GRPCProvider is a Provider that has access to a raw gRPC connection.
type GRPCProvider interface {
	Provider

	// InitGRPC initializes the provider with a gRPC connection.
	InitGRPC(conn *grpc.ClientConn)
}
