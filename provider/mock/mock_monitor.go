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

// TestMonitor is a mock of provider Monitor interface.
type TestMonitor struct {
	TestLogger
}

// NewMockMonitor is used to create a new mock monitor for testings
func NewMockMonitor() *TestMonitor {
	m := &TestMonitor{
		TestLogger: TestLogger{logrus.WithField("sensor", "test")}}
	return m
}
