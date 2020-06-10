// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"errors"
	"flag"
	"fmt"
	"strings"
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
