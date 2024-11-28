APP_NAME=ai-proxy
VERSION=$(shell git describe --tags --always --long --dirty)
BUILD_DIR=./cmd
OUTPUT_DIR=bin

# Detect if the OS is Windows
ifeq ($(OS),Windows_NT)
    SHELL := cmd.exe
    SET_ENV = set
	SEP = &
	RM = del /F /Q
	PATHSEP = \\
else
    SHELL := /bin/bash
    SET_ENV = env
	SEP = &&
	RM = rm -f
	PATHSEP = /
endif

# Build for Windows (64-bit)
windows: 	
	$(SET_ENV) GOOS=windows$(SEP) $(SET_ENV) GOARCH=amd64$(SEP) go build -o $(OUTPUT_DIR)$(PATHSEP)$(APP_NAME)_windows.exe -ldflags="-s -w -X main.version=$(VERSION)" $(BUILD_DIR)

# Build for Linux (64-bit)
linux: 
	$(SET_ENV) GOOS=linux$(SEP) $(SET_ENV) GOARCH=amd64$(SEP) go build -o $(OUTPUT_DIR)$(PATHSEP)$(APP_NAME)_linux -ldflags="-s -w -X main.version=$(VERSION)" $(BUILD_DIR)

# Build for macOS (ARM 64-bit)
macos-arm: 	
	$(SET_ENV) GOOS=darwin$(SEP) $(SET_ENV) GOARCH=arm64$(SEP) go build -o $(OUTPUT_DIR)$(PATHSEP)$(APP_NAME)_darwin_arm64 -ldflags="-s -w -X main.version=$(VERSION)" $(BUILD_DIR)

# Default target to build all versions
all: windows linux macos-arm

# Clean up built binaries
clean:
	$(RM) $(OUTPUT_DIR)$(PATHSEP)$(APP_NAME)_*

test:
	go test -v -cover ./...

.PHONY: test all clean windows linux macos-arm