// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc/credentials"
)

const (
	authHeader = "Authorization"
	bearerFmt  = "Bearer: %s"
)

type accessTokenAuth struct {
	bearerToken string
}

func (a *accessTokenAuth) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {
	return map[string]string{
		authHeader: a.bearerToken,
	}, nil
}

func (a *accessTokenAuth) RequireTransportSecurity() bool { return true }

// NewAccessTokenCredential constructs a new credential from a token
func NewAccessTokenCredential(token string) credentials.PerRPCCredentials {
	return &accessTokenAuth{bearerToken: fmt.Sprintf(bearerFmt, token)}
}
