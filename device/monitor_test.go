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
	<-dm.logCh

	// Monitor should log out nothing, since current level is lower than Debug level.
	dm.Debugf("Datasource monitor msg-2")
	got = b.String()
	if got != "" {
		t.Fatalf("Monitor last message should be empty, got: %v", got)
	}
	b.Reset()
	select { // make sure we are not pushing to the channel if loglevel is not enabled
	case d := <-dm.logCh:
		t.Fatal("unexpected message:", d)
	default:
	}

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

// This suite contains test cases to verify the functionality of the DatasourceMetrics struct,
// which manages the creation and manipulation of metrics within a datasource.
// It includes tests for creating metrics, setting metric values, and error handling.
func TestDatasourceMonitorMetrics(t *testing.T) {
	dm := newDatasourceMonitor(logrus.WithField("test", t.Name()), logrus.InfoLevel)

	// Test case to verify that the metric map is initially empty.
	if len(dm.metricMap) != 0 {
		t.Fatal("Metric map should be empty")
	}

	err := dm.CreateMetric("test1", "Number", "test1Description")
	if err != nil {
		t.Fatal("Failed to create metric")
	}

	// Test case to verify that metric can be created successfully.
	// It creates a new metric with the specified name, type, and description.
	// Then, it checks if the metric exists in the metric map.
	metricObj, ok := dm.metricMap["test1"]
	if len(dm.metricMap) != 1 || !ok {
		t.Fatal("Metric map should have test1 metric")
	}

	// Test case to update existing metric
	err = dm.CreateMetric("test1", "Number", "test1Description")
	if err == nil {
		t.Fatal("Metric should not be updated or modified")
	}

	// create few other metrics
	err = dm.CreateMetric("test2", "Seconds", "test2Description")
	if err != nil {
		t.Fatal("Failed to create metric")
	}
	err = dm.CreateMetric("test3", "MB", "test3Description")
	if err != nil {
		t.Fatal("Failed to create metric")
	}

	// Test case to verify error handling when setting a value for a non-existing metric.
	// It will set a value for a metric that doesn't exist in the metric map.
	err = dm.SetMetricInt("test4", 12)
	if !dm.metricMap["test4"].isChanged || err != nil {
		t.Fatal("Failed to create metric")
	}

	// Test case to verify int metric get added successfully without any error
	err = dm.SetMetricInt("test1", 12)
	if !dm.metricMap["test1"].isChanged || err != nil {
		t.Fatal("Failed to update metric")
	}

	// Test case will validate if isChanged flag will not change if we set same
	// value and it will only update if new value received in set request
	dm.metricMap["test1"] = metricInfo{
		value:     int64(14),
		isChanged: false,
	}
	err = dm.SetMetricInt("test1", 14)
	if dm.metricMap["test1"].isChanged || err != nil {
		t.Fatal("Failed to update metric")
	}
	err = dm.SetMetricInt("test1", 15)
	if !dm.metricMap["test1"].isChanged || err != nil {
		t.Fatal("Failed to update metric")
	}
	// Test case to verify successful updating of an existing metric with an integer value.
	// It checks if the value of the metric in the metric map is updated accordingly.
	err = dm.SetMetricInt("test1", 16)
	metricObj, ok = dm.metricMap["test1"]
	if err != nil || !ok || metricObj.value != int64(16) {
		t.Fatal("Incorrect metric data received")
	}

	// Test case to verify metric value update if it was not initialized.
	// It checks if the value of the metric in the metric map is updated accordingly.
	err = dm.CreateMetric("testNew", "test", "testNewDescription")
	if err != nil {
		t.Fatal("Failed to create metric")
	}
	err = dm.IncMetricInt("testNew", 4)
	metricObj, ok = dm.metricMap["testNew"]
	if err != nil || !ok || metricObj.value != int64(4) {
		t.Fatal("Incorrect metric data received")
	}

	// Test case to verify successful incrementing of an existing metric with an integer value.
	// It checks if the value of the metric in the metric map is updated accordingly.
	err = dm.IncMetricInt("testNew", 4)
	metricObj, ok = dm.metricMap["testNew"]
	if err != nil || !ok || metricObj.value != int64(8) {
		t.Fatal("Incorrect metric data received")
	}

	// Test case to verify successful updating of an existing metric with a string value.
	// It attempts to update the value of an existing metric with a
	// string value, which should fail. The test verifies that the correct error is returned
	// when attempting to set an incorrect value type.
	err = dm.SetMetricString("test1", "15")
	if err == nil ||
		err.Error() != "Error: Metric type mismatch for metric name:test1, "+
			"Expected type:int64 and Received type:string" {
		t.Fatal("Failed to update metric")
	}
}
