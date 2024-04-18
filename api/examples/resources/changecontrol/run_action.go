// Copyright (c) 2024 Arista Networks, Inc.  All rights reserved.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.
//
// Create Change Control for specific action(s) and monitor for completion
// Example usage:
//  CLOUDVISION_REGIONAL_REDIRECT=false go run run_action.go --server www.cv-dev.corp.arista.io:443
//		--token-file temp-test/token.txt --action-args temp-test/actions-args.json
//
// This is a go conversion of the python script found at:
// github.com/aristanetworks/cloudvision-python/blob/trunk
//		examples/resources/changecontrol/run_action.py

package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"os"
	"time"

	changecontrol_v1 "github.com/aristanetworks/cloudvision-go/api/arista/changecontrol.v1"
	fmp_wrappers "github.com/aristanetworks/cloudvision-go/api/fmp"
	timestamp "google.golang.org/protobuf/types/known/timestamppb"
	wrappers "google.golang.org/protobuf/types/known/wrapperspb"

	cv_grpc "github.com/aristanetworks/cloudvision-go/grpc"
	"github.com/aristanetworks/glog"
)

var (
	RPC_TIMEOUT time.Duration = 30 // in seconds
	// Declare args values, bound in init()
	server     = ""
	actionArgs = ""
	tokenFile  = ""
	certFile   = ""
)

func createConnection(ctx context.Context, server string, tokenFile string,
	certFile string) (*grpc.ClientConn, error) {
	/*
	   Creates a grpc secure connection with the provided info

	   Args:
	       server: 	  Address of the CV instance to create the connection to
	       tokenFile: Path to a file containing a token
	       certFile:  Path to a file containing the cert (optional)

	   Returns:
	       grpc.ClientConn: ClientConn that rAPI stubs can use
	*/
	glog.Infof("Creating Client Connection...")
	// Create Auth using optional cert file and token
	auth, err := cv_grpc.NewTokenAuth(tokenFile, certFile)
	if err != nil {
		glog.Errorf("Failed to generate authentication from token and cert file paths: %v", err)
		return nil, err
	}
	timeoutCtx, cancel := context.WithTimeout(ctx, RPC_TIMEOUT*time.Second)
	defer cancel()
	connCreds, err := cv_grpc.DialWithAuth(timeoutCtx, server, auth)
	if err != nil {
		glog.Errorf("Failed to create grpc connection: %v", err)
		return nil, err
	}
	glog.Infof("Client Connection created!")
	return connCreds, nil
}

func loadArgs(path string) map[string]map[string]string {
	/*
	   Loads the actions and args from a file

	   Args:
	       path: 	path to a json file containing the actions and args
	            	for the change control, provided as a flag at run time

	   Returns:
	       actionAndArgs:	A map mimicing a python
	                    	Dict[string, Dict[string, string]] that is loaded
	                      	with all the action IDs and matching args

	*/
	glog.Infof("Loading args from file..")
	content, err := os.ReadFile(path)
	if err != nil {
		glog.Fatalf("Failed to read actions and args file: %v", err)
	}

	var actionsAndArgs map[string]map[string]string
	err = json.Unmarshal(content, &actionsAndArgs)
	if err != nil {
		glog.Fatalf("Failed to unmarshal actions and args: %v", err)
	}
	glog.Infof("Args and actions loaded!")
	return actionsAndArgs

}

