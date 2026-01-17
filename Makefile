.PHONY: build-plugin build-image run stop clean all

IMAGE_NAME := kong-plugin-poc
IMAGE_TAG := latest

build-plugin:  ## No-op for now
	@echo "Plugin build step (no-op)"

build-image: build-plugin
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f kong/Dockerfile .

run: build-image
	@echo "Starting Kong..."
	docker run -d \
		--name kong-dev \
		-p 8000:8000 \
		-p 8001:8001 \
		-p 8443:8443 \
		-p 8444:8444 \
		-e KONG_DATABASE=off \
		-e KONG_DECLARATIVE_CONFIG=/kong.yml \
		-e KONG_ADMIN_LISTEN="0.0.0.0:8001" \
		-v $(PWD)/kong/kong.yml:/kong.yml \
		$(IMAGE_NAME):$(IMAGE_TAG)

logs:
	docker logs kong-dev

stop:
	docker stop kong-dev || true
	docker rm kong-dev || true

clean: stop
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

all: build-plugin build-image
