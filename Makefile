.PHONY:clean
clean:
	@./scripts/clean.sh

.PHONY:build
build: clean
	@./scripts/build.sh

.PHONY:install
install:
	@./scripts/install.sh

.PHONY: test
test:
	@go clean -testcache && go test -v ./...

.PHONY: coverage
coverage:
	@./scripts/generate_test_coverage_report.sh $(mode)

.PHONY: badge
badge:
	@gopherbadger -md="README.md" -png=false

.PHONY:fmt
fmt:
	@goimports -w $(shell find . -name "*.go") && echo "go files formatted"

.PHONY:commit
commit: fmt
	@git add . && git commit --amend --no-edit

.PHONY: docs
docs:
	@go run ./scripts/doc_gen.go

.PHONY: version
version:
	@./scripts/get_version.sh

.PHONY: module_refresh
module_refresh:
	@GOPROXY=proxy.golang.org go list -m github.com/fuskovic/networker/v2@latest