// Copyright (c) 2021 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

// Package redirector provides the functionality to redirect to region specific
// apiserver, based on the request context. This package and package gen which
// contains the generated protobuf used for communication are currently internal
// and are subject to change.
package redirector

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/aristanetworks/cloudvision-go/internal/redirector/gen"
)

// Resolve will return a list of targets that the client can connect to the
// instance of cloudvision it is registered in. In case the functionality is not
// implemented serverside, we return an empty list of targets.
func Resolve(ctx context.Context, ep string, opts ...grpc.DialOption) ([]string, error) {
	conn, err := grpc.DialContext(ctx, ep, opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to dial %s for regional redirection: %w", ep, err)
	}

	// systemID is empty as we want to return the default assignment for the
	// organization bearing the token.
	resp, err := pb.NewRedirectorClient(conn).GetAssignment(ctx, &pb.GetAssignmentRequest{})
	if err != nil {
		// Handle a special case for unimplemented and not found, which means
		// that the functionality is not yet available serverside.
		if c := status.Code(err); c == codes.Unimplemented || c == codes.NotFound {
			return []string{}, nil
		}
		return nil, fmt.Errorf("unable to get list of hosts: %w", err)
	}
	return strings.Split(resp.GetAssignment().GetClusters(), ","), nil
}
