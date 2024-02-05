// Copyright (c) 2023 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package main

//  go run api/examples/configlet/cmd/main.go -token `cat ~/tmp/token`  -server cvp686:443 -v=3

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	configlet_v1 "github.com/aristanetworks/cloudvision-go/api/arista/configlet.v1"
	studio_v1 "github.com/aristanetworks/cloudvision-go/api/arista/studio.v1"
	cApi "github.com/aristanetworks/cloudvision-go/api/examples/configlet"
	ws "github.com/aristanetworks/cloudvision-go/api/examples/workspace"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/aristanetworks/cloudvision-go/api/fmp"
	cvgrpc "github.com/aristanetworks/cloudvision-go/grpc"
	"github.com/aristanetworks/glog"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var (
	serverFlag = flag.String(
		"server",
		"127.0.0.1:12054",
		"Address of CloudVision",
	)
	tokenFlag = flag.String(
		"token",
		"",
		"auth token",
	)
	deviceFlag = flag.String(
		"device",
		"E186A7D64EE72E219D6E64E4C23A2DEB",
		"Device Id",
	)
	bodyFlag = flag.String(
		"body",
		"alias exampleAlias example",
		"Configlet Body",
	)
	cleanupFlag = flag.Bool(
		"cleanup",
		false,
		"cleanup resources",
	)
)