func addCC(ctx context.Context, connection *grpc.ClientConn, ccID string,
	actionsAndArgs map[string]map[string]string) (*timestamp.Timestamp, error) {
	/*
		   Creates a Change Control of the given ID and returns the Timestamp for approval

		   Args:
		       connection:		The GRPC connection that can be used by rAPI stubs
		       ccID:			The ID of the Change Control to create
		       actionsAndArgs:	A string map of the actions and their respective arguments.
		                    	The top level keys should be for each action in the order to be
								executed, and each of the contents of those entries a str->str
								dictionary of the arguments to be passed.

		   Returns:
		       Timestamp: 		Timestamp of the write that can be used to approve the CC
	*/

	glog.Infof("Creating Change Control with ID %s", ccID)
	ccName := "run_action script created change"
	rootStageId := "stage-root"
	var rootStageRows = []*fmp_wrappers.RepeatedString{}
	var stageConfigMapDict = make(map[string]*changecontrol_v1.StageConfig)
	for actionID, args := range actionsAndArgs {
		currActionID := "stage-action " + actionID
		action := changecontrol_v1.Action{
			Name: &wrappers.StringValue{Value: actionID},
			Args: &fmp_wrappers.MapStringString{Values: args},
		}
		rootStageRows = append(rootStageRows,
			&fmp_wrappers.RepeatedString{Values: []string{currActionID}})
		stageConfigMapDict[currActionID] = &changecontrol_v1.StageConfig{
			Name:   &wrappers.StringValue{Value: ("Scheduled action " + actionID)},
			Action: &action,
		}
	}

	stageConfigMapDict[rootStageId] = &changecontrol_v1.StageConfig{
		Name: &wrappers.StringValue{Value: (ccName + " Root")},
		Rows: &changecontrol_v1.RepeatedRepeatedString{Values: rootStageRows},
	}
	stageConfigMap := changecontrol_v1.StageConfigMap{
		Values: stageConfigMapDict,
	}
	changeConfig := changecontrol_v1.ChangeConfig{
		Name:        &wrappers.StringValue{Value: ccName},
		RootStageId: &wrappers.StringValue{Value: rootStageId},
		Stages:      &stageConfigMap,
		Notes:       &wrappers.StringValue{Value: "Created and managed by script"},
	}
	key := changecontrol_v1.ChangeControlKey{Id: &wrappers.StringValue{Value: ccID}}
	setReq := changecontrol_v1.ChangeControlConfigSetRequest{
		Value: &changecontrol_v1.ChangeControlConfig{
			Key:    &key,
			Change: &changeConfig,
		},
	}

	ccServClient := changecontrol_v1.NewChangeControlConfigServiceClient(connection)
	timeoutCtx, cancel := context.WithTimeout(ctx, RPC_TIMEOUT*time.Second)
	defer cancel()
	resp, err := ccServClient.Set(timeoutCtx, &setReq)
	if err != nil {
		glog.Errorf("Failed to create change control %s with error: %v", ccID, err)
		return nil, err
	}
	glog.Infof("Change Control %s created successfully", ccID)
	return resp.Time, nil
}

func approveCC(ctx context.Context, connection *grpc.ClientConn, ccID string,
	ts *timestamp.Timestamp) error {
	/*
		   Approves a Change Control of the given ID and Timestamp

		   Args:
		       connection (grpc.ClientConn):	The GRPC connection that can be
			   									used by rAPI stubs
		       ccID (string):	The ID of the Change Control to approve
		       ts (Timestamp):	The Timestamp of the Change Control to approve
	*/

	glog.Infof("Approving Change Control %s", ccID)
	key := changecontrol_v1.ChangeControlKey{Id: &wrappers.StringValue{Value: ccID}}
	setReq := changecontrol_v1.ApproveConfigSetRequest{
		Value: &changecontrol_v1.ApproveConfig{
			Key: &key,
			Approve: &changecontrol_v1.FlagConfig{
				Value: &wrappers.BoolValue{Value: true},
			},
			// NOTE: TS needs to match that of the cc update in the DB
			Version: ts,
		},
	}

	ccAprvClient := changecontrol_v1.NewApproveConfigServiceClient(connection)
	timeoutCtx, cancel := context.WithTimeout(ctx, RPC_TIMEOUT*time.Second)
	defer cancel()
	_, err := ccAprvClient.Set(timeoutCtx, &setReq)
	if err != nil {
		glog.Errorf("Failed to approve change control %s with error: %v", ccID, err)
		return err
	}
	glog.Infof("Change Control %s approved successfully", ccID)
	return nil

}

