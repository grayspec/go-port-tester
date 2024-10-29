# Makefile for cross-compilation

BINARY_NAME_SERVER=server
BINARY_NAME_CLIENT=client
BUILD_DIR=build

.PHONY: all clean build

# Windows
build-windows:
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(BINARY_NAME_SERVER).exe ./server/server.go
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/windows/$(BINARY_NAME_CLIENT).exe ./client/client.go

# MacOS
build-macos:
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/macos/$(BINARY_NAME_SERVER) ./server/server.go
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/macos/$(BINARY_NAME_CLIENT) ./client/client.go

# Linux
build-linux:
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(BINARY_NAME_SERVER) ./server/server.go
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/linux/$(BINARY_NAME_CLIENT) ./client/client.go

# All builds
all: build-windows build-macos build-linux

# Clean build artifacts
clean:
	rm -rf $(BUILD_DIR)
