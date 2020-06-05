// Copyright (c) 2020 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package grpc

import (
	"context"
	"reflect"
	"testing"
)

func TestNewAccessTokenCredential(t *testing.T) {
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
