stop:
	@echo "stopping containers..."
	-@docker stop $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing containers..."
	-@docker rm $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing image..."
	-@docker rmi networker

binary:
	@echo "building binary..." && \
	go build -o networker cmd/networker/*.go

windows_binary:
	@echo "building binary for windows" &&\
	GOOS=windows GOARCH=386 go build -o networker.exe cmd/networker/*.go

release: binary windows_binary
	@echo "release binaries ready"

image:
	@echo "building image..."
	-@docker build --no-cache -t networker .

start: stop image
	@echo "starting container..."
	-@docker run -it --rm --network host networker

lint:
	@echo "linting all go files..."  &&\
	goimports -w $(shell find . -name "*.go")

append_commit: lint
	@git add .
	@git commit --amend --no-edit
	@echo "appended commit"

remove_binaries:
	@echo "removing any old binaries if they exist..."
	-@rm networker networker.exe

reset: remove_binary binary
	clear