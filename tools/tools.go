// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

//go:build tools
// +build tools

// Dummy package to declare external build dependencies.

package tools

import (
	// Runtime used by gomock tests / generate commands
	// and mockgen tool used by go:generate
	_ "github.com/golang/mock/mockgen"
	_ "github.com/golang/mock/mockgen/model"
)