func toWrapper(str string) *wrapperspb.StringValue {
	if str == "" {
		return nil
	}
	return wrapperspb.String(str)
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if *serverFlag == "" {
		glog.Fatal("-server is required")
	}
	if *deviceFlag == "" {
		glog.Fatal("-device is required")
	}
	if *tokenFlag == "" {
		glog.Fatal("-token is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	glog.Infof("server: %v", *serverFlag)
	conn, err := connect(ctx, *serverFlag, *tokenFlag)
	if err != nil {
		glog.Fatalf("Error dialing server: %s", err)
	}

	wsAPIUtils := ws.NewWsAPIUtils(conn)

	wsID := uuid.NewString()

	_, err = wsAPIUtils.CreateWS(ctx, wsID, "Sample Configlet Workspace",
		"Sample Configlet Workspace Description")
	if err != nil {
		glog.Errorf("Error creating workspace %v", err)
		return
	}
	err = wsAPIUtils.WaitForWSAvailable(ctx, wsID)
	if err != nil {
		glog.Errorf("Error waiting for workspace creation %v", err)
		return
	}

	if *cleanupFlag {
		configlets, assignments, err := getResources(ctx, conn)
		if err != nil {
			return
		}
		logResources(configlets, assignments)
		err = cleanupResources(ctx, wsID, conn, configlets, assignments)
	} else {
		err = createResources(ctx, wsID, *bodyFlag, conn)
	}
	ws, err1 := wsAPIUtils.BuildAndSubmitWS(ctx, wsID)
	if err1 != nil {
		glog.Errorf("Error while building/submitting workspace %v: %v", wsID, err1)
		return
	}
	glog.Infof("Got change control ids %+v", ws.CcIds.GetValues())

	configlets, assignments, err2 := getResources(ctx, conn)
	if err2 != nil {
		return
	}
	logResources(configlets, assignments)
}

func logResources(configlets []*configlet_v1.Configlet,
	assignments []*configlet_v1.ConfigletAssignment) {
	for _, c := range configlets {
		glog.Infof("Got mainline configlet %v %v", c.DisplayName, c.Key)
	}

	glog.Infof("Got %d mainline assignments", len(assignments))
}

func cleanupResources(ctx context.Context, wsID string, conn *grpc.ClientConn,
	configlets []*configlet_v1.Configlet, assignments []*configlet_v1.ConfigletAssignment) error {
	for _, c := range configlets {
		if !strings.Contains(c.DisplayName.String(), "sample configlet") {
			glog.Infof("Skipping configlet %v", c.DisplayName.String())
			continue
		}
		_, err := cApi.SetConfigletConfig(
			ctx, conn,
			&configlet_v1.ConfigletConfig{
				Key: &configlet_v1.ConfigletKey{
					WorkspaceId: toWrapper(wsID),
					ConfigletId: c.Key.ConfigletId,
				},
				Remove: wrapperspb.Bool(true),
			})
		if err != nil {
			glog.Errorf("error setting configlet config for workspace %v configlet %v",
				wsID, c.Key.ConfigletId)
			return err
		}
	}
	for _, a := range assignments {
		if !strings.Contains(a.DisplayName.String(), "sample assignment") {
			glog.Infof("Skipping assignment %v", a.DisplayName.String())
			continue
		}
		glog.Infof("Deleting assignment %v", a.DisplayName.String())
		_, err := cApi.SetConfigletAssignmentConfig(
			ctx, conn,
			&configlet_v1.ConfigletAssignmentConfig{
				Key: &configlet_v1.ConfigletAssignmentKey{
					WorkspaceId:           toWrapper(wsID),
					ConfigletAssignmentId: a.Key.ConfigletAssignmentId,
				},
				Remove: wrapperspb.Bool(true),
			})
		if err != nil {
			glog.Errorf("error setting configlet assignment config for workspace %v assignment %v",
				wsID, a.Key.ConfigletAssignmentId.String())
			return err
		}
	}
	return nil
}

func getResources(ctx context.Context, conn *grpc.ClientConn) (
	[]*configlet_v1.Configlet, []*configlet_v1.ConfigletAssignment, error) {
	configlets, err2 := cApi.GetAllConfiglets(ctx, conn, "", 1)
	if err2 != nil {
		glog.Errorf("Error while getting configlets from mainline %v", err2)
		return nil, nil, nil
	}

	assignments, err3 := cApi.GetAllConfigletAssignments(ctx, conn, "")
	if err3 != nil {
		glog.Errorf("Error while getting assignments from mainline %v", err3)
		return nil, nil, nil
	}
	return configlets, assignments, nil
}

func createResources(ctx context.Context, wsID, body string, conn *grpc.ClientConn) error {
	configletID := uuid.NewString()
	assignmentID := uuid.NewString()
	_, err := cApi.SetConfigletConfig(
		ctx, conn,
		&configlet_v1.ConfigletConfig{
			Key: &configlet_v1.ConfigletKey{
				WorkspaceId: toWrapper(wsID),
				ConfigletId: toWrapper(configletID),
			},
			DisplayName: toWrapper("sample configlet"),
			Body:        toWrapper(body),
		})
	if err != nil {
		glog.Errorf("Error setting configlet config %v", err)
		return err
	}

	_, err = cApi.SetConfigletAssignmentConfig(
		ctx, conn,
		&configlet_v1.ConfigletAssignmentConfig{
			Key: &configlet_v1.ConfigletAssignmentKey{
				WorkspaceId:           toWrapper(wsID),
				ConfigletAssignmentId: toWrapper(assignmentID),
			},
			DisplayName:        toWrapper("sample assignment"),
			Query:              toWrapper(fmt.Sprintf("device:%s", *deviceFlag)),
			ConfigletIds:       &fmp.RepeatedString{Values: []string{configletID}},
			ChildAssignmentIds: nil, // single level tree in this example
			MatchPolicy:        configlet_v1.MatchPolicy_MATCH_POLICY_MATCH_ALL,
		})
	if err != nil {
		glog.Errorf("Error setting assignment config %v", err)
		return err
	}

	inputsConfig := studio_v1.NewInputsConfigServiceClient(conn)
	inputValuesJSON := []byte(fmt.Sprintf(`{"configletAssignmentRoots": ["%s"]}`, assignmentID))
	_, err = inputsConfig.Set(ctx,
		&studio_v1.InputsConfigSetRequest{
			Value: &studio_v1.InputsConfig{
				Key: &studio_v1.InputsKey{
					StudioId:    wrapperspb.String("studio-static-configlet"),
					WorkspaceId: wrapperspb.String(wsID),
					Path:        &fmp.RepeatedString{},
				},
				Inputs: &wrapperspb.StringValue{Value: string(inputValuesJSON)},
			}})
	if err != nil {
		glog.Errorf("Error setting inputs config %v", err)
		return err
	}
	return nil
}

func buildDetails(ctx context.Context, wsID, buildID string, ws ws.WsAPIUtils) {
	allBuildDetails, err := ws.GetAllBuildDetails(ctx, wsID, buildID)
	if err != nil {
		glog.Fatalf("Error while getting build details from mainline %v", err)
		return
	}
	for deviceID, deviceBuildResult := range allBuildDetails {
		glog.Infof("Build details of device %v", deviceID)
		cvr := deviceBuildResult.ConfigValidationResult
		glog.Infof("Diff Summary: %v", cvr.Summary)
		glog.Infof("Config Errors: %v", cvr.Errors)
		glog.Infof("Config Warnings: %v", cvr.Warnings)
		glog.Infof("Config Sources: %+v", cvr.ConfigSources)
	}
}

// TODO: move this to cloudvision-go:grpc/
func connect(ctx context.Context, target, token string) (*grpc.ClientConn, error) {
	tlsConf := cvgrpc.TLSConfig()
	tlsConf.InsecureSkipVerify = true
	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConf)),
		grpc.WithPerRPCCredentials(cvgrpc.NewAccessTokenCredential(token)),
	}
	return grpc.DialContext(ctx, target, dialOpts...)
}
