// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package grpc contains utilities for interacting with CloudVision's gRPC APIs
package grpc

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

func TestAccessTokenCredential(t *testing.T) {
	cred := NewAccessTokenCredential("token")
	expectedMd := map[string]string{
		"Authorization": "Bearer: token",
	}
	md, err := cred.GetRequestMetadata(context.Background(), "this/is/a/test/uri")
	if err != nil {
		t.Fatalf("got unexpected error when retrieving metadata: %s", err)
	}
	if !reflect.DeepEqual(expectedMd, md) {
		t.Fatalf("got metadata %v but expected %v", md, expectedMd)
	}
}

func TestAuthFlag(t *testing.T) {
	tcases := []struct {
		value string
		auth  *Auth
		err   error
	}{{
		value: "???",
		err:   errors.New("unknown authentication scheme: ???"),
	}, {
		value: "token,foo",
		auth:  &Auth{typ: "token", tokenFile: "foo"},
	}, {
		value: "token,foo,ca.crt",
		auth:  &Auth{typ: "token", tokenFile: "foo", caFile: "ca.crt"},
	}, {
		value: "token",
		err:   errors.New("wrong number of parameters for token authentication"),
	}, {
		value: "token,foo,ca.crt,bar",
		err:   errors.New("wrong number of parameters for token authentication"),
	}}
	for _, tcase := range tcases {
		t.Run(tcase.value, func(t *testing.T) {
			a := new(Auth)
			err := a.Set(tcase.value)
			if !reflect.DeepEqual(tcase.err, err) {
				t.Fatalf("expected error %v but got %v", tcase.err, err)
			}
			if err == nil {
				if !reflect.DeepEqual(tcase.auth, a) {
					t.Fatalf("expected to set %#v but got %#v", tcase.auth, a)
				}
				str := a.String()
				if tcase.value != str {
					t.Fatalf("expected string form to be %s but got %s",
						tcase.value, str)
				}
			}
		})
	}
}
