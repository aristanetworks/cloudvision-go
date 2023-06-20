// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package grpc

import (
	"context"
	"errors"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/aristanetworks/cloudvision-go/api/arista/redirector.v1"
	"github.com/aristanetworks/cloudvision-go/api/fmp"
	"github.com/aristanetworks/cloudvision-go/internal/redirector"
)

func TestResolveRedirection(t *testing.T) {
	testCases := map[string]struct {
		env     string
		in      string
		targets [][]string
		want    string
		err     error
		wantErr error
	}{
		"no env": {
			env: "",
			in:  "arista.io",
			targets: [][]string{
				{"host1.region1.arista.io", "host2.region1.arista.io"},
				{"host1.region2.arista.io", "host2.region2.arista.io"},
			},
			want:    "host1.region1.arista.io",
			err:     nil,
			wantErr: nil,
		},
		"disabled": {
			env: "false",
			in:  "arista.io",
			targets: [][]string{
				{"host1.region1.arista.io", "host2.region1.arista.io"},
				{"host1.region2.arista.io", "host2.region2.arista.io"},
			},
			want:    "arista.io",
			err:     nil,
			wantErr: nil,
		},
		"any error": {
			env: "true",
			in:  "arista.io",
			targets: [][]string{
				{"host1.region1.arista.io", "host2.region1.arista.io"},
				{"host1.region2.arista.io", "host2.region2.arista.io"},
			},
			want:    "",
			err:     io.EOF,
			wantErr: io.EOF,
		},
		"no targets": {
			env:     "true",
			in:      "arista.io",
			targets: [][]string{},
			want:    "",
			err:     ErrNoClusterTargets,
			wantErr: ErrNoClusterTargets,
		},
		"unimplemented": {
			env:     "true",
			in:      "arista.io",
			targets: [][]string{},
			want:    "arista.io",
			err:     redirector.ErrUnimplemented,
			wantErr: nil,
		},
		"good": {
			env: "true",
			in:  "arista.io",
			targets: [][]string{
				{"host1.region1.arista.io", "host1.region1.arista.io"},
				{"host1.region2.arista.io", "host2.region2.arista.io"},
			},
			want:    "host1.region1.arista.io",
			err:     nil,
			wantErr: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			os.Setenv("CLOUDVISION_REGIONAL_REDIRECT", tc.env)
			defer os.Setenv("CLOUDVISION_REGIONAL_REDIRECT", "")

			if got, err := resolveRedirection(
				context.Background(),
				tc.in,
				grpc.WithContextDialer(dialer(tc.targets, tc.err)),
				// on a real server we use TLS for caller resolution.
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			); got != tc.want || !errorIs(err, tc.wantErr) {
				t.Errorf("resolveRedirection() = %q, %v; want %q, %v",
					got, err, tc.want, tc.err)
			}
		})
	}
}

// The code below is copied from internal/redirector/redirector_test.go as we
// want to test the same behaviour, but this time how well it integrates.

func dialer(want [][]string, err error) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)

	server := grpc.NewServer()
	pb.RegisterAssignmentServiceServer(server, &assignmentService{
		want: want,
		err:  err,
	})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}

type assignmentService struct {
	pb.UnimplementedAssignmentServiceServer

	want [][]string
	err  error
}

func (s *assignmentService) GetOne(
	ctx context.Context, req *pb.AssignmentRequest) (*pb.AssignmentResponse, error) {
	if s.err != nil {
		// Special unimplemented case
		if s.err == redirector.ErrUnimplemented {
			return nil, status.Error(codes.Unimplemented, "unimplemented")
		}
		return nil, s.err
	}

	clusters := make([]*pb.Cluster, 0, len(s.want))
	for _, targets := range s.want {
		clusters = append(clusters, &pb.Cluster{
			Hosts: &fmp.RepeatedString{
				Values: targets,
			},
		})
	}

	return &pb.AssignmentResponse{
		Value: &pb.Assignment{
			Clusters: &pb.Clusters{
				Values: clusters,
			},
		},
	}, nil
}

// errorsIs is adapts the stdlib errors.Is and also allows to compare strings
// for handling gRPC errors
func errorIs(a, b error) bool {
	if errors.Is(a, b) {
		return true
	}

	if (a == nil && b != nil) || (a != nil && b == nil) {
		return false
	}

	return strings.Contains(a.Error(), b.Error())
}
