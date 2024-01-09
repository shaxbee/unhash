ifndef _include_go_mk
_include_go_mk := 1

include makefiles/base.mk

### Variables

GOLANGCILINT_VERSION ?= v1.55.2
GOTESTSUM_VERSION ?= v1.11.0
GOCOV_VERSION ?= v1.1.0
GOCOV_HTML_VERSION ?= v1.4.0

GOTESTPKG ?= ./...

# Go coverage directory
# GOCOVERDIR := build/coverage
GOCOVERPKG ?= ./... # Go coverage packages

# JUnit report file
# JUNIT_FILE := build/junit.xml

# Cobertura coverage report file
# CODECOV_FILE := build/coverage.xml

# HTML coverage report file
# CODECOV_HTMLFILE := build/coverage.html

### Targets

.PHONY: generate-go format-go lint-go test-go integration-test-go e2e-test-go coverage-go

generate: generate-go
format: format-go
lint: lint-go
test: test-go
integration-test: integration-test-go
e2e-test: e2e-test-go
coverage: coverage-go

### Tools

# Install golangci-lint
GOLANGCILINT_ROOT := $(BINDIR)/golangci-lint-$(GOLANGCILINT_VERSION)
GOLANGCILINT := $(GOLANGCILINT_ROOT)/golangci-lint

$(GOLANGCILINT):
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOLANGCILINT_ROOT) $(GOLANGCILINT_VERSION)

# Install gotestsum
GOTESTSUM_ROOT := $(BINDIR)/gotestsum-$(GOTESTSUM_VERSION)
GOTESTSUM := $(GOTESTSUM_ROOT)/gotestsum

$(GOTESTSUM):
	GOBIN=$(abspath $(GOTESTSUM_ROOT)) go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)

# Install gocov, gocov-xml
GOCOV_ROOT := $(BINDIR)/gocov-$(GOCOV_VERSION)
GOCOV := $(GOCOV_ROOT)/gocov
GOCOV_XML := $(GOCOV_ROOT)/gocov-xml

$(GOCOV):
	@mkdir -p $(GOCOV_ROOT)
	GOBIN=$(abspath $(GOCOV_ROOT)) go install github.com/axw/gocov/gocov@$(GOCOV_VERSION)

$(GOCOV_XML):
	@mkdir -p $(GOCOV_ROOT)
	GOBIN=$(abspath $(GOCOV_ROOT))  go install github.com/AlekSi/gocov-xml@$(GOCOV_VERSION)

# Install gocov-html
GOCOV_HTML_ROOT := $(BINDIR)/gocov-html-$(GOCOV_HTML_VERSION)
GOCOV_HTML := $(GOCOV_HTML_ROOT)/gocov-html

$(GOCOV_HTML):
	@mkdir -p $(GOCOV_HTML_ROOT)
	GOBIN=$(abspath $(GOCOV_HTML_ROOT)) go install github.com/matm/gocov-html/cmd/gocov-html@$(GOCOV_HTML_VERSION)

### Implementation

generate-go:
	go generate ./...

format-go:
	go fmt ./...

lint-go: $(GOLANGCILINT)
	$(GOLANGCILINT) run

GOTEST := go test -v

test-go:
	$(GOTEST) -short $(GOTESTPKG) $(GOTESTARGS)

integration-test-go:
	$(GOTEST) -count 1 $(GOTESTPKG) $(GOTESTARGS)

e2e-test-go:
	$(GOTEST) -count 1 ./e2e $(GOTESTARGS)

# if JUNIT_FILE is set generate JUnit reports
ifneq ($(strip $(JUNIT_FILE)),)
test-go integration-test-go e2e-test-go: $(GOTESTSUM)
test integration-test e2e-test: $(JUNIT_FILE)

GOTESTOUT := $(TMPDIR)/test-results.json
GOTEST := $(GOTESTSUM) --format standard-verbose --jsonfile $(GOTESTOUT) --

$(JUNIT_FILE): $(GOTESTSUM)
	@mkdir -p $(dir $(JUNIT_FILE))
	$(GOTESTSUM) --junitfile $(JUNIT_FILE) --raw-command cat $(GOTESTOUT) &>/dev/null
endif # JUNIT_FILE

ifneq ($(filter coverage,$(MAKECMDGOALS)),)
GOCOVERDIR ?= $(TMPDIR)/coverage
GOCOVEROUT ?= $(TMPDIR)/go-cover.out
endif

ifneq ($(strip $(GOCOVERDIR)),)
GOTEST += -coverpkg=$(GOCOVERPKG) -covermode=atomic
GOTESTARGS += -test.gocoverdir=$(abspath $(GOCOVERDIR))

test-go integration-test-go e2e-test-go: $(GOCOVERDIR)

coverage-go: $(GOCOVEROUT)
	go tool covdata func -i $(abspath $(GOCOVERDIR)) -pkg $(GOCOVERPKG)

$(GOCOVEROUT): $(GOCOVERDIR)
	go tool covdata textfmt -i $(abspath $(GOCOVERDIR)) -o $(GOCOVEROUT) -pkg $(GOCOVERPKG)

$(GOCOVERDIR):
	@mkdir -p $(GOCOVERDIR)
endif # GOCOVERDIR

# if CODECOV_FILE is set generate Cobertura coverage report
ifneq ($(strip $(CODECOV_FILE)),)
.PHONY: $(CODECOV_FILE)
coverage: $(CODECOV_FILE)
$(CODECOV_FILE): $(GOCOV) $(GOCOV_XML) $(GOCOVEROUT)
	@mkdir -p $(dir $(CODECOV_FILE))
	$(GOCOV) convert $(GOCOVEROUT) | $(GOCOV_XML) >$(CODECOV_FILE)
endif # CODECOV_FILE

# if CODECOV_HTMLFILE is set generate HTML coverage report
ifneq ($(strip $(CODECOV_HTMLFILE)),)
.PHONY: $(CODECOV_HTMLFILE)
coverage: $(CODECOV_HTMLFILE)
$(CODECOV_HTMLFILE): $(GOCOV) $(GOCOV_HTML) $(GOCOVEROUT)
	@mkdir -p $(dir $(CODECOV_HTMLFILE))
	$(GOCOV) convert $(GOCOVEROUT) | $(GOCOV_HTML) -t kit >$(CODECOV_HTMLFILE)
endif # CODECOV_HTMLFILE

endif # _include_go_mk