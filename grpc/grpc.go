// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DialWithAuth dials a gRPC endpoint, target, with the provided
// authentication config and dial options.
func DialWithAuth(ctx context.Context, target string, auth *Auth, opts ...grpc.DialOption) (
	*grpc.ClientConn, error) {

	authOpts, err := auth.Configure()
	if err != nil {
		return nil, fmt.Errorf("failed to configuration authentication scheme: %s", err)
	}
	opts = append(opts, grpc.WithBlock())
	return grpc.DialContext(ctx, target, append(opts, authOpts...)...)
}

// DialWithToken dials a gRPC endpoint, target, with the provided
// token and dial options.
func DialWithToken(ctx context.Context, target, token string, opts ...grpc.DialOption) (
	*grpc.ClientConn, error) {

	opts = append(opts,
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(NewAccessTokenCredential(token)),
		grpc.WithTransportCredentials(credentials.NewTLS(TLSConfig())),
	)
	return grpc.DialContext(ctx, target, opts...)
}
