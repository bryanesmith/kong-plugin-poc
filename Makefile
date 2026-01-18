.PHONY: build-plugin build-image run stop clean all

IMAGE_NAME := kong-plugin-poc
IMAGE_TAG := latest

build-mcp:
	@echo "Building MCP server..."
	cd mcp_server && go mod tidy
	cd mcp_server && go build -o mcp_server main.go

run-mcp: build-mcp
	@echo "Starting MCP server..."
	cd mcp_server && ./mcp_server

test-mcp: build-mcp
	@echo "Testing MCP server..."
	cd mcp_server && ./test_request.sh | ./mcp_server

stop-mcp:
	@echo "Stopping MCP server..."
	killall mcp_server || true

clean-mcp:
	@echo "Cleaning MCP server..."
	cd mcp_server && rm -f mcp_server || true

build-plugin:  ## No-op for now
	@echo "Plugin build step (no-op)"

build-kong: build-plugin
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f kong/Dockerfile .

run-kong: build-image
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

stop-kong:
	docker stop kong-dev || true
	docker rm kong-dev || true

logs-kong:
	docker logs kong-dev

clean-kong: stop-kong
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

#all: build-plugin build-image
