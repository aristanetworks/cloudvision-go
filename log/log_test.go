// Copyright (c) 2022 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package log

import (
	"context"
	"os"
	"path"
	"testing"

	"golang.org/x/sync/errgroup"
)

func TestLogFile(t *testing.T) {
	errg, ctx := errgroup.WithContext(context.Background())

	if err := InitLogging(path.Join(os.TempDir(), "temporary-test-file.log"), ctx); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		i := i
		errg.Go(func() error {
			Log(errg).Infof("Testing %d", i)
			Log(ctx).Infof("Testing log %d", i)
			return nil
		})
	}

	if err := errg.Wait(); err != nil {
		t.Fatal(err)
	}
}
