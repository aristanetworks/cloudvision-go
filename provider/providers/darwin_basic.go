// Copyright (c) 2018 Arista Networks, Inc.  All rights reserved.
// Arista Networks, Inc. Confidential and Proprietary.
// Subject to Arista Networks, Inc.'s EULA.
// FOR INTERNAL USE ONLY. NOT FOR DISTRIBUTION.

package providers

import (
	"arista/entity"
	"arista/provider"
	"arista/schema"
	"arista/types"
	"fmt"
	"os"
	"time"

	"github.com/aristanetworks/goarista/path"
)

type darwinProvider struct {
	provider.ReadOnly
	// Closed when we're done initialization
	ready chan struct{}
	// Closed when we want to stop Run()
	done chan struct{}

	// Sampling period
	period time.Duration
}

func (p *darwinProvider) WaitForNotification() {
	<-p.ready
}

func (p *darwinProvider) Stop() {
	<-p.ready
	close(p.done)
}

func (p *darwinProvider) updateStats() {
}

func setSystemConfig(root types.Entity) types.Entity {
	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	systemType := types.NewEntityType("::system::config")
	systemType.AddAttribute("hostname", types.StringType)
	systemType.AddAttribute("domain-name", types.StringType)
	systemType.AddAttribute("login-banner", types.StringType)
	systemType.AddAttribute("motd-banner", types.StringType)
	data := map[string]interface{}{"hostname": hostname}
	systemConfig, err := entity.MakeDirsWithAttributes(root,
		path.New("system", "config"), nil, systemType, data)
	if err != nil {
		panic(fmt.Errorf("Failed to create /system/config: %s", err))
	}
	return systemConfig
}

func (p *darwinProvider) Run(s *schema.Schema, root types.Entity, ch chan<- types.Notification) {
	setSystemConfig(root)
	close(p.ready)
	tick := time.NewTicker(p.period)
	defer tick.Stop()
	for {
		select {
		case <-tick.C:
			p.updateStats()
		case <-p.done:
			return
		}
	}
}

// NewDarwinProvider returns a read-only basic darwin provider that pushes data
// following the OpenConfig convention
func NewDarwinProvider() provider.Provider {
	return &darwinProvider{
		ready:  make(chan struct{}),
		done:   make(chan struct{}),
		period: time.Second * 10,
	}
}