func executeCC(ctx context.Context, connection *grpc.ClientConn, ccID string) error {
	/*
		   Executes and approved Change Control of the given ID

		   Args:
		       connection (grpc.ClientConn):	The GRPC connection that can be
			   									used by rAPI stubs
		       ccID (string):	The ID of the Change Control to approve
	*/

	glog.Infof("Executing Change Control %s", ccID)
	key := changecontrol_v1.ChangeControlKey{Id: &wrappers.StringValue{Value: ccID}}
	setReq := changecontrol_v1.ChangeControlConfigSetRequest{
		Value: &changecontrol_v1.ChangeControlConfig{
			Key: &key,
			Start: &changecontrol_v1.FlagConfig{
				Value: &wrappers.BoolValue{Value: true},
			},
		},
	}
	ccExecClient := changecontrol_v1.NewChangeControlConfigServiceClient(connection)
	timeoutCtx, cancel := context.WithTimeout(ctx, RPC_TIMEOUT*time.Second)
	defer cancel()
	_, err := ccExecClient.Set(timeoutCtx, &setReq)
	if err != nil {
		glog.Errorf("Failed to execute change control %s with error: %v", ccID, err)
		return err
	}
	glog.Infof("Change Control %s executed successfully", ccID)
	return nil
}

func subscribeToCCStatus(ctx context.Context, connection *grpc.ClientConn, ccID string) error {
	/*
		   Subscribes to a Change Control and monitors it until completion

		   Args:
		       connection (grpc.ClientConn):	The GRPC connection that can be used
			   									by rAPI stubs
		       ccID (string):	The ID of the Change Control to approve
	*/

	glog.Infof("Subscribing to change control %s to monitor for completion", ccID)
	key := changecontrol_v1.ChangeControlKey{Id: &wrappers.StringValue{Value: ccID}}
	subReq := changecontrol_v1.ChangeControlStreamRequest{}
	subReq.PartialEqFilter = append(subReq.PartialEqFilter,
		&changecontrol_v1.ChangeControl{Key: &key})

	ccServiceClient := changecontrol_v1.NewChangeControlServiceClient(connection)
	timeoutCtx, cancel := context.WithTimeout(ctx, RPC_TIMEOUT*time.Second)
	defer cancel()
	responses, err := ccServiceClient.Subscribe(timeoutCtx, &subReq)
	if err != nil {
		glog.Errorf("Failed to get subscribtion responses: %v", err)
		return err
	}
	// The program will stay in this loop either until the change control completes
	// or an error occurs with fetching the most recent response
	for {
		resp, err := responses.Recv()
		if err != nil {
			glog.Errorf("Failed to recieve from subscription: %v", err)
			return err
		}
		if resp.Value != nil {
			if resp.Value.Status ==
				changecontrol_v1.ChangeControlStatus_CHANGE_CONTROL_STATUS_COMPLETED {
				if resp.Value.Error != nil && resp.Value.Error.Value != "" {
					err := resp.Value.Error.Value
					glog.Infof("Change Control %s completed with error: %v",
						ccID, err)
					return errors.New(err)
				}
				glog.Infof("Change Control %s completed successfully", ccID)
				return nil
			}
		}
	}
}

func main() {
	flag.Parse()
	ctx := context.Background()
	actionsAndArgs := loadArgs(actionArgs)

	connection, err := createConnection(ctx, server, tokenFile, certFile)
	if err != nil {
		glog.Fatalf("Failed to create connection: %v", err)
	}
	defer connection.Close()

	ccID := uuid.New().String()

	ts, err := addCC(ctx, connection, ccID, actionsAndArgs)
	if err != nil {
		glog.Fatalf("Failed to add change control: %v", err)
	}

	err = approveCC(ctx, connection, ccID, ts)
	if err != nil {
		glog.Fatalf("Failed to approve change control: %v", err)
	}

	err = executeCC(ctx, connection, ccID)
	if err != nil {
		glog.Fatalf("Failed to execute change control: %v", err)
	}

	err = subscribeToCCStatus(ctx, connection, ccID)
	if err != nil {
		glog.Fatalf("Failed to subscribe to change control: %v", err)
	}
}

func init() {
	flag.StringVar(&server, "server", "", "CloudVision server to "+
		"connect to in <host>:<port> format")
	help := ("path to json file of the actions and arguments to be run, e.g. actionsAndArgs.json." +
		" Top level keys should be the action IDs, with each actionID entry containing" +
		" the string arguments for that action. Actions will be executed serially in" +
		" the order defined. The same ActionID may not be executed more than once.")
	flag.StringVar(&actionArgs, "action-args", "", help)
	flag.StringVar(&tokenFile, "token-file", "", "file with access token")
	flag.StringVar(&certFile, "cert-file", "", "(optional) certificate to use as root CA")
}
