// Copyright (c) 2020 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package workspace

import (
	"context"
	"fmt"
	"io"

	workspace_v1 "github.com/aristanetworks/cloudvision-go/api/arista/workspace.v1"
	"github.com/aristanetworks/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// WsAPIUtils is the interface for workspace API utils
type WsAPIUtils interface {
	CreateWS(ctx context.Context, wsID, displayName,
		description string) (*workspace_v1.WorkspaceConfig, error)
	WaitForWSAvailable(ctx context.Context, wsID string) error
	BuildWS(ctx context.Context, wsID,
		buildID string) (*workspace_v1.WorkspaceConfig, error)
	SubmitWS(ctx context.Context, wsID string, wait, force bool) (
		*workspace_v1.WorkspaceConfig, *workspace_v1.Workspace, error)
	AbandonWS(ctx context.Context,
		wsID string) (*workspace_v1.WorkspaceConfig, error)
	DeleteWS(ctx context.Context, wsID string) (
		*workspace_v1.WorkspaceKey, error)
	GetWS(ctx context.Context, wsID string) (
		*workspace_v1.Workspace, error)
	GetWSConfig(ctx context.Context, wsID string) (
		*workspace_v1.WorkspaceConfig, error)
	GetAllWSConfigs(ctx context.Context,
	) ([]*workspace_v1.WorkspaceConfig, error)
	GetBuild(ctx context.Context,
		wsID, buildID string) (*workspace_v1.WorkspaceBuild, error)
	GetAllBuildDetails(ctx context.Context, wsID, buildID string) (
		map[string]*workspace_v1.WorkspaceBuildDetails, error)
	WaitForBuildToFinish(ctx context.Context,
		wsID, buildID string) (*workspace_v1.WorkspaceBuild, error)
	SubscribeWsState(ctx context.Context,
		wsID string) (workspace_v1.WorkspaceService_SubscribeClient, error)
	BuildAndSubmitWS(ctx context.Context,
		wsID string) (*workspace_v1.Workspace, error)
	RebaseWS(ctx context.Context, wsID string) (*workspace_v1.Workspace, error)
}

type wsAPIUtils struct {
	wsConfigServiceClient       workspace_v1.WorkspaceConfigServiceClient
	wsServiceClient             workspace_v1.WorkspaceServiceClient
	wsBuildServiceClient        workspace_v1.WorkspaceBuildServiceClient
	wsBuildDetailsServiceClient workspace_v1.WorkspaceBuildDetailsServiceClient
}

// NewWsAPIUtils creates an instance of WsAPIUtils
func NewWsAPIUtils(conn grpc.ClientConnInterface) WsAPIUtils {
	return &wsAPIUtils{
		wsConfigServiceClient:       workspace_v1.NewWorkspaceConfigServiceClient(conn),
		wsServiceClient:             workspace_v1.NewWorkspaceServiceClient(conn),
		wsBuildServiceClient:        workspace_v1.NewWorkspaceBuildServiceClient(conn),
		wsBuildDetailsServiceClient: workspace_v1.NewWorkspaceBuildDetailsServiceClient(conn),
	}
}

// CreateWS creates a workspace
func (wsUtils *wsAPIUtils) CreateWS(ctx context.Context,
	wsID, displayName, description string) (*workspace_v1.WorkspaceConfig, error) {
	resp, err := wsUtils.wsConfigServiceClient.Set(ctx, &workspace_v1.WorkspaceConfigSetRequest{
		Value: &workspace_v1.WorkspaceConfig{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID),
			},
			DisplayName: wrapperspb.String(displayName),
			Description: wrapperspb.String(description),
		},
	})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig Set: %+v\n", resp)
	return resp.GetValue(), err
}

// WaitForWSAvailable wait for workspace state to be created
func (wsUtils *wsAPIUtils) WaitForWSAvailable(ctx context.Context, wsID string) error {
	stream, err := wsUtils.SubscribeWsState(ctx, wsID)
	if err != nil {
		return err
	}
	for {
		streamResp, err := stream.Recv()
		if err != nil {
			return err
		}
		if streamResp.Value != nil {
			break
		}
	}
	err = stream.CloseSend()
	return err
}

