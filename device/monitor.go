// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

const (
	logChCapacity = 100
)

var (
	logMapping = map[string]logrus.Level{
		"LOG_LEVEL_ERROR": logrus.ErrorLevel,
		"LOG_LEVEL_INFO":  logrus.InfoLevel,
		"LOG_LEVEL_DEBUG": logrus.DebugLevel,
		"LOG_LEVEL_TRACE": logrus.TraceLevel,
	}
)

type datasourceLogger struct {
	log    *logrus.Logger
	logger *logrus.Entry
	logCh  chan string
}

// Infof logs and records message in info level.
func (l *datasourceLogger) Infof(format string, args ...interface{}) {
	l.logf(logrus.InfoLevel, format, args...)
}

// Errorf logs internal error message
func (l *datasourceLogger) Errorf(format string, args ...interface{}) {
	l.logf(logrus.ErrorLevel, format, args...)
}

// Debugf logs message only when in debug level.
func (l *datasourceLogger) Debugf(format string, args ...interface{}) {
	l.logf(logrus.DebugLevel, format, args...)
}

// Tracef logs message only when in trace level.
func (l *datasourceLogger) Tracef(format string, args ...interface{}) {
	l.logf(logrus.TraceLevel, format, args...)
}

func (l *datasourceLogger) logf(level logrus.Level,
	format string, args ...interface{}) {
	l.logger.Logf(level, format, args...)
	if l.log.IsLevelEnabled(level) {
		l.logCh <- fmt.Sprintf(format, args...)
	}
}

// DatasourceMonitor passes datasource context to device/manager
type datasourceMonitor struct {
	datasourceLogger
}

// SetLoggerLevel sets monitor logger level
func (dm *datasourceMonitor) SetLoggerLevel(level logrus.Level) {
	dm.log.SetLevel(level)
}

// PackageFile will find the last package and filename.
// Example: my/package/here/file.go -> here/file.go
func PackageFile(file string) string {
	pkgAndFile := file
	// Find last package and file name
	idx := strings.LastIndex(file, "/")
	if idx != -1 {
		// grab only file name
		pkgAndFile = file[idx+1:]
		// grab file package if available
		idx := strings.LastIndex(file[:idx], "/")
		if idx != 1 {
			pkgAndFile = file[idx+1:]
		}
	}
	return pkgAndFile
}

// newDatasourceMonitor returns a new datasource monitor for the datasource
func newDatasourceMonitor(logEntry *logrus.Entry, loglevel logrus.Level) *datasourceMonitor {
	logger := logrus.New()
	logger.Level = loglevel
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			// We need to skip the stack up to the original call.
			// Library or code changes could influence this number.
			// 0 is here, 1 is one level up etc.
			// Because there is a deterministic code path from the log call to here
			// we can find how many functions are in the stack and skip those.
			const skip = 8
			_, file, line, ok := runtime.Caller(skip)
			if !ok {
				file = f.File
				line = f.Line
			}
			fname := PackageFile(file)
			return "", fmt.Sprintf("%s:%d", fname, line)
		},
	})

	logEntry = logEntry.Dup()
	logEntry.Level = loglevel
	logEntry.Logger = logger
	dm := &datasourceMonitor{
		datasourceLogger: datasourceLogger{
			log:    logger,
			logger: logEntry}}
	dm.logCh = make(chan string, logChCapacity)
	dm.SetLoggerLevel(loglevel)
	return dm
}
