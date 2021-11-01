.PHONY:clean
clean:
	@./scripts/clean.sh

.PHONY:build
build: clean
	@./scripts/build.sh

.PHONY:install
install:
	@./scripts/install.sh

.PHONY:image
image:
	@echo "building image" && docker build --no-cache -t networker .

.PHONY:container
container: clean image
	@echo "starting container" && docker run -it --rm --network host networker

.PHONY:fmt
fmt:
	@goimports -w $(shell find . -name "*.go") && echo "go files formatted"

.PHONY:commit
commit: fmt
	@git add . && git commit --amend --no-edit