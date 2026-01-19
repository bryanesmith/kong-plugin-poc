.PHONY: build clean

PLUGIN_NAME := mcp-bridge
OUTPUT := $(PLUGIN_NAME)

build:
	@echo "Building Kong Go plugin (embedded server) for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(OUTPUT) .

clean:
	rm -f $(OUTPUT)
