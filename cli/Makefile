# GENERAL
AppName := TAIGO CLI
BINARY_NAME := taigo
COPYRIGHT_TEXT := Taigo CLI | MIT License 2021 - $(shell date '+%Y')
BUILD_TIME := $(shell date '+%c')

# Configure OS-specific binary suffix
ifeq ($(OS),Windows_NT)
	BINARY_SUFFIX := .exe
else
	BINARY_SUFFIX :=
endif

# Get/Set BUILD_TYPE
ifdef BUILD_TYPE
	BUILD_TYPE := $(BUILD_TYPE)
else
#   SET YOUR CUSTOM DEFAULT RELEASE TYPE BELOW
	BUILD_TYPE := release
endif

# Get/Set LATEST_GIT_TAG
LATEST_GIT_TAG := $(shell git describe --tags --abbrev=0)
ifndef LATEST_GIT_TAG
	LATEST_GIT_TAG := no-version
endif

# Get/Set LATEST_GIT_COMMIT
LATEST_GIT_COMMIT := $(shell git log -n 1 --pretty=format:"%H")
ifndef LATEST_GIT_COMMIT
	LATEST_GIT_COMMIT := commit-id-could-not-be-retrieved
endif

# TECHNICAL VARS
GO_LD_FLAGS := -ldflags "-X 'main.AppName=$(AppName)' -X 'main.AppBuildDate=$(BUILD_TIME)' -X 'main.AppBuildType=$(BUILD_TYPE)' -X 'main.AppCopyrightText=$(COPYRIGHT_TEXT)' -X 'main.AppSemVersion=$(LATEST_GIT_TAG)' -X 'main.GitCommit=$(LATEST_GIT_COMMIT)'"
export PATH := $(shell go env GOPATH)/bin:$(PATH)

# **MAKE TARGETS**
info:
	@echo "Choose from the following targets:"
	@echo "  * build (defaults to your host/OS settings)"
	@echo "  * build-darwin"
	@echo "  * build-linux"
	@echo "  * build-windows"
	@echo "    -----------------------------------------------"
	@echo "  * build-all"

build:
	@echo "Building for host system: $(shell go env GOOS)/$(shell go env GOARCH)"
	$(eval BUILD_TYPE := internal/verification)
	$(eval GO_LD_FLAGS := -ldflags "-X 'main.AppName=$(AppName)' -X 'main.AppBuildDate=$(BUILD_TIME)' -X 'main.AppBuildType=$(BUILD_TYPE)' -X 'main.AppCopyrightText=$(COPYRIGHT_TEXT)' -X 'main.AppSemVersion=$(LATEST_GIT_TAG)' -X 'main.GitCommit=$(LATEST_GIT_COMMIT)'")
	@go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)$(BINARY_SUFFIX)

build-darwin:
	@GOOS=darwin GOARCH=amd64 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	@GOOS=darwin GOARCH=arm64 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-darwin-arm64

build-linux:
	@GOOS=linux GOARCH=386 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-linux-386
	@GOOS=linux GOARCH=amd64 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	@GOOS=linux GOARCH=arm go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-linux-arm
	@GOOS=linux GOARCH=arm64 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-linux-arm64

build-windows: generate-win-versioninfo
	@GOOS=windows GOARCH=386 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-windows-386.exe
	@GOOS=windows GOARCH=amd64 go build $(GO_LD_FLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe

build-enterprise-linux:
	@GOOS=linux GOARCH=amd64 go build $(GO_LD_FLAGS) --tags enterprise -o dist/$(BINARY_NAME).enterprise-linux-amd64

build-all: clean build-darwin build-linux build-windows

clean:
	rm -rf ./dist/*
	rm -rf ./resource_*.syso

generate-win-versioninfo:
	@go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest
	goversioninfo -platform-specific=true -file-version "$(LATEST_GIT_TAG).0" -product-version "$(LATEST_GIT_TAG)" -copyright "$(COPYRIGHT_TEXT)" -private-build "$(LATEST_GIT_TAG)" 
