// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

type configTestCase struct {
	description string
	input       []byte
	expect      []Config
	err         error
}

func TestReadConfigsFromBytes(t *testing.T) {
	testCases := []configTestCase{
		{description: "basic device",
			input: []byte(`
         -  Device: test
            Options:
               a: b
               c: d`),
			expect: []Config{Config{Device: "test",
				Options: map[string]string{"a": "b", "c": "d"}}},
		},
		{description: "device with no ID",
			input: []byte(`
         -  Options:
               a: b
               c: d`),
			err: errors.New("Device must be specified"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			configs, err := readConfigsFromBytes(testCase.input)
			if err == nil && testCase.err != nil {
				t.Fatalf("Expect error %v but got nil", testCase.err)
			}
			if err != nil && testCase.err == nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if err != nil && !reflect.DeepEqual(configs, testCase.expect) {
				t.Fatalf("Mismatched configs:\n got:\n %+v\n expect:\n %+v",
					configs, testCase.expect)
			}
		})
	}
}
