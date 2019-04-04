# Copyright (c) 2019 Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the COPYING file.

gendump: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO) build $(GOLDFLAGS) -o gendump-$(GOPKGVERSION)

include ../../../../Makefile

clean:
	rm -f gendump-*

.PHONY: gendump clean

