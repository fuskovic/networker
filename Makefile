stop:
	@echo "stopping containers..."
	-@docker stop $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing containers..."
	-@docker rm $(docker ps -aq | grep networker | awk '{ print $1 }')
	@echo "removing image..."
	-@docker rmi networker

image:
	@echo "building image..."
	-@docker build -t networker .

start: stop image
	@echo "starting container..."
	-@docker run -it --rm --network host networker