// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package redirector

import (
	"context"
	"errors"
	"fmt"

	pb "github.com/aristanetworks/cloudvision-go/api/arista/redirector.v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

var (
	// ErrUnimplemented is returned if the redirector is not implemented on the
	// server. This is a dedicated error as status.Status does not support error
	// unwrapping. https://github.com/grpc/grpc-go/issues/2934
	ErrUnimplemented = errors.New("unimplemented")
)

// Resolve will return a slice of clusters, where each cluster is a slice of
// targets that the client can connect to.
func Resolve(ctx context.Context, ep string, opts ...grpc.DialOption) ([][]string, error) {
	conn, err := grpc.DialContext(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to dial %s for redirection information: %w", ep, err)
	}

	resp, err := pb.NewAssignmentServiceClient(conn).GetOne(ctx, &pb.AssignmentRequest{
		// Query with empty key because we use the cert for system_id resolution.
		Key: &pb.AssignmentKey{
			SystemId: wrapperspb.String(""),
		},
	})
	if err != nil {
		if status.Code(err) == codes.Unimplemented {
			err = ErrUnimplemented
		}
		return nil, fmt.Errorf("unable to get redirection assignments: %w", err)
	}

	clusters := resp.GetValue().GetClusters().GetValues()
	res := make([][]string, 0, len(clusters))

	for _, cluster := range clusters {
		res = append(res, cluster.GetHosts().GetValues())
	}

	return res, nil
}
