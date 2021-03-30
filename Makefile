# directory for built binaries
BUILD_DIR = bin

UNAME_S := $(shell uname -s | tr A-Z a-z)

GO_PKGS := $(shell go list ./...)

# service name
SERVICE_NAME = autobots

build: build/$(UNAME_S) ## build binaries based on the OS
	@ln -s $(SERVICE_NAME).$(UNAME_S) $(BUILD_DIR)/$(SERVICE_NAME) || true

build/$(UNAME_S):
	@echo "$@"
	@rm -rf bin/*
	@GOOS=$(UNAME_S) GO111MODULE=on go build -o $(BUILD_DIR)/$(SERVICE_NAME).$(UNAME_S) github.com/ContextLogic/$(SERVICE_NAME)

clean: ## clean ups
	@rm -rf $(BUILD_DIR)
