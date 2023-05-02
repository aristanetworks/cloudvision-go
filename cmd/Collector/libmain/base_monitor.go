// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package libmain

import (
	log "github.com/sirupsen/logrus"
)

// BaseLogger only logs std info or error.
type BaseLogger struct {
	logger *log.Entry
}

// Infof logs and records info message.
func (l *BaseLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Errorf logs internal error message.
func (l *BaseLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debugf logs message only when in debug level.
func (l *BaseLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Tracef logs message only when in trace level.
func (l *BaseLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

// BaseMonitor is a not fully functional monitor
// for device creation in cmd/Collector.
type BaseMonitor struct {
	BaseLogger
}

// NewBaseMonitor returns a new noop monitor for the collector.
func NewBaseMonitor(logger *log.Entry) *BaseMonitor {
	nm := &BaseMonitor{
		BaseLogger: BaseLogger{logger}}
	return nm
}
