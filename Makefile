# Run docker container with protobuf image mounted
DOCKER_PROTO_IMAGE = registry-gitlab.i.wish.com/contextlogic/tooling-image/protobuf/master:v0.0.2
DOCKER_RUN = docker run -w /root/protobuf -v `pwd`:/root/protobuf $(DOCKER_PROTO_IMAGE)

# Compile protobuf files
PBC = protoc -I$(SOURCE_DIR) \
			-I/root/include \
			-I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
			--go_out=plugins=grpc:$(TARGET_DIR) $(PROTO_FILES)

# Compile grpc gateway
PBGW = protoc -I$(SOURCE_DIR) \
				-I/root/include \
				-I/go/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
				--grpc-gateway_out=logtostderr=true:$(TARGET_DIR)

# protobuf definition files
PROTO_FILES := $(shell find api -type f -iname "*.proto")
GO_FILES := find . -type f -iname "*.go" | grep -v '^./vendor'
GO_PKGS := $(shell go list ./... | grep -v 'api/proto_gen')

# target directory for compiled protobuf definitions
TARGET_DIR = api/proto_gen

api/proto_gen:
	mkdir api/proto_gen


# source directory for protobuf deifinitions
SOURCE_DIR = api/proto

# directory for built binaries
BUILD_DIR = bin

VERSION := $(shell git describe --tags 2> /dev/null || echo "unreleased")
V_DIRTY := $(shell git describe --exact-match HEAD 2> /dev/null > /dev/null || echo "-unreleased")
GIT     := $(shell git rev-parse --short HEAD)
DIRTY   := $(shell git diff-index --quiet HEAD 2> /dev/null > /dev/null || echo "-dirty")

UNAME_S := $(shell uname -s | tr A-Z a-z)

# service name
SERVICE_NAME = autobots

help: ## Show this help
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {sub("\\\\n",sprintf("\n%22c"," "), $$2);printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

all: clean proto build ## clean, compile protobuf definitions and build the binary

compose: clean proto

clean: ## clean ups
	@rm -rf $(TARGET_DIR) $(BUILD_DIR)

# compile protobuf definitions to code-gen *.pb.go
# using a docker container of deps of protoc, protoc-gen-go installed.
proto: $(TARGET_DIR) $(SOURCE_DIR) ## compile protobuf definitons using docker container
	@if [ "$(docker images -q ${DOCKER_PROTO_IMAGE} 2> /dev/null)" != "" ]; then \
		docker login registry-gitlab.i.wish.com -u ${GITLAB_USERNAME} -p ${GITLAB_ACCESS_TOKEN}; \
	fi
	@$(DOCKER_RUN) $(PBC)

gateway: $(TARGET_DIR) $(SOURCE_DIR) ## compile protobuf definitons using docker container
	@if [ "$(docker images -q ${DOCKER_PROTO_IMAGE} 2> /dev/null)" != "" ]; then \
		docker login registry-gitlab.i.wish.com -u ${GITLAB_USERNAME} -p ${GITLAB_ACCESS_TOKEN}; \
	fi
	@for FILE in $(PROTO_FILES); do \
		$(DOCKER_RUN) $(PBGW) $$FILE; \
	done
vendor: go.mod go.sum ## pull the vendor pkgs for deps
	@GO111MODULE=on go mod vendor

build: $(TARGET_DIR) build/$(UNAME_S) ## build binaries based on the OS
	@ln -s $(SERVICE_NAME).$(UNAME_S) $(BUILD_DIR)/$(SERVICE_NAME) || true

build/$(UNAME_S):
	@echo "$@"
	@GOOS=$(UNAME_S) CGO_ENABLED=0 GO111MODULE=on go build -o $(BUILD_DIR)/$(SERVICE_NAME).$(UNAME_S) -ldflags \
   		"-X github.com/ContextLogic/$(SERVICE_NAME).Version=$(VERSION)$(V_DIRTY) \
   		 -X github.com/ContextLogic/$(SERVICE_NAME).Git=$(GIT)$(DIRTY)" \
   		github.com/ContextLogic/$(SERVICE_NAME)

run: vendor ## run the server from source code
	@go run main.go server -c config/service.json

watch: ## Watch .go files for changes and rerun make run (requires entr, see https://github.com/clibs/entr)
	@$(GO_FILES) | entr -rc $(MAKE) run

## Unit test and test coverager
test: ## run unit tests
	@echo "running unit tests"
	@go test -v -cover $(GO_PKGS) 2> /dev/null

coverage: ## run test with coverage
	@go test -coverprofile=/tmp/cover $(GO_PKGS)
	@go tool cover -html=/tmp/cover -o coverage.html
	@rm /tmp/cover

## Lint and format code
GOIMPORTS ?= goimports -local=github.com/ContextLogic/$(SERVICE_NAME)
GETGOLINT := $(shell go get -u golang.org/x/lint/golint 2> /dev/null)
GOLINT := golint
GOFMT := gofmt

imports:
	@$(GOIMPORTS) -w $(shell $(GO_FILES))

fmt:
	@$(GOFMT) -w -s $(shell $(GO_FILES))

lint:
	@$(GETGOLINT)
	@for FILE in $(shell $(GO_FILES)); do $(GOLINT) $$FILE;  done;

