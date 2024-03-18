// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package provider

import (
	"context"
)

// SensorMetadata structure holds sensor metadata
type SensorMetadata struct {
	SensorIP       string
	SensorHostname string
}

// A Logger logs and records messages
type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Tracef(format string, args ...interface{})
}

// A MetricCollector will collect and publish metrics
type MetricCollector interface {
	SetMetricString(name string, value string) error
	SetMetricFloat(name string, value float64) error
	SetMetricInt(name string, value int64) error
	CreateMetric(name string, valueUnit string, description string) error
}

// A Monitor owns the state and metrics information.
// Datasource could extract information through this interface.
type Monitor interface {
	Logger
	MetricCollector
}

// A Provider "owns" some states on a target device streams out notifications on any
// changes to those states.
type Provider interface {

	// Run() kicks off the provider.  This method does not return until ctx
	// is cancelled or an error is encountered,
	// and is thus usually invoked by doing `go provider.Run()'.
	Run(ctx context.Context) error
}

// SensorMetadataProvider gives way for provider to initialize itself with sensor metadata.
type SensorMetadataProvider interface {
	// Init provider with metadata about sensor environment.
	Init(sensorMD *SensorMetadata)
}
