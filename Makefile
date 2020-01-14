stop:
	@echo "stopping containers..."
	-@docker stop $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing containers..."
	-@docker rm $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing image..."
	-@docker rmi networker

image:
	@echo "building image..."
	-@docker build --no-cache -t networker .

start: stop image
	@echo "starting container..."
	-@docker run -it --rm --network host networker

lint:
	@goimports -w $(shell find . -type f -name *.go)
	@echo "go files have been linted"

append_commit: lint
	@git add .
	@git commit --amend --no-edit
	@echo "appended commit"