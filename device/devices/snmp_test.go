// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package devices

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/aristanetworks/cloudvision-go/device"
	psnmp "github.com/aristanetworks/cloudvision-go/provider/snmp"
	"github.com/gosnmp/gosnmp"
)

type optionsTestCase struct {
	name             string
	options          map[string]string
	expectedVersion  gosnmp.SnmpVersion
	expectedV3Params *psnmp.V3Params
	expectedError    error
}

func runOptionsTest(t *testing.T, tc optionsTestCase) {
	so, err := device.SanitizedOptions(options, tc.options)
	if err != nil {
		if tc.expectedError == nil {
			t.Fatalf("Error sanitizing options: %v", err)
		} else if tc.expectedError.Error() != err.Error() {
			t.Fatalf("Unexpected error sanitizing options: %v", err)
		}
		return
	}

	s, err := newSnmp(context.Background(), so)
	if err != nil {
		if tc.expectedError == nil {
			t.Fatalf("Unexpected error in newSnmp: %v", err)
		}
		if err.Error() != tc.expectedError.Error() {
			t.Fatalf("Unexpected error in newSnmp: %v (expected: %v)",
				err, tc.expectedError)
		}
		return
	}
	ss := s.(*snmp)
	if ss.v != tc.expectedVersion {
		t.Fatalf("Expected version %v, got %v", tc.expectedVersion, ss.v)
	}
	if tc.expectedV3Params == nil {
		if ss.v3Params != nil {
			t.Fatalf("Expected no v3 params, got %v", ss.v3Params)
		}
		return
	}
	if ss.v3Params.SecurityModel != tc.expectedV3Params.SecurityModel {
		t.Fatalf("Expected security model %v, got %v",
			tc.expectedV3Params.SecurityModel, ss.v3Params.SecurityModel)
	}
	if ss.v3Params.Level != tc.expectedV3Params.Level {
		t.Fatalf("Expected security level %v, got %v",
			tc.expectedV3Params.Level, ss.v3Params.Level)
	}
	if !reflect.DeepEqual(ss.v3Params.UsmParams, tc.expectedV3Params.UsmParams) {
		t.Fatalf("Expected v3 params %v, got %v", tc.expectedV3Params.UsmParams,
			ss.v3Params.UsmParams)
	}
}

// Use default value unless a key is of form "k=v", in which case use
// the specified value.
func selectOpt(keys ...string) map[string]string {
	options := map[string]string{
		"v":       "3",
		"address": "1.1.1.1",
		"l":       "authPriv",
		"a":       "sha",
		"A":       "apass",
		"x":       "des",
		"X":       "xpass",
		"u":       "user",
		"mibs":    "/a/b/c",
	}
	out := make(map[string]string)
	for _, k := range keys {
		ss := strings.Split(k, "=")
		vspec := ""
		if len(ss) > 1 {
			k = ss[0]
			vspec = ss[1]
		}
		if v, ok := options[k]; ok {
			if len(vspec) > 0 {
				out[k] = vspec
			} else {
				out[k] = v
			}
		}
	}
	return out
}

func usmParams(keys ...string) *gosnmp.UsmSecurityParameters {
	usm := &gosnmp.UsmSecurityParameters{}
	for _, k := range keys {
		switch k {
		case "a":
			usm.AuthenticationProtocol = gosnmp.SHA
		case "A":
			usm.AuthenticationPassphrase = "apass"
		case "x":
			usm.PrivacyProtocol = gosnmp.DES
		case "X":
			usm.PrivacyPassphrase = "xpass"
		case "u":
			usm.UserName = "user"
		}
	}
	return usm
}

func TestOptions(t *testing.T) {
	for _, tc := range []optionsTestCase{
		{
			name: "v2 sane",
			options: map[string]string{
				"v":       "2c",
				"c":       "public",
				"address": "1.1.1.1",
				"mibs":    "/a/b/c",
			},
			expectedVersion: gosnmp.Version2c,
		},
		{
			name: "v2 missing community",
			options: map[string]string{
				"v":       "2c",
				"address": "1.1.1.1",
				"mibs":    "/a/b/c",
			},
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: community string required for version 2c"),
		},
		{
			name:            "v3 authPriv sane",
			options:         selectOpt("v", "address", "l", "a", "A", "x", "X", "u", "mibs"),
			expectedVersion: gosnmp.Version3,
			expectedV3Params: &psnmp.V3Params{
				SecurityModel: gosnmp.UserSecurityModel,
				Level:         gosnmp.AuthPriv,
				UsmParams:     usmParams("a", "A", "x", "X", "u"),
			},
		},
		{
			name:            "v3 authNoPriv",
			options:         selectOpt("v", "address", "l=authNoPriv", "a", "A", "u", "mibs"),
			expectedVersion: gosnmp.Version3,
			expectedV3Params: &psnmp.V3Params{
				SecurityModel: gosnmp.UserSecurityModel,
				Level:         gosnmp.AuthNoPriv,
				UsmParams:     usmParams("a", "A", "u"),
			},
		},
		{
			name:            "v3 noAuthNoPriv",
			options:         selectOpt("v", "address", "l=noAuthNoPriv", "u", "mibs"),
			expectedVersion: gosnmp.Version3,
			expectedV3Params: &psnmp.V3Params{
				SecurityModel: gosnmp.UserSecurityModel,
				Level:         gosnmp.NoAuthNoPriv,
				UsmParams:     usmParams("u"),
			},
		},
		{
			name:    "v3 no username",
			options: selectOpt("v", "address", "l", "mibs"),
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: v3 is configured, so a username is required"),
		},
		{
			name:    "v3 auth missing auth proto",
			options: selectOpt("v", "address", "l", "A", "x", "X", "u", "mibs"),
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: auth is configured, so an authentication protocol " +
				"must be specified"),
		},
		{
			name:    "v3 priv missing priv proto",
			options: selectOpt("v", "address", "l", "a", "A", "X", "u", "mibs"),
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: privacy is configured, so a privacy protocol " +
				"must be specified"),
		},
		{
			name:    "v3 auth missing auth key",
			options: selectOpt("v", "address", "l", "a", "x", "X", "u", "mibs"),
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: auth is configured, so an authentication " +
				"key must be specified"),
		},
		{
			name:          "no mibs",
			options:       selectOpt("v", "address", "l", "a", "A", "x", "X", "u"),
			expectedError: errors.New("Required option 'mibs' not provided"),
		},
		{
			name: "0 poll interval",
			options: map[string]string{
				"v":            "2c",
				"c":            "public",
				"address":      "1.1.1.1",
				"mibs":         "/a/b/c",
				"pollInterval": "0s",
			},
			expectedError: errors.New("Configuration error for device " +
				"1.1.1.1: poll interval must be greater than 0 seconds"),
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			runOptionsTest(t, tc)
		})
	}
}
