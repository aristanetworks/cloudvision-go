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
	level  log.Level
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

// BaseMetricCollector provides a base implementation for metric management.
type BaseMetricCollector struct{}

// SetMetricString sets the value of the metric with the specified name to the provided
// string value. This method returns nil as there is no default implementation.
func (m *BaseMetricCollector) SetMetricString(name string, value string) error { return nil }

// SetMetricFloat sets the value of the metric with the specified name to the provided
// float64 value. This method returns nil as there is no default implementation.
func (m *BaseMetricCollector) SetMetricFloat(name string, value float64) error { return nil }

// SetMetricInt sets the value of the metric with the specified name to the provided int64 value.
// This method returns nil as there is no default implementation.
func (m *BaseMetricCollector) SetMetricInt(name string, value int64) error { return nil }

// CreateMetric creates a new metric with the specified name, unit, and description.
func (m *BaseMetricCollector) CreateMetric(name string, valueUnit string,
	description string) error {
	return nil
}

// BaseMonitor is a not fully functional monitor
// for device creation in cmd/Collector.
type BaseMonitor struct {
	BaseLogger
	BaseMetricCollector
}

// NewBaseMonitor returns a new noop monitor for the collector.
func NewBaseMonitor(logger *log.Entry) *BaseMonitor {
	nm := &BaseMonitor{
		BaseLogger:          BaseLogger{logger: logger, level: log.InfoLevel},
		BaseMetricCollector: BaseMetricCollector{},
	}
	return nm
}
