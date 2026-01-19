.PHONY: build-plugin build-image run stop clean all

IMAGE_NAME := kong-plugin-poc
IMAGE_TAG := latest

build-mcp:
	@echo "Building MCP server for Linux..."
	cd mcp_server && GOOS=linux GOARCH=amd64 go build -o mcp_server .

build-proxy:
	@echo "Building MCP HTTP proxy for Linux..."
	cd mcp-http-proxy && GOOS=linux GOARCH=amd64 go build -o mcp-http-proxy .

test-mcp: build-mcp
	@echo "Testing MCP server..."
	cd mcp_server && ./test_request.sh | ./mcp_server

clean-mcp:
	@echo "Cleaning MCP server..."
	cd mcp_server && rm -f mcp_server || true

clean-plugin:
	@echo "Cleaning Kong plugin..."
	cd kong-plugin-mcp && make clean

build-plugin:
	@echo "Building Kong MCP plugin..."
	cd kong-plugin-mcp && make build

build-kong: build-mcp build-proxy build-plugin
	@echo "Building Docker image..."
	docker build -t $(IMAGE_NAME):$(IMAGE_TAG) -f kong/Dockerfile .

run-kong: build-kong
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

test-kong: 
	@echo "Testing Kong MCP integration..."
	./test_kong_mcp.sh

clean-kong: stop-kong
	docker rmi $(IMAGE_NAME):$(IMAGE_TAG) || true

clean: clean-mcp clean-plugin clean-kong

all: build-kong