// BuildWS starts a build of a workspace
func (wsUtils *wsAPIUtils) BuildWS(ctx context.Context,
	wsID, buildID string) (*workspace_v1.WorkspaceConfig, error) {
	resp, err := wsUtils.wsConfigServiceClient.Set(ctx, &workspace_v1.WorkspaceConfigSetRequest{
		Value: &workspace_v1.WorkspaceConfig{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID)},
			Request: workspace_v1.Request_REQUEST_START_BUILD,
			RequestParams: &workspace_v1.RequestParams{
				RequestId: wrapperspb.String(buildID),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig Set (START_BUILD): %+v\n", resp)
	return resp.GetValue(), nil
}

// SubmitWS submits a workspace
func (wsUtils *wsAPIUtils) SubmitWS(ctx context.Context, wsID string,
	wait, force bool) (
	*workspace_v1.WorkspaceConfig, *workspace_v1.Workspace, error) {

	submitRequest := workspace_v1.Request_REQUEST_SUBMIT
	if force {
		submitRequest = workspace_v1.Request_REQUEST_SUBMIT_FORCE
	}
	uuidStruct, err := uuid.NewRandom()
	if err != nil {
		glog.Errorf("Error creating uuid struct %+v", err)
		return nil, nil, err
	}
	submitID := "submit-" + uuidStruct.String()
	resp, err := wsUtils.wsConfigServiceClient.Set(ctx, &workspace_v1.WorkspaceConfigSetRequest{
		Value: &workspace_v1.WorkspaceConfig{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID)},
			Request: submitRequest,
			RequestParams: &workspace_v1.RequestParams{
				RequestId: wrapperspb.String(submitID),
			},
		},
	})
	if err != nil {
		return nil, nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig Set (SUBMIT): %+v\n", resp)
	wsConfig := resp.GetValue()
	if !wait {
		return wsConfig, nil, nil
	}

	// Wait for submit to finish, and return Workspace state
	type result struct {
		workspace *workspace_v1.Workspace
		err       error
	}
	ch := make(chan result)
	go func() {
		stream, err := wsUtils.wsServiceClient.Subscribe(ctx, &workspace_v1.WorkspaceStreamRequest{
			PartialEqFilter: []*workspace_v1.Workspace{
				{
					Key: &workspace_v1.WorkspaceKey{
						WorkspaceId: wrapperspb.String(wsID),
					},
				},
			},
		})
		if err != nil {
			ch <- result{nil, err}
			return
		}
		var resp *workspace_v1.WorkspaceStreamResponse
		for {
			resp, err = stream.Recv()
			if err != nil {
				break
			}
			if resp != nil && resp.GetValue() != nil {
				respValue := resp.GetValue()
				if respValue.Responses != nil && respValue.Responses.Values != nil &&
					respValue.Responses.Values[submitID] != nil &&
					respValue.Responses.Values[submitID].Status ==
						workspace_v1.ResponseStatus_RESPONSE_STATUS_FAIL {
					err = fmt.Errorf("submit failed")
					glog.Errorf("submit failed reason %s",
						respValue.Responses.Values[submitID].Message.Value)
					break
				}
				if respValue.State ==
					workspace_v1.WorkspaceState_WORKSPACE_STATE_SUBMITTED {
					break
				}
			}
		}
		ch <- result{resp.GetValue(), err}
	}()
	resultVal := <-ch
	return wsConfig, resultVal.workspace, resultVal.err
}

// AbandonWS abandons a workspace
func (wsUtils *wsAPIUtils) AbandonWS(ctx context.Context,
	wsID string) (*workspace_v1.WorkspaceConfig, error) {
	uuidStruct, err := uuid.NewRandom()
	if err != nil {
		glog.Errorf("Error creating uuid struct %+v", err)
		return nil, err
	}

	abandonID := "abandon-" + uuidStruct.String()

	resp, err := wsUtils.wsConfigServiceClient.Set(ctx, &workspace_v1.WorkspaceConfigSetRequest{
		Value: &workspace_v1.WorkspaceConfig{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID)},
			Request: workspace_v1.Request_REQUEST_ABANDON,
			RequestParams: &workspace_v1.RequestParams{
				RequestId: wrapperspb.String(abandonID),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig Set (ABANDON): %+v\n", resp)
	// Not waiting for it to complete
	return resp.GetValue(), nil
}

// DeleteWS deletes a workspace
func (wsUtils *wsAPIUtils) DeleteWS(ctx context.Context, wsID string) (
	*workspace_v1.WorkspaceKey, error) {
	resp, err := wsUtils.wsConfigServiceClient.Delete(ctx,
		&workspace_v1.WorkspaceConfigDeleteRequest{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID)},
		})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig Delete: %+v\n", resp)
	return resp.GetKey(), nil
}

// GetWS returns a Workspace resource
func (wsUtils *wsAPIUtils) GetWS(ctx context.Context, wsID string) (
	*workspace_v1.Workspace, error) {
	resp, err := wsUtils.wsServiceClient.GetOne(ctx, &workspace_v1.WorkspaceRequest{
		Key: &workspace_v1.WorkspaceKey{
			WorkspaceId: wrapperspb.String(wsID),
		},
	})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from Workspace GetOne(): %+v\n", resp)
	return resp.GetValue(), nil
}

// GetWSConfig returns a WorkspaceConfig resource
func (wsUtils *wsAPIUtils) GetWSConfig(ctx context.Context, wsID string) (
	*workspace_v1.WorkspaceConfig, error) {
	resp, err := wsUtils.wsConfigServiceClient.GetOne(ctx, &workspace_v1.WorkspaceConfigRequest{
		Key: &workspace_v1.WorkspaceKey{
			WorkspaceId: wrapperspb.String(wsID),
		},
	})
	if err != nil {
		return nil, err
	}
	glog.V(4).Infof("Response from WorkspaceConfig GetOne(): %+v\n", resp)
	return resp.GetValue(), nil
}

// GetAllWSConfigs retrieves all WorkspaceConfig resources and returns them
func (wsUtils *wsAPIUtils) GetAllWSConfigs(ctx context.Context) (
	[]*workspace_v1.WorkspaceConfig, error) {
	stream, err := wsUtils.wsConfigServiceClient.GetAll(ctx,
		&workspace_v1.WorkspaceConfigStreamRequest{})
	if err != nil {
		return nil, err
	}
	var wsConfigs []*workspace_v1.WorkspaceConfig
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return wsConfigs, err
		}
		glog.V(4).Infof("From WorkspaceConfig GetAll stream: %+v", resp)
		wsConfigs = append(wsConfigs, resp.GetValue())
	}
	return wsConfigs, nil
}

