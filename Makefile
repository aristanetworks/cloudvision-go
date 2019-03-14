# Copyright (c) 2019 Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the COPYING file.

GO := go

TEST_TIMEOUT := 60s
GOTEST_FLAGS := -cover -race -count 1
GOLINT := golint
GOFMT := gofmt
GODIRS := find . -type d ! -path './.git/*' ! -path './vendor/*'
GOFILES := find . -name '*.go' ! -path './vendor/*' ! -name '*.pb.go'
GOPKGVERSION := $(shell git rev-parse HEAD)
GOLDFLAGS := -ldflags="-s -w -X arista/version.Version=$(GOPKGVERSION)"
# Supply defaults if not provided
GOOS ?= linux
GOARCH ?= 386

all: install

# XXX TODO: This is a change to test the cloudvision-go-to-github Jenkins job.
# It should be removed ASAP.

install:
	$(GO) install ./...

check: vet fmtcheck lint test

jenkins: check

fmtcheck:
	@if ! which $(GOFMT) >/dev/null; then echo Please install $(GOFMT); exit 1; fi
	goimports=`$(GOFILES) | xargs $(GOFMT) -l 2>&1`; \
	if test -n "$$goimports"; then echo Check the following files for coding style AND USE goimports; echo "$$goimports"; \
		if test "$(shell $(GO) version | awk '{ print $$3 }')" != "devel"; then exit 1; fi; \
	fi
	$(GOFILES) -exec ./check_line_len.awk {} +

vet:
	$(GO) vet ./...

lint:
	lint=`$(GODIRS) | xargs -L 1 $(GOLINT) | fgrep -v .pb.go`; if test -n "$$lint"; then echo "$$lint"; exit 1; fi
# The above is ugly, but unfortunately golint doesn't exit 1 when it finds
# lint.  See https://github.com/golang/lint/issues/65

test:
	$(GO) test $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) ./...

.PHONY: all check fmtcheck jenkins lint test vet
