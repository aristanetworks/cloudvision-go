// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

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
	sanitized, err := sanitizedOptions(&testCase.devInfo, testCase.config)
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

var h1 = `Help options for device test:
  a
	option a is a required option
  b
	option b is not required (default stuff)` + "\n"

var h2 = `Help options for device anothertest:
  x
	Sets something or other. (default true)
  y
	Pointless but required.
  zzzz
	Put the device to sleep. (default false)` + "\n"

var altOpt = map[string]Option{
	"x": Option{
		Description: "Sets something or other.",
		Default:     "true",
		Required:    false,
	},
	"y": Option{
		Description: "Pointless but required.",
		Default:     "",
		Required:    true,
	},
	"zzzz": Option{
		Description: "Put the device to sleep.",
		Default:     "false",
		Required:    false,
	},
}

var anotherDeviceInfo = deviceInfo{
	name:    "anothertest",
	creator: NewTestDevice,
	options: altOpt,
}

type helpTestCase struct {
	description        string
	devInfo            deviceInfo
	expectedHelpString string
}

func runHelpTest(t *testing.T, testCase helpTestCase) {
	dh := help(testCase.devInfo)
	if dh != testCase.expectedHelpString {
		t.Fatalf("In test %s, generated help string did not match expected:\n%s\n%s",
			testCase.description, dh, testCase.expectedHelpString)
	}
}

func runHelpTests(t *testing.T, testCases []helpTestCase) {
	for _, tc := range testCases {
		runHelpTest(t, tc)
	}
}

func TestHelp(t *testing.T) {
	testCases := []helpTestCase{
		{
			description:        "TestDevice",
			devInfo:            testDeviceInfo,
			expectedHelpString: h1,
		},
		{
			description:        "anotherTestDevice",
			devInfo:            anotherDeviceInfo,
			expectedHelpString: h2,
		},
	}
	runHelpTests(t, testCases)
}
