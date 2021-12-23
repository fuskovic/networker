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

.PHONY: update_badge
update_badge:
	@gopherbadger -md="README.md" -png=false

.PHONY:fmt
fmt:
	@goimports -w $(shell find . -name "*.go") && echo "go files formatted"

.PHONY:commit
commit: fmt
	@git add . && git commit --amend --no-edit

.PHONY: readme
readme:
	@./scripts/update_readme.sh