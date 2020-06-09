// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	authHeader = "Authorization"
	bearerFmt  = "Bearer: %s"
)

// acceptableCipherSuites is a list of safe ciphersuites that can should be used for
// TLS connections.
// The order is important. The ciphersuites should be listed from most preferred to least preferred.
// CipherSuites with 3DES should be avoided because of SWEET32 vulnerability.
var acceptableCipherSuites = []uint16{
	tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
	tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
}

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

// DialWithTokenAndCert dials a gRPC endpoint, target, with the credentials in tokenFile.
// certFile is used as the root CA if supplied, else the host's root CA set is used.
func DialWithTokenAndCert(ctx context.Context, target,
	tokenFile, certFile string) (*grpc.ClientConn, error) {
	dat, err := ioutil.ReadFile(tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read token file: %s", err)
	}
	token := strings.TrimSpace(string(dat))
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}
	callCredential := NewAccessTokenCredential(token)

	tlsConf := &tls.Config{
		CipherSuites: acceptableCipherSuites,
		MinVersion:   tls.VersionTLS12,
	}

	if certFile != "" {
		c, err := ioutil.ReadFile(certFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read cert file: %s", err)
		}

		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(c) {
			return nil, errors.New("failed to add cert gfile to pool")
		}
		tlsConf.RootCAs = pool
	}

	return grpc.DialContext(ctx, target,
		grpc.WithBlock(),
		grpc.WithPerRPCCredentials(
			callCredential,
		),
		grpc.WithTransportCredentials(
			credentials.NewTLS(tlsConf),
		),
	)
}
