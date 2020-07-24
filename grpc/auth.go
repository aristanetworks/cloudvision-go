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
	"flag"
	"fmt"
	"io/ioutil"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// AuthFlagUsage indicates how to use the authentication config flag.
const AuthFlagUsage = "Authentication scheme used to connect to CloudVision. " +
	"Possible values\n:" +
	"\t\"token,{token_file}[,{ca_file}]\": client-side certificate with token-based " +
	"authentication. Uses host's root CA if {ca_file} is not provided.\n"

// AuthFlag adds an authentication config flag with the name "auth".
func AuthFlag() *Auth {
	a := new(Auth)
	flag.Var(a, "auth", AuthFlagUsage)
	return a
}

// NewTokenAuth creates a new token authentication config. If caFile
// is not provided, the host's root CA will be used.
func NewTokenAuth(tokenFile, caFile string) (*Auth, error) {
	if tokenFile == "" {
		return nil, errors.New("tokenFile is required")
	}
	return &Auth{
		typ:       "token",
		tokenFile: tokenFile,
		caFile:    caFile,
	}, nil
}

// Auth holds an authentication scheme used to connect to CloudVision.
type Auth struct {
	typ       string
	tokenFile string
	caFile    string
}

const (
	authHeader = "Authorization"
	bearerFmt  = "Bearer %s"
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

type accessTokenCred struct {
	bearerToken string
}

// NewAccessTokenCredential constructs a new per-RPC credential from a token.
func NewAccessTokenCredential(token string) credentials.PerRPCCredentials {
	return &accessTokenCred{bearerToken: fmt.Sprintf(bearerFmt, token)}
}

func (a *accessTokenCred) GetRequestMetadata(ctx context.Context,
	uri ...string) (map[string]string, error) {
	return map[string]string{
		authHeader: a.bearerToken,
	}, nil
}

func (a *accessTokenCred) RequireTransportSecurity() bool { return true }

// CAFile returns the CA certificate file specified in Auth.
func (a *Auth) CAFile() string {
	return a.caFile
}

// ClientCredentials returns dial options corresponding to the client credentials.
func (a *Auth) ClientCredentials() ([]grpc.DialOption, error) {
	switch a.typ {
	case "token":
		b, err := ioutil.ReadFile(a.tokenFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read token file: %s", err)
		}
		token := strings.TrimSpace(string(b))
		if token == "" {
			return nil, errors.New("token cannot be empty")
		}
		return []grpc.DialOption{
			grpc.WithPerRPCCredentials(NewAccessTokenCredential(token)),
		}, nil
	default:
		return nil, fmt.Errorf("unknown authentication scheme: %s", a.typ)
	}
}

// Configure returns the authentication configuration as a series of gRPC dial options.
func (a *Auth) Configure() ([]grpc.DialOption, error) {
	cfg := &tls.Config{
		CipherSuites: acceptableCipherSuites,
		MinVersion:   tls.VersionTLS12,
	}
	var opts []grpc.DialOption

	clCreds, err := a.ClientCredentials()
	if err != nil {
		return nil, err
	}
	opts = append(opts, clCreds...)
	if a.caFile != "" {
		b, err := ioutil.ReadFile(a.caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read ca file: %s", err)
		}
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM(b) {
			return nil, errors.New("failed to add ca file to pool")
		}
		cfg.RootCAs = pool
	}
	return append(opts, grpc.WithTransportCredentials(credentials.NewTLS(cfg))), nil
}

var _ flag.Value = new(Auth) // flag.Value implementation check

// Set implements flag.Value.
func (a *Auth) Set(v string) error {
	s := strings.Split(v, ",")
	switch s[0] {
	case "token":
		if len(s) < 2 || len(s) > 3 {
			return errors.New("wrong number of parameters for token authentication")
		}
		a.tokenFile = s[1]
		if len(s) == 3 {
			a.caFile = s[2]
		}
	default:
		return fmt.Errorf("unknown authentication scheme: %s", s[0])
	}
	a.typ = s[0]
	return nil
}

func (a *Auth) String() string {
	s := a.typ + ","
	switch a.typ {
	case "token":
		s += a.tokenFile
	}
	if a.caFile != "" {
		s += "," + a.caFile
	}
	return s
}
