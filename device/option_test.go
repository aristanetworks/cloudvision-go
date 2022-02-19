// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aristanetworks/cloudvision-go/provider"
	"github.com/pkg/errors"
)

var testDeviceOptions = map[string]Option{
	"a": Option{
		Description: "option a is a required option",
		Default:     "",
		Pattern:     "[xyz]{3}",
		Required:    true,
	},
	"b": Option{
		Description: "option b is not required",
		Default:     "stuff",
		Required:    false,
	},
	"c": Option{
		Description: "option c is not required but has a pattern",
		Pattern:     "[xyz]{3}",
	},
}

var testDeviceInfo = registrationInfo{
	name:    "test",
	creator: NewTestDevice,
	options: testDeviceOptions,
}

type testDevice struct {
	deviceID string
}

func (td testDevice) Alive() (bool, error) {
	return true, nil
}

func (td testDevice) DeviceID() (string, error) {
	if td.deviceID == "" {
		return "0a0a.0a0a.0a0a", nil
	}
	return td.deviceID, nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

func (td testDevice) Type() string {
	return ""
}

func (td testDevice) IPAddr() string {
	return "192.168.5.6"
}

// NewTestDevice returns a dummy device for testing.
func NewTestDevice(map[string]string) (Device, error) {
	return testDevice{}, nil
}

type optionsTestCase struct {
	description    string
	devInfo        registrationInfo
	config         map[string]string
	expectedConfig map[string]string
	expectedError  error
	shouldPass     bool
}

func runOptionsTest(t *testing.T, testCase optionsTestCase) {
	sanitized, err := SanitizedOptions(testCase.devInfo.options,
		testCase.config)
	if testCase.shouldPass && err != nil {
		t.Fatalf("Error sanitizing options in test %s", testCase.description)
	}
	if !testCase.shouldPass && err == nil {
		t.Fatalf("No error sanitizing options in test %s", testCase.description)
	} else if !testCase.shouldPass {
		if testCase.expectedError != nil &&
			err.Error() != testCase.expectedError.Error() {
			t.Fatalf("Expected error '%s'; got '%s'", testCase.expectedError, err)
		}
		return
	}

	if !reflect.DeepEqual(sanitized, testCase.expectedConfig) {
		t.Fatalf("In test %s, sanitized config %s did not match expected %s",
			testCase.description, sanitized, testCase.expectedConfig)
	}
}

func runOptionsTests(t *testing.T, testCases []optionsTestCase) {
	for _, tc := range testCases {
		runOptionsTest(t, tc)
	}
}

func TestOptions(t *testing.T) {
	testCases := []optionsTestCase{
		{
			description: "sane options",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"a": "xyz",
				"b": "jkl",
			},
			expectedConfig: map[string]string{
				"a": "xyz",
				"b": "jkl",
				"c": "",
			},
			shouldPass: true,
		},
		{
			description: "default options",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"a": "xyz",
			},
			expectedConfig: map[string]string{
				"a": "xyz",
				"b": "stuff",
				"c": "",
			},
			shouldPass: true,
		},
		{
			description: "extra options",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"a": "zyx",
				"c": "ghi",
			},
			expectedConfig: nil,
			shouldPass:     false,
		},
		{
			description: "missing options",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"b": "jkl",
			},
			expectedConfig: nil,
			shouldPass:     false,
		},
		{
			description: "non-matching option",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"a": "xxxx",
			},
			expectedConfig: nil,
			expectedError: errors.New("Value for option 'a' ('xxxx') " +
				"does not match regular expression '[xyz]{3}'"),
			shouldPass: false,
		},
		{
			description: "default option invalid according to pattern",
			devInfo: registrationInfo{
				name:    "test2",
				creator: NewTestDevice,
				options: map[string]Option{
					"a": Option{
						Description: "option a is not required",
						Default:     "abc",
						Pattern:     "[xyz]{3}",
					},
				},
			},
			config:         map[string]string{},
			expectedConfig: nil,
			expectedError: errors.New("Default value ('abc') for option 'a' " +
				"does not match regular expression '[xyz]{3}'"),
		},
	}
	runOptionsTests(t, testCases)
}

func stringSliceEqual(s1, s2 []string) bool {
	if len(s1) != len(s2) {
		return false
	}
	sort.Strings(s1)
	sort.Strings(s2)
	for i, s := range s1 {
		if s != s2[i] {
			return false
		}
	}
	return true
}

func equalOneOf(s1 string, s2 []string) bool {
	for _, s := range s2 {
		if s1 == s {
			return true
		}
	}
	return false
}

