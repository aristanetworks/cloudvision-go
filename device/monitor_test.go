// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestPackageFile(t *testing.T) {
	for _, tc := range []struct {
		input    string
		expected string
	}{
		{
			input:    "a/b/c/d/e.go",
			expected: "d/e.go",
		},
		{
			input:    "d/e.go",
			expected: "d/e.go",
		},
		{
			input:    "e.go",
			expected: "e.go",
		},
		{
			input:    "",
			expected: "",
		},
	} {
		t.Run(tc.input, func(t *testing.T) {
			out := PackageFile(tc.input)
			if out != tc.expected {
				t.Fatalf("expected: %v, got: %v", tc.expected, out)
			}
		})
	}
}

func TestDatasourceMonitor(t *testing.T) {
	dm := newDatasourceMonitor(logrus.WithField("test", t.Name()), logrus.InfoLevel)
	var b strings.Builder
	dm.log.Out = &b

	// Monitor should log out Info level message.
	msg1 := "Datasource monitor msg-1"
	dm.Infof(msg1)
	got := b.String()
	if !strings.Contains(got, msg1) {
		t.Fatalf("Monitor last message should be: %s, got %s", msg1, got)
	}
	// Make sure we are capturing the correct source file and line
	if !strings.Contains(got, `file="device/monitor_test.go:52"`) {
		t.Fatalf("Monitor last message has wrong file and line: %s, got %s", msg1, got)
	}
	if !strings.Contains(got, `level=info`) {
		t.Fatalf("Monitor last message has wrong level, got %s", got)
	}
	b.Reset()

	// Monitor should log out nothing, since current level is lower than Debug level.
	dm.Debugf("Datasource monitor msg-2")
	got = b.String()
	if got != "" {
		t.Fatalf("Monitor last message should be empty, got: %v", got)
	}
	b.Reset()

	// Set Monitor to Debug level.
	msg3 := "Datasource monitor msg-3"
	dm.SetLoggerLevel(logrus.DebugLevel)
	dm.Debugf(msg3)
	got = b.String()
	if !strings.Contains(got, msg3) {
		t.Fatalf("Monitor last message should be: %s, got: %v", msg3, got)
	}
	if !strings.Contains(got, `level=debug`) {
		t.Fatalf("Monitor last message has wrong level, got %s", got)
	}
	b.Reset()

	// Monitor should log out Error level message, since current level is higher.
	msg4 := "Datasource monitor msg-4"
	dm.Errorf(msg4)
	got = b.String()
	if !strings.Contains(got, msg4) {
		t.Fatalf("Monitor last message should be: %s, got: %v", msg4, got)
	}
	if !strings.Contains(got, `level=error`) {
		t.Fatalf("Monitor last message has wrong level, got %s", got)
	}
	b.Reset()

	// Trace level
	dm.SetLoggerLevel(logrus.TraceLevel)
	msg5 := "Datasource monitor msg-5"
	dm.Tracef(msg5)
	got = b.String()
	if !strings.Contains(got, msg5) {
		t.Fatalf("Monitor last message should be: %s, got: %v", msg5, got)
	}
	if !strings.Contains(got, `level=trace`) {
		t.Fatalf("Monitor last message has wrong level, got %s", got)
	}
	b.Reset()
}
