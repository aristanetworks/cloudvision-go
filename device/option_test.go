// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"reflect"
	"testing"
)

var testDeviceInfo = deviceInfo{
	name:    "test",
	creator: NewTestDevice,
	options: TestDeviceOptions,
}

type optionsTestCase struct {
	description    string
	devInfo        deviceInfo
	config         map[string]string
	expectedConfig map[string]string
	shouldPass     bool
}

func runOptionsTest(t *testing.T, testCase optionsTestCase) {
	sanitized, err := sanitizedOptions(testCase.devInfo.options,
		testCase.config)
	if testCase.shouldPass && err != nil {
		t.Fatalf("Error sanitizing options in test %s", testCase.description)
	}
	if !testCase.shouldPass && err == nil {
		t.Fatalf("No error sanitizing options in test %s", testCase.description)
	} else if !testCase.shouldPass {
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
			},
			shouldPass: true,
		},
		{
			description: "extra options",
			devInfo:     testDeviceInfo,
			config: map[string]string{
				"a": "xyz",
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
	}
	runOptionsTests(t, testCases)
}
