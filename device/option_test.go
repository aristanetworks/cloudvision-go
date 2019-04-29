// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"reflect"
	"testing"

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

type testDevice struct{}

func (td testDevice) Alive() (bool, error) {
	return true, nil
}

func (td testDevice) DeviceID() (string, error) {
	return "0a0a.0a0a.0a0a", nil
}

func (td testDevice) Providers() ([]provider.Provider, error) {
	return nil, nil
}

// NewTestDevice returns a dummy device for testing.
func NewTestDevice(Config) (Device, error) {
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
