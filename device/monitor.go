// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"sync"

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

// metricInfo represents information about a metric, including its value, unit, description,
// and whether it has been updated with a new value.
type metricInfo struct {
	value       any
	unit        string
	description string
	isChanged   bool
}

// metricCollector represents a collection of metrics with thread-safe access.
// It implements BaseMetricCollector to create, set and update the metric.
type metricCollector struct {
	mu        sync.RWMutex
	metricMap map[string]metricInfo
}

// SetMetricString sets the value of the metric with the specified name to the provided
// string value. It delegates the storage of the metric value to the storeMetric method.
// Returns an error if the storage operation fails.
func (m *metricCollector) SetMetricString(name string, value string) error {
	return m.storeMetric(name, value)
}

// SetMetricFloat sets the value of the metric with the specified name to the provided
// float64 value. It delegates the storage of the metric value to the storeMetric method.
// Returns an error if the storage operation fails.
func (m *metricCollector) SetMetricFloat(name string, value float64) error {
	return m.storeMetric(name, value)
}

// SetMetricInt sets the value of the metric with the specified name to the provided int64 value.
// It delegates the storage of the metric value to the storeMetric method.
// Returns an error if the storage operation fails.
func (m *metricCollector) SetMetricInt(name string, value int64) error {
	return m.storeMetric(name, value)
}

// CreateMetric creates a new metric with the specified name and unit in the datasource.
// This method initializes a new metric object with the provided name and unit,
// and adds it to the datasource's metric map for future reference.
func (m *metricCollector) CreateMetric(name string, valueUnit string, description string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.metricMap[name]
	if !ok {
		m.metricMap[name] = metricInfo{
			value:       nil,
			unit:        valueUnit,
			description: description,
		}
	} else {
		return fmt.Errorf("Error: Metric already created with name:%s", name)
	}
	return nil
}

// storeMetric updates the value of the metric with the specified name in the datasource's
// metric map. It takes the metric name and value as parameters, where the value can be of any
// type (string, float64, int64). The method first checks if the metric exists in the metric map
// and if it exists then verifies that current metric value type and received metric value type
// should match else return error.
// If the metric does not exist, it returns an error indicating that the metric is not found.
// Returns nil if the metric value is successfully updated.
func (m *metricCollector) storeMetric(name string, value any) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	metricInfoValue, ok := m.metricMap[name]

	if ok {
		if metricInfoValue.value == nil ||
			reflect.TypeOf(metricInfoValue.value) == reflect.TypeOf(value) {
			if metricInfoValue.value != value {
				metricInfoValue.value = value
				metricInfoValue.isChanged = true
				m.metricMap[name] = metricInfoValue
			}
		} else {
			return fmt.Errorf("Error: Metric type mismatch for metric name:%s, "+
				"Expected type:%T and Received type:%T", name,
				metricInfoValue.value, value)
		}
	} else {
		m.metricMap[name] = metricInfo{
			value:     value,
			isChanged: true,
		}
	}
	return nil
}

// DatasourceMonitor passes datasource context to device/manager
type datasourceMonitor struct {
	datasourceLogger
	metricCollector
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
			logger: logEntry},
		metricCollector: metricCollector{metricMap: make(map[string]metricInfo, 0)},
	}
	dm.logCh = make(chan string, logChCapacity)
	dm.SetLoggerLevel(loglevel)
	return dm
}
