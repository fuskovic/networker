.PHONY: clean
clean:
	-rm ./tls/*.pem

PHONY: cert
cert: clean
	@./generate_self_signed_cert.sh

.PHONY: image
image: 
	-docker rmi networker_playground_server --force
	docker build -t networker_playground_server --no-cache .

.PHONY: container
container: cert image
	-docker rm "$(shell docker ps -a | grep nwserver | awk '{ print $$1 }')" --force
	docker run -d \
		--name nwserver \
		-e NETWORKER_PLAYGROUND_PORT=80 \
		-p 80:80 \
		-v $(shell git rev-parse --show-toplevel)/playground/tls:/app/tls/ \
		networker_playground_server

.PHONY: deployment
deployment:
	kubectl create deployment \
		--dry-run=client \
		-o yaml \
		--image networker_playground_server \
		networker_servers > deployment.dev.yaml

.PHONY:servers
servers: cert image
	-docker volume rm playground_tls
	-docker-compose down
	docker-compose up

.PHONY:balance
balance:
	networker balance \
		--targets 127.0.0.1:3001,127.0.0.1:3002,127.0.0.1:3003 \
		--tls \
		--key ${PWD}/tls/key.pem \
		--cert ${PWD}/tls/cert.pem \
		--port 80

