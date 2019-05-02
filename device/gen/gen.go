// Copyright (c) 2019 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package gen

// To run this command you need protoc:
// brew install protobuf

//go:generate go get -u github.com/golang/protobuf/protoc-gen-go github.com/mwitkow/go-proto-validators/protoc-gen-govalidators
//go:generate protoc --proto_path=. --proto_path=${GOPATH}/src --govalidators_out=. --proto_path=${GOPATH}/src/github.com/google/protobuf/src --go_out=plugins=grpc:. inventory.proto
