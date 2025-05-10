GO_TOOLS := go run -modfile ./tools/go.mod
GOLANCI_LINT = $(GO_TOOLS) github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2
GOFUMPT := $(GO_TOOLS) mvdan.cc/gofumpt
GO_TEST = $(GO_TOOLS) gotest.tools/gotestsum --format pkgname

#   🔨 TOOLS       #
##@ Tools

prep: prep/tools

prep/tools:
	@if ! command -v golangci-lint >/dev/null 2>&1; then \
		echo "golangci-lint is not installed. Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest; \
	fi

	@if ! command -v copywrite >/dev/null 2>&1; then \
		echo "copywrite is not installed. Installing copywrite..."; \
		go install github.com/hashicorp/copywrite@latest; \
	fi

#   ⛹🏽‍ License   #


license: license/headers/check

license/headers/check:
	copywrite headers --plan

license/headers/apply:
	copywrite headers

test/ci: test/unit

test/unit:
	mkdir -p build/reports
	$(GO_TEST) --junitfile build/reports/test-unit.xml -- -race ./... -count=1 -short -cover -coverprofile build/reports/unit-test-coverage.out
