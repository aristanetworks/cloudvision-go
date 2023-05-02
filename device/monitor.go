// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"github.com/sirupsen/logrus"
)

type datasourceLogger struct {
	logger *logrus.Entry
}

// Infof logs and records message in info level.
func (l *datasourceLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Errorf logs internal error message
func (l *datasourceLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debugf logs message only when in debug level.
func (l *datasourceLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Tracef logs message only when in trace level.
func (l *datasourceLogger) Tracef(format string, args ...interface{}) {
	l.logger.Tracef(format, args...)
}

type datasourceMonitor struct {
	datasourceLogger
}
