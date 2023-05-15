// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"sync/atomic"

	log "github.com/sirupsen/logrus"
)

const (
	logChCapacity = 100
)

var (
	logMapping = map[string]log.Level{
		"LOG_LEVEL_ERROR": log.ErrorLevel,
		"LOG_LEVEL_INFO":  log.InfoLevel,
		"LOG_LEVEL_DEBUG": log.DebugLevel,
		"LOG_LEVEL_TRACE": log.TraceLevel,
	}
)

type datasourceLogger struct {
	logger *log.Entry
	level  log.Level
	logCh  chan string
}

// Infof logs and records message in info level.
func (l *datasourceLogger) Infof(format string, args ...interface{}) {
	l.logf(log.InfoLevel, format, args...)
}

// Errorf logs internal error message
func (l *datasourceLogger) Errorf(format string, args ...interface{}) {
	l.logf(log.ErrorLevel, format, args...)
}

// Debugf logs message only when in debug level.
func (l *datasourceLogger) Debugf(format string, args ...interface{}) {
	l.logf(log.DebugLevel, format, args...)
}

// Tracef logs message only when in trace level.
func (l *datasourceLogger) Tracef(format string, args ...interface{}) {
	l.logf(log.TraceLevel, format, args...)
}

func (l *datasourceLogger) getLevel() log.Level {
	return log.Level(atomic.LoadUint32((*uint32)(&l.level)))
}

func (l *datasourceLogger) isLevelEnabled(level log.Level) bool {
	return l.getLevel() >= level
}

func (l *datasourceLogger) logf(level log.Level,
	format string, args ...interface{}) {
	if l.isLevelEnabled(level) {
		l.logger.Logf(log.ErrorLevel, format, args...)
		l.logCh <- fmt.Sprintf(format, args...)
	}
}

// DatasourceMonitor passes datasource context to device/manager
type datasourceMonitor struct {
	datasourceLogger
}

// SetLoggerLevel sets monitor logger level
func (dm *datasourceMonitor) SetLoggerLevel(level log.Level) {
	atomic.StoreUint32((*uint32)(&dm.level), uint32(level))
}

// newDatasourceMonitor returns a new datasource monitor for the datasource
func newDatasourceMonitor(logger *log.Entry, loglevel log.Level) *datasourceMonitor {
	dm := &datasourceMonitor{
		datasourceLogger: datasourceLogger{
			logger: logger}}
	dm.logCh = make(chan string, logChCapacity)
	dm.SetLoggerLevel(loglevel)
	return dm
}
