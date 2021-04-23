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
	@GOOS=$(UNAME_S) GO111MODULE=on go build -o $(BUILD_DIR)/$(SERVICE_NAME).$(UNAME_S) github.com/ContextLogic/$(SERVICE_NAME)/pkg

vendor: go.mod go.sum ## pull the vendor pkgs for deps
	@GO111MODULE=on go mod vendor

docker:
	@docker build --build-arg ITA_JOB_TOKEN=${GITLAB_ACCESS_TOKEN} \
		--build-arg ITA_PROJECT_NAME=$(SERVICE_NAME) -f Dockerfile.local -t contextlogic/$(SERVICE_NAME) .

clean: ## clean ups
	@rm -rf $(BUILD_DIR)

update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif