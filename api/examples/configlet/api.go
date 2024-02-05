// Copyright (c) 2022 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package configlet

import (
	"context"
	"io"
	"time"

	sc_v1 "github.com/aristanetworks/cloudvision-go/api/arista/configlet.v1"
	atime "github.com/aristanetworks/cloudvision-go/api/arista/time"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/aristanetworks/glog"
	"google.golang.org/grpc"
	wpb "google.golang.org/protobuf/types/known/wrapperspb"
)

// SetConfigletConfig issues a Set on the ConfigletConfig resource
func SetConfigletConfig(ctx context.Context, conn *grpc.ClientConn,
	scConfig *sc_v1.ConfigletConfig) (
	*sc_v1.ConfigletConfig, error) {
	c := sc_v1.NewConfigletConfigServiceClient(conn)
	resp, err := c.Set(ctx, &sc_v1.ConfigletConfigSetRequest{
		Value: scConfig,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetConfigletConfig retrieves a ConfigletConfig resource
func GetConfigletConfig(ctx context.Context, conn *grpc.ClientConn, wsid, configid string) (
	*sc_v1.ConfigletConfig, error) {
	c := sc_v1.NewConfigletConfigServiceClient(conn)
	resp, err := c.GetOne(ctx, &sc_v1.ConfigletConfigRequest{
		Key: &sc_v1.ConfigletKey{
			WorkspaceId: wpb.String(wsid),
			ConfigletId: wpb.String(configid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetAllConfigletConfigs retrieves all ConfigletConfig resources
func GetAllConfigletConfigs(ctx context.Context, conn *grpc.ClientConn, wsid string) (
	[]*sc_v1.ConfigletConfig, error) {
	c := sc_v1.NewConfigletConfigServiceClient(conn)
	stream, err := c.GetAll(ctx, &sc_v1.ConfigletConfigStreamRequest{
		PartialEqFilter: []*sc_v1.ConfigletConfig{
			{
				Key: &sc_v1.ConfigletKey{
					WorkspaceId: wpb.String(wsid),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var scConfigs []*sc_v1.ConfigletConfig
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			break
		}
		scConfigs = append(scConfigs, resp.GetValue())
	}
	return scConfigs, nil
}

// DeleteConfigletConfig deletes a ConfigletConfig resource
func DeleteConfigletConfig(ctx context.Context, conn *grpc.ClientConn,
	wsid, configid string) (
	*sc_v1.ConfigletConfigDeleteResponse, error) {
	c := sc_v1.NewConfigletConfigServiceClient(conn)
	resp, err := c.Delete(ctx, &sc_v1.ConfigletConfigDeleteRequest{
		Key: &sc_v1.ConfigletKey{
			WorkspaceId: wpb.String(wsid),
			ConfigletId: wpb.String(configid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetConfiglet retrieves a Configlet resource
func GetConfiglet(ctx context.Context, conn *grpc.ClientConn, wsid, configid string) (
	*sc_v1.Configlet, error) {
	c := sc_v1.NewConfigletServiceClient(conn)
	resp, err := c.GetOne(ctx, &sc_v1.ConfigletRequest{
		Key: &sc_v1.ConfigletKey{
			WorkspaceId: wpb.String(wsid),
			ConfigletId: wpb.String(configid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetAllConfiglets retrieves all Configlet resources
func GetAllConfiglets(ctx context.Context, conn *grpc.ClientConn,
	wsid string, ts int64) ([]*sc_v1.Configlet, error) {
	c := sc_v1.NewConfigletServiceClient(conn)
	req := sc_v1.ConfigletStreamRequest{
		PartialEqFilter: []*sc_v1.Configlet{
			{
				Key: &sc_v1.ConfigletKey{
					WorkspaceId: wpb.String(wsid),
				},
			},
		},
		Time: &atime.TimeBounds{Start: timestamppb.New(time.Unix(ts, 0)),
			End: timestamppb.New(time.Now())},
	}
	stream, err := c.GetAll(ctx, &req)
	if err != nil {
		return nil, err
	}
	var scConfigs []*sc_v1.Configlet
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			break
		}
		scConfigs = append(scConfigs, resp.GetValue())
	}
	return scConfigs, nil
}

// SetConfigletAssignmentConfig issues a Set on the ConfigletAssignmentConfig resource
func SetConfigletAssignmentConfig(ctx context.Context, conn *grpc.ClientConn,
	cnConfig *sc_v1.ConfigletAssignmentConfig) (
	*sc_v1.ConfigletAssignmentConfig, error) {
	c := sc_v1.NewConfigletAssignmentConfigServiceClient(conn)
	resp, err := c.Set(ctx, &sc_v1.ConfigletAssignmentConfigSetRequest{
		Value: cnConfig,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetConfigletAssignmentConfig retrieves a ConfigletAssignmentConfig resource
func GetConfigletAssignmentConfig(ctx context.Context, conn *grpc.ClientConn, wsid, nodeid string) (
	*sc_v1.ConfigletAssignmentConfig, error) {
	c := sc_v1.NewConfigletAssignmentConfigServiceClient(conn)
	resp, err := c.GetOne(ctx, &sc_v1.ConfigletAssignmentConfigRequest{
		Key: &sc_v1.ConfigletAssignmentKey{
			WorkspaceId:           wpb.String(wsid),
			ConfigletAssignmentId: wpb.String(nodeid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetAllConfigletAssignmentConfigs retrieves all ConfigletAssignmentConfig resources
func GetAllConfigletAssignmentConfigs(ctx context.Context, conn *grpc.ClientConn, wsid string) (
	[]*sc_v1.ConfigletAssignmentConfig, error) {
	c := sc_v1.NewConfigletAssignmentConfigServiceClient(conn)
	stream, err := c.GetAll(ctx, &sc_v1.ConfigletAssignmentConfigStreamRequest{
		PartialEqFilter: []*sc_v1.ConfigletAssignmentConfig{
			{
				Key: &sc_v1.ConfigletAssignmentKey{
					WorkspaceId: wpb.String(wsid),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var cnConfigs []*sc_v1.ConfigletAssignmentConfig
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			break
		}
		cnConfigs = append(cnConfigs, resp.GetValue())
	}
	return cnConfigs, nil
}

// DeleteConfigletAssignmentConfig deletes a ConfigletAssignmentConfig resource
func DeleteConfigletAssignmentConfig(ctx context.Context, conn *grpc.ClientConn,
	wsid, nodeid string) (
	*sc_v1.ConfigletAssignmentConfigDeleteResponse, error) {
	c := sc_v1.NewConfigletAssignmentConfigServiceClient(conn)
	resp, err := c.Delete(ctx, &sc_v1.ConfigletAssignmentConfigDeleteRequest{
		Key: &sc_v1.ConfigletAssignmentKey{
			WorkspaceId:           wpb.String(wsid),
			ConfigletAssignmentId: wpb.String(nodeid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetConfigletAssignment retrieves a ConfigletAssignment resource
func GetConfigletAssignment(ctx context.Context, conn *grpc.ClientConn, wsid, nodeid string) (
	*sc_v1.ConfigletAssignment, error) {
	c := sc_v1.NewConfigletAssignmentServiceClient(conn)
	resp, err := c.GetOne(ctx, &sc_v1.ConfigletAssignmentRequest{
		Key: &sc_v1.ConfigletAssignmentKey{
			WorkspaceId:           wpb.String(wsid),
			ConfigletAssignmentId: wpb.String(nodeid),
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// GetAllConfigletAssignments retrieves all ConfigletAssignment resources
func GetAllConfigletAssignments(ctx context.Context, conn *grpc.ClientConn, wsid string) (
	[]*sc_v1.ConfigletAssignment, error) {
	c := sc_v1.NewConfigletAssignmentServiceClient(conn)
	stream, err := c.GetAll(ctx, &sc_v1.ConfigletAssignmentStreamRequest{
		PartialEqFilter: []*sc_v1.ConfigletAssignment{
			{
				Key: &sc_v1.ConfigletAssignmentKey{
					WorkspaceId: wpb.String(wsid),
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	var cns []*sc_v1.ConfigletAssignment
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			glog.Error(err)
			break
		}
		cns = append(cns, resp.GetValue())
	}
	return cns, nil
}
