// Copyright (c) 2024 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import "context"

// MetricTracker tracks various metrics for datasources.
type MetricTracker interface {
	// TrackDatasources tracks how many unique datasources exist in the sensor
	TrackDatasources(ctx context.Context, numDatasource int)
	// TrackDatasourceErrors tracks how many errors the datasources encounter. This metric is
	// partitioned by the type of the datasource and the error. This metric is monotonically
	// increasing
	TrackDatasourceErrors(ctx context.Context, typ string, errorType string)
	// TrackDatasourceDeploys tracks how often datasources of a particular type are deployed.
	// This metric is monotonically increasing
	TrackDatasourceDeploys(ctx context.Context, typ string)
	// TrackDatasourceRestarts tracks how often datasources of a particular type are restarted.
	// This metric is monotonically increasing
	TrackDatasourceRestarts(ctx context.Context, typ string)
}

// noopMetricTracker represents a no operation MetricTracker
type noopMetricTracker struct {
}

func (mt noopMetricTracker) TrackDatasources(ctx context.Context, numDatasource int) {}
func (mt noopMetricTracker) TrackDatasourceErrors(
	ctx context.Context, typ string, errorType string) {
}
func (mt noopMetricTracker) TrackDatasourceDeploys(ctx context.Context, typ string)  {}
func (mt noopMetricTracker) TrackDatasourceRestarts(ctx context.Context, typ string) {}