// GetBuild retrieves a WS build result
func (wsUtils *wsAPIUtils) GetBuild(ctx context.Context,
	wsID, buildID string) (*workspace_v1.WorkspaceBuild, error) {
	key := workspace_v1.WorkspaceBuildKey{
		WorkspaceId: wrapperspb.String(wsID),
		BuildId:     wrapperspb.String(buildID),
	}
	resp, err := wsUtils.wsBuildServiceClient.GetOne(ctx, &workspace_v1.WorkspaceBuildRequest{
		Key: &key,
	})
	if err != nil {
		return nil, err
	}
	return resp.GetValue(), nil
}

// WaitForBuildToFinish waits for a build to finish, i.e. succeed or fail.
// If the build is not found, it also waits for it with the assumption that it will start soon.
// Note that nil err does not mean build pass.
func (wsUtils *wsAPIUtils) WaitForBuildToFinish(ctx context.Context,
	wsID, buildID string) (*workspace_v1.WorkspaceBuild, error) {
	stream, err := wsUtils.wsBuildServiceClient.Subscribe(ctx,
		&workspace_v1.WorkspaceBuildStreamRequest{
			PartialEqFilter: []*workspace_v1.WorkspaceBuild{
				{
					Key: &workspace_v1.WorkspaceBuildKey{
						WorkspaceId: wrapperspb.String(wsID),
						BuildId:     wrapperspb.String(buildID),
					},
				},
			},
		})
	if err != nil {
		return nil, err
	}
	var resp *workspace_v1.WorkspaceBuildStreamResponse
	for {
		resp, err = stream.Recv()
		if err != nil {
			return nil, err
		}
		if resp != nil && resp.GetValue() != nil &&
			resp.GetValue().State != workspace_v1.BuildState_BUILD_STATE_UNSPECIFIED &&
			resp.GetValue().State != workspace_v1.BuildState_BUILD_STATE_IN_PROGRESS {
			break
		}
	}
	err = stream.CloseSend()
	if err != nil {
		return resp.GetValue(), err
	}
	return resp.GetValue(), nil
}

// GetAllBuildDetails returns build details of the given build, indexed by device ID
func (wsUtils *wsAPIUtils) GetAllBuildDetails(ctx context.Context,
	wsID, buildID string) (map[string]*workspace_v1.WorkspaceBuildDetails, error) {
	getAllResp, err := wsUtils.wsBuildDetailsServiceClient.GetAll(ctx,
		&workspace_v1.WorkspaceBuildDetailsStreamRequest{
			PartialEqFilter: []*workspace_v1.WorkspaceBuildDetails{
				{Key: &workspace_v1.WorkspaceBuildDetailsKey{
					WorkspaceId: &wrapperspb.StringValue{Value: wsID},
					BuildId:     &wrapperspb.StringValue{Value: buildID},
				}},
			},
		})
	if err != nil {
		return nil, err
	}
	allBuildDetails := map[string]*workspace_v1.WorkspaceBuildDetails{}
	for {
		resp, err := getAllResp.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return allBuildDetails, err
		}
		allBuildDetails[resp.GetValue().GetKey().GetDeviceId().GetValue()] = resp.GetValue()
	}
	return allBuildDetails, nil
}