func TestGetOptionHelpers(t *testing.T) {
	for _, tc := range []struct {
		name         string
		optionString string
		expected     string
	}{
		{
			name:         "string basic",
			optionString: "basic",
			expected:     "basic",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.optionString}
			os, err := GetStringOption("x", om)
			if err != nil {
				t.Fatalf("error in GetStringOption: %s", err)
			}
			if os != tc.expected {
				t.Fatalf("unexpected value from GetStringOption: %s", os)
			}
		})
	}

	for _, tc := range []struct {
		name         string
		optionString string
		expected     bool
		expectedErr  bool
	}{
		{
			name:         "true",
			optionString: "true",
			expected:     true,
		},
		{
			name:         "false",
			optionString: "false",
			expected:     false,
		},
		{
			name:         "rtue",
			optionString: "rtue",
			expectedErr:  true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.optionString}
			b, err := GetBoolOption("x", om)
			if (err == nil) == tc.expectedErr {
				t.Fatalf("expected error: %v; got: %s", tc.expectedErr, err)
			}
			if err != nil {
				return
			}
			if b != tc.expected {
				t.Fatalf("expected '%v', got '%v'", tc.expected, b)
			}
		})
	}

	for _, tc := range []struct {
		name            string
		address         string
		expectedAddress []string
		errorExpected   bool
	}{
		{
			name:            "1.1.1.1",
			address:         "1.1.1.1",
			expectedAddress: []string{"1.1.1.1"},
		},
		{
			name:          "1.1.1.1111",
			address:       "1.1.1.1111",
			errorExpected: true,
		},
		{
			name:          "port",
			address:       "1.1.1.1:123",
			errorExpected: true,
		},
		{
			name:            "hostname",
			address:         "localhost",
			expectedAddress: []string{"127.0.0.1", "::1"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.address}
			addr, err := GetAddressOption("x", om)
			if (err == nil) == tc.errorExpected {
				t.Fatalf("expected error: %v; got: %s", tc.errorExpected, err)
			}
			if err != nil {
				return
			}
			if !equalOneOf(addr, tc.expectedAddress) {
				t.Fatalf("expected address: %s; got: %s", tc.expectedAddress, addr)
			}
		})
	}

	for _, tc := range []struct {
		name          string
		portString    string
		expectedPort  string
		errorExpected bool
	}{
		{
			name:         "123",
			portString:   "123",
			expectedPort: "123",
		},
		{
			name:          "1234567",
			portString:    "1234567",
			errorExpected: true,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.portString}
			port, err := GetPortOption("x", om)
			if (err == nil) == tc.errorExpected {
				t.Fatalf("expected error: %v; got: %s", tc.errorExpected, err)
			}
			if err != nil {
				return
			}
			if port != tc.expectedPort {
				t.Fatalf("expected port: %s; got: %s", tc.expectedPort, port)
			}
		})
	}

	for _, tc := range []struct {
		name             string
		durationString   string
		expectedDuration time.Duration
		errorExpected    bool
	}{
		{
			name:             "10s",
			durationString:   "10s",
			expectedDuration: time.Second * 10,
		},
		{
			name:             "1m",
			durationString:   "1m",
			expectedDuration: time.Second * 60,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.durationString}
			duration, err := GetDurationOption("x", om)
			if (err == nil) == tc.errorExpected {
				t.Fatalf("expected error: %v; got: %s", tc.errorExpected, err)
			}
			if err != nil {
				return
			}
			if duration != tc.expectedDuration {
				t.Fatalf("expected: %s; got: %s", tc.expectedDuration, duration)
			}
		})
	}

	for _, tc := range []struct {
		name          string
		optionString  string
		expectedList  []string
		errorExpected bool
	}{
		{
			name:         "single value",
			optionString: "/a/b/c",
			expectedList: []string{"/a/b/c"},
		},
		{
			name:         "multiple values",
			optionString: "/a/b/c,/d/e/f,/g/h/i",
			expectedList: []string{"/a/b/c", "/d/e/f", "/g/h/i"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			om := map[string]string{"x": tc.optionString}
			ss, err := GetStringListOption("x", om)
			if (err == nil) == tc.errorExpected {
				t.Fatalf("expected error: %v; got: %s", tc.errorExpected, err)
			}
			if err != nil {
				return
			}
			if !stringSliceEqual(ss, tc.expectedList) {
				t.Fatalf("expected: %s; got: %s",
					strings.Join(tc.expectedList, " "), strings.Join(ss, " "))
			}
		})
	}
}
