// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/aristanetworks/cloudvision-go/internal/redirector"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ErrNoClusterTargets is returned when no clusters are returned.
var ErrNoClusterTargets = errors.New("no redirection cluster targets")

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

	target, err = resolveRedirection(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to redirect: %w", err)
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

	target, err := resolveRedirection(ctx, target, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to redirect: %w", err)
	}

	return grpc.DialContext(ctx, target, opts...)
}

// resolveRedirection uses the resolver to to get the regional endpoint. If the redirection
// fails we just return the client specified endpoint itself.
func resolveRedirection(
	ctx context.Context, target string, opts ...grpc.DialOption,
) (string, error) {
	// Allow the client to disable redirection
	if strings.ToLower(os.Getenv("CLOUDVISION_REGIONAL_REDIRECT")) == "false" {
		return target, nil
	}

	targets, err := redirector.Resolve(ctx, target, opts...)
	if err != nil {
		// Special case when we get back the Unimplemented error. This means
		// that the functionality is not yet available serverside and we should
		// just return the original target.
		if errors.Is(err, redirector.ErrUnimplemented) {
			return target, nil
		}

		return "", err
	}

	// Prevent a panic if targets is a nil slice and return an error
	// instead.
	if len(targets) < 1 || len(targets[0]) < 1 {
		return "", ErrNoClusterTargets
	}
	return targets[0][0], nil
}
