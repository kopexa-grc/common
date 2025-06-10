GO_TEST = go tool gotest.tools/gotestsum --format pkgname
GO_FUMPT = go tool mvdan.cc/gofumpt
GO_LINT = go tool github.com/golangci/golangci-lint/v2/cmd/golangci-lint

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

	@if ! command -v lefthook >/dev/null 2>&1; then \
		echo "lefthook is not installed. Installing lefthook..."; \
		curl -1sLf https://gobinaries.com/evilmartians/lefthook/install | sh; \
	fi

	@if [ ! -f .lefthook.yml ]; then \
		echo "Installing lefthook configuration..."; \
		lefthook install; \
	fi

prep/tools/mockgen:
	@if ! command -v mockgen >/dev/null 2>&1; then \
		echo "mockgen is not installed. Installing mockgen..."; \
		go install go.uber.org/mock/mockgen@latest; \
	fi

#   🧹 Formatting   #
##@ Formatting

fmt:
	$(GO_FUMPT) -w .

fmt/check:
	$(GO_FUMPT) -d .

#   🔍 Linting     #
##@ Linting

lint:
	$(GO_LINT) run ./...

lint/fix:
	$(GO_LINT) run --fix ./...

#   ⛹🏽‍ License   #
##@ License

license: license/headers/check

license/headers/check:
	copywrite headers --plan

license/headers/apply:
	copywrite headers

#   🧪 Testing     #
##@ Testing

test/ci: test/unit

test/unit:
	mkdir -p build/reports
	$(GO_TEST) --junitfile build/reports/test-unit.xml -- -race ./... -count=1 -short -cover -coverprofile build/reports/unit-test-coverage.out

test/coverage: test/unit
	go tool cover -func=build/reports/unit-test-coverage.out

test/coverage/html: test/unit
	go tool cover -html=build/reports/unit-test-coverage.out -o build/reports/coverage.html

#   🏗️ Building    #
##@ Building

build: license/headers/check
	go build ./...

clean:
	rm -rf build/
	go clean

#   🚀 All-in-One  #
##@ All

all: fmt lint test/unit

#   🔄 Git Hooks   #
##@ Git Hooks

hooks/install:
	lefthook install

hooks/uninstall:
	lefthook uninstall
