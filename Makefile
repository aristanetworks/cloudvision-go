GO := go

TEST_TIMEOUT := 30s
GOTEST_FLAGS := -v
GOLINT := golint
GOFMT := gofmt
GODIRS := find . -type d ! -path './.git/*' ! -path './vendor/*'
GOFILES := find . -name '*.go' ! -path './vendor/*' ! -name '*.pb.go'

DEVICE_DIR := ./device
DEVICE_INTEGRATION_DIR := $(DEVICE_DIR)/integration
PROVIDER_DIR := ./provider
DEVICE_PLUGIN_DIR := $(DEVICE_DIR)/plugins
DEVICE_TEST_DIR := $(DEVICE_DIR)/test
PLUGIN_SO_BASE_DIR := $(DEVICE_PLUGIN_DIR)/build
# To test loading plugins recursively.
PLUGIN_SO_DIR := $(DEVICE_PLUGIN_DIR)/build/test1/test2
PLUGIN_NAME := test.so
PLUGIN_SRC_FILES := $(DEVICE_PLUGIN_DIR)/test_plugin.go
DEVICE_INTEGRATION_FLAGS = -test.short -cover -tags="integration device"

all: install

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

plugins:
	mkdir -p $(PLUGIN_SO_DIR)
	# Build a dummy plugin for integration testing.
	$(GO) build $(GOLDFLAGS) -buildmode=plugin -o $(PLUGIN_SO_DIR)/$(PLUGIN_NAME) $(PLUGIN_SRC_FILES)

clean-plugins:
	rm -rf $(PLUGIN_SO_BASE_DIR)

integration:
	$(MAKE) plugins
	$(GO) test $(DEVICE_INTEGRATION_DIR)/... $(DEVICE_INTEGRATION_FLAGS)
	$(MAKE) clean-plugins

COVER_PKGS := `find . -name '*_test.go' ! -path "./.git/*" ! -path "./vendor/*" | xargs -I{} dirname {} | sort -u`

.PHONY: all check fmtcheck jenkins lint test vet plugins clean-plugins
