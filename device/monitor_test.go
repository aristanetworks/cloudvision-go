// Copyright (c) 2023 Arista Networks, Inc.
// Use of this source code is governed by the Apache License 2.0
// that can be found in the COPYING file.

package device

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
)

func TestDatasourceMonitor(t *testing.T) {
	logger, hook := test.NewNullLogger()

	dm := newDatasourceMonitor(logger.WithField("sensor", "test"))

	// Set Monitor to Info level.
	dm.SetLoggerLevel(logrus.InfoLevel)
	if dm.getLevel() != logrus.InfoLevel {
		t.Fatal("Wrong loglevel!")
	}

	// Monitor should log out Info level message.
	msg1 := "Datasource monitor msg-1"
	dm.Infof(msg1)
	if hook.LastEntry().Message != msg1 {
		t.Fatalf("Monitor last message should be: %s", msg1)
	}

	// Monitor should log out nothing, since current level is lower than Debug level.
	msg2 := "Datasource monitor msg-1"
	dm.Debugf(msg2)
	if hook.LastEntry().Message != msg1 {
		t.Fatalf("Monitor last message should be: %s", msg1)
	}

	// Set Monitor to Debug level.
	dm.SetLoggerLevel(logrus.DebugLevel)
	dm.Debugf(msg2)
	if hook.LastEntry().Message != msg2 {
		t.Fatalf("Monitor last message should be: %s", msg2)
	}

	// Monitor should log out Error level message, since current level is higher.
	msg3 := "Datasource monitor msg-3"
	dm.Errorf(msg3)
	if hook.LastEntry().Message != msg3 {
		t.Fatalf("Monitor last message should be: %s", msg2)
	}
}