// SubscribeWsState issues a subscribe on the Workspace state resource
func (wsUtils *wsAPIUtils) SubscribeWsState(ctx context.Context,
	wsID string) (workspace_v1.WorkspaceService_SubscribeClient, error) {
	return wsUtils.wsServiceClient.Subscribe(ctx, &workspace_v1.WorkspaceStreamRequest{
		PartialEqFilter: []*workspace_v1.Workspace{
			{
				Key: &workspace_v1.WorkspaceKey{
					WorkspaceId: wrapperspb.String(wsID),
				},
			},
		},
	})
}

// BuildAndSubmitWS builds and submits a workspace if the build is successful
func (wsUtils *wsAPIUtils) BuildAndSubmitWS(ctx context.Context,
	wsID string) (*workspace_v1.Workspace, error) {
	uuidStruct, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error generating a random number for build ID: %v", err)
	}
	buildID := "build-" + uuidStruct.String()
	_, err = wsUtils.BuildWS(ctx, wsID, buildID)
	if err != nil {
		glog.Errorf("WS %s, starting build %s failed: %+v", wsID, buildID, err)
		return nil, err
	}
	br, err := wsUtils.WaitForBuildToFinish(ctx, wsID, buildID)
	if err != nil {
		glog.Errorf("WS %s, build %s did not finish after wait: %+v", wsID, buildID, err)
		return nil, err
	}
	if br.State != workspace_v1.BuildState_BUILD_STATE_SUCCESS {
		msg := br.GetError().GetValue()
		err := fmt.Errorf("WS %s, build %s did not succeed (%s): %q", wsID, buildID, br.State,
			msg)
		glog.Error(err)
		return nil, err
	}

	// Submit the workspace and wait for the submission to be completed
	_, ws, err := wsUtils.SubmitWS(ctx, wsID, true, false)
	if err != nil {
		glog.Errorf("WS %s, submission failed: %q", wsID, err)
		return nil, err
	}
	return ws, nil
}

func (wsUtils *wsAPIUtils) RebaseWS(ctx context.Context, wsID string) (
	*workspace_v1.Workspace, error) {
	uuidStruct, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("error creating uuid struct: %v", err)
	}
	rebaseID := "rebase-" + uuidStruct.String()
	_, err = wsUtils.wsConfigServiceClient.Set(ctx, &workspace_v1.WorkspaceConfigSetRequest{
		Value: &workspace_v1.WorkspaceConfig{
			Key: &workspace_v1.WorkspaceKey{
				WorkspaceId: wrapperspb.String(wsID)},
			Request: workspace_v1.Request_REQUEST_REBASE,
			RequestParams: &workspace_v1.RequestParams{
				RequestId: wrapperspb.String(rebaseID),
			},
		},
	})
	if err != nil {
		return nil, err
	}
	// Wait for rebase to finish, and return Workspace state
	type result struct {
		workspace *workspace_v1.Workspace
		err       error
	}
	ch := make(chan result)
	go func() {
		stream, err := wsUtils.wsServiceClient.Subscribe(ctx, &workspace_v1.WorkspaceStreamRequest{
			PartialEqFilter: []*workspace_v1.Workspace{
				{
					Key: &workspace_v1.WorkspaceKey{
						WorkspaceId: wrapperspb.String(wsID),
					},
				},
			},
		})
		if err != nil {
			ch <- result{nil, err}
			return
		}
		var resp *workspace_v1.WorkspaceStreamResponse
		for {
			resp, err = stream.Recv()
			if err != nil {
				break
			}
			if resp != nil && resp.GetValue() != nil {
				respValue := resp.GetValue()
				if respValue.Responses != nil && respValue.Responses.Values != nil &&
					respValue.Responses.Values[rebaseID] != nil {
					if respValue.Responses.Values[rebaseID].Status ==
						workspace_v1.ResponseStatus_RESPONSE_STATUS_FAIL {
						err = fmt.Errorf("rebase failed")
						glog.Errorf("rebase failed reason %s",
							respValue.Responses.Values[rebaseID].Message.Value)
						break
					}
					if respValue.Responses.Values[rebaseID].Status ==
						workspace_v1.ResponseStatus_RESPONSE_STATUS_SUCCESS {
						break
					}
				}
			}
		}
		ch <- result{resp.GetValue(), err}
	}()
	resultVal := <-ch
	return resultVal.workspace, resultVal.err
}
