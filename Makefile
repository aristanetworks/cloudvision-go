# Copyright (c) 2019 Arista Networks, Inc.
# Use of this source code is governed by the Apache License 2.0
# that can be found in the COPYING file.

GO := go

GOFMT := gofmt
GODIRS := go list ./... | grep -v ".*/smi"
GOFILES := find . -name '*.go' ! -name '*.pb.go' ! -name '*gen.go' ! -name '*.l.go' ! -name '*.y.go'
GOPKGVERSION := $(shell git rev-parse HEAD)
GOLDFLAGS := -ldflags="-s -w -X github.com/aristanetworks/cloudvision-go/version.Version=$(GOPKGVERSION) -X github.com/aristanetworks/cloudvision-go/version.CollectorVersion=$(GOPKGVERSION)"
# Supply defaults if not provided
GOOS ?= linux
GOARCH ?= 386

TEST_TIMEOUT := 60s
GOTEST_FLAGS := -cover -race -count 1
LINT_GOGC := GOGC=50 # Reduce golangci-lint's memory usage
LINT := $(LINT_GOGC) golangci-lint run
LINTFLAGS ?= --deadline=10m --skip-files="device/gen/gen\.go$$" --skip-files="provider/snmp/smi/.*" --exclude-use-default=false --print-issued-lines --print-linter-name --out-format=colored-line-number --disable-all --max-same-issues=0 --max-issues-per-linter=0
LINTCONFIG := --config golangci.yml
LINTNEWCONFIG := --config golangci-new.yml
LINTEXTRAFLAGS ?=
# XXX TODO: We may want to only lint changed packages in the future. For now it's not a big deal to lint everything.
LINT_PKGS := ./...

all: install

install:
	$(GO) install ./...

check: vet fmtcheck test

gen:
	cd device/gen && go generate

smigen:
	cd provider/snmp/smi && golex -o lexer.l.go lexer.l && goyacc -o parser.y.go -v "" parser.y

fmtcheck:
	@if ! which $(GOFMT) >/dev/null; then echo Please install $(GOFMT); exit 1; fi
	goimports=`$(GOFILES) | xargs $(GOFMT) -l 2>&1`; \
	if test -n "$$goimports"; then echo Check the following files for coding style AND USE goimports; echo "$$goimports"; \
		if test "$(shell $(GO) version | awk '{ print $$3 }')" != "devel"; then exit 1; fi; \
	fi
	$(GOFILES) -exec ./check_line_len.awk {} +
	./check_copyright_notice.sh

# Ignore smi package for now.
vet:
	$(GO) vet $(shell $(GODIRS))

lint:
	@# golangci installed from source doesn't support version, so don't fail the target
	-$(LINT) --version
	$(LINT) $(LINTFLAGS) --disable-all --enable=unused --tests=false $(LINTEXTRAFLAGS) $(LINT_PKGS)
	$(LINT) $(LINTCONFIG) $(LINTFLAGS) $(LINTEXTRAFLAGS) $(LINT_PKGS)
	$(LINT) $(LINTNEWCONFIG) $(LINTFLAGS) $(LINTEXTRAFLAGS) $(LINT_PKGS)
#	lint=`$(GODIRS) | xargs -L 1 $(GOLINT) | fgrep -v .pb.go`; if test -n "$$lint"; then echo "$$lint"; exit 1; fi
# The above is ugly, but unfortunately golint doesn't exit 1 when it finds
# lint.  See https://github.com/golang/lint/issues/65

test:
	$(GO) test $(GOTEST_FLAGS) -timeout=$(TEST_TIMEOUT) ./...

.PHONY: all check fmtcheck lint test vet gen smigen
