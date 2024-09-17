# Stuart docker registry URL
REGISTRY_URL ?= "nicolasverle"
BUILDTIME ?= $(shell date -u '+%Y-%m-%d_%H:%M:%S')
COMMIT ?= $(shell git rev-parse HEAD)

# Image URL to use all building/pushing image targets
IMG ?= "${REGISTRY_URL}/casestudy:${COMMIT}"

## Tool Binaries
KUSTOMIZE ?= $(shell which kustomize)

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk commands is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Development

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

##@ Build

.PHONY: build
build: fmt vet ## Build manager binary.
	go build -o bin/linkextractor main.go

.PHONY: release
release: fmt vet ## Build manager binary.
	CGO_ENABLED=0 go build -ldflags "-s -w" -a -trimpath  -o bin/linkextractor main.go

.PHONY: docker-login
docker-login:
	@docker login ${REGISTRY_URL} -u ${DOCKER_REGISTRY_USER} -p ${DOCKER_REGISTRY_PASSWORD}

.PHONY: docker-build
docker-build: release 
	docker build -t ${IMG} .

.PHONY: docker-push
docker-push: docker-build 
	docker push ${IMG}

# ##@ Deployment

ifndef ignore-not-found
  ignore-not-found = false
endif

.PHONY: generate
generate: 
	$(KUSTOMIZE) build manifests/test 

.PHONY: deploy
deploy: 
	$(KUSTOMIZE) build manifests/test | kubectl apply -f -

