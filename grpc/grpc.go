// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/aristanetworks/cloudvision-go/internal/redirector"
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
	opts = append(opts, authOpts...)

	if os.Getenv("CLOUDVISION_REGIONAL_REDIRECT") != "false" {
		targets, err := redirector.Resolve(ctx, target, opts...)
		if err != nil {
			return nil, err
		}
		// pick the first host until we have HA, if there are no hosts returned
		// just continue with the original target.
		if len(targets) > 1 {
			target = targets[0]
		}
	}

	return grpc.DialContext(ctx, target, opts...)
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

	if os.Getenv("CLOUDVISION_REGIONAL_REDIRECT") != "false" {
		targets, err := redirector.Resolve(ctx, target, opts...)
		if err != nil {
			return nil, err
		}
		// pick the first host until we have HA, if there are no hosts returned
		// just continue with the original target.
		if len(targets) > 1 {
			target = targets[0]
		}
	}

	return grpc.DialContext(ctx, target, opts...)
}
