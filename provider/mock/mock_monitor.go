// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package mock

import (
	"github.com/sirupsen/logrus"
)

// TestLogger is a mock of Monitor Logger interface.
type TestLogger struct {
	logger *logrus.Entry
}

// Infof logs and records message in datasource state.
func (l *TestLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Errorf logs internal error message
func (l *TestLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debugf logs message only when in debug level.
func (l *TestLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Tracef logs message only when in trace level.
func (l *TestLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

// TestMetricCollector provides a mock implementation of a metrics collection
// interface for testing purposes.
type TestMetricCollector struct{}

// SetMetricString sets a string metric value for the given metric name.
// It returns nil indicating success or an error if the operation fails.
func (m *TestMetricCollector) SetMetricString(name string, value string) error { return nil }

// SetMetricFloat sets a float64 metric value for the given metric name.
// It returns nil indicating success or an error if the operation fails.
func (m *TestMetricCollector) SetMetricFloat(name string, value float64) error { return nil }

// SetMetricInt sets an int64 metric value for the given metric name.
// It returns nil indicating success or an error if the operation fails.
func (m *TestMetricCollector) SetMetricInt(name string, value int64) error { return nil }

// CreateMetric creates a new metric with the specified name, unit, and description.
func (m *TestMetricCollector) CreateMetric(name string, valueUnit string,
	description string) error {
	return nil
}

// TestMonitor is a mock of provider Monitor interface.
type TestMonitor struct {
	TestLogger
	TestMetricCollector
}

// NewMockMonitor is used to create a new mock monitor for testings
func NewMockMonitor() *TestMonitor {
	m := &TestMonitor{
		TestLogger: TestLogger{logrus.WithField("sensor", "test")}}
	return m
}
