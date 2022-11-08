# SPDX-FileCopyrightText: (C) 2022 Intel Corporation
# SPDX-License-Identifier: LicenseRef-Intel

# If any command in a pipe has nonzero status, return that status
SHELL = bash -o pipefail

export CGO_ENABLED=1
export GO111MODULE=on

.PHONY: build
PLATFORM ?= --platform linux/x86_64

VERSION                     ?= $(shell cat ./VERSION)

DOCKER_REPOSITORY           ?= 
ROC-ADAPTER_IMAGE_NAME      ?= sra-adapter
ROC-ADAPTER_VERSION         ?= latest

IMG_REGISTRY := amr-registry.caas.intel.com/one-intel-edge/maestro-i/${ROC-ADAPTER_IMAGE_NAME}:${ROC-ADAPTER_VERSION}

## Docker labels. Only set ref and commit date if committed
DOCKER_LABEL_VCS_URL        ?= $(shell git remote get-url $(shell git remote))
DOCKER_LABEL_VCS_REF        = $(shell git rev-parse HEAD)
DOCKER_LABEL_BUILD_DATE     ?= $(shell date -u "+%Y-%m-%dT%H:%M:%SZ")
DOCKER_LABEL_COMMIT_DATE    = $(shell git show -s --format=%cd --date=iso-strict HEAD)

DOCKER_EXTRA_ARGS           ?=
DOCKER_BUILD_ARGS ?= \
        ${DOCKER_EXTRA_ARGS} \
        --build-arg org_label_schema_version="${VERSION}" \
        --build-arg org_label_schema_vcs_url="${DOCKER_LABEL_VCS_URL}" \
        --build-arg org_label_schema_vcs_ref="${DOCKER_LABEL_VCS_REF}" \
        --build-arg org_label_schema_build_date="${DOCKER_LABEL_BUILD_DATE}" \
        --build-arg org_opencord_vcs_commit_date="${DOCKER_LABEL_COMMIT_DATE}" \
        --build-arg org_opencord_vcs_dirty="${DOCKER_LABEL_VCS_DIRTY}" \
		--build-arg LOCAL_AETHER_MODELS=${LOCAL_AETHER_MODELS}

PKG               := github.com/intel-innersource/frameworks.edge.one-intel-edge.springboard.reference-implementation.roc-adapters.sra

BINDIR            := build

# Note: GOPRIVATE is set in Jenkinsfile or in developer environemnt
GOPATH            := $(shell go env GOPATH)
GOCMD             := go
GOLANG_CLI_LINT_VERSION := v1.49.0
GOCOVER_COBERTURA_VERSION := latest
MOCKGEN_VERSION := v1.6.0

all: go-build docker-build

test-dependency:
	go install github.com/golang/mock/mockgen@$(MOCKGEN_VERSION)
	go install github.com/boumenot/gocover-cobertura@$(GOCOVER_COBERTURA_VERSION)
	go generate ./...

# @HELP run various tests
test: go-build go-test go-lint go-deps images
	# license check is temporarily disabled

# @HELP build the adapter
go-build:
	mkdir -p $(BINDIR)
	go build -o $(BINDIR)/_output/roc-adapter ./cmd/roc-adapter	

# @HELP format go code using go fmt
go-format:
	@go fmt $(shell sh -c "go list ./...")

# @HELP run go lint
go-lint:
	${GOPATH}/bin/golangci-lint --version || go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANG_CLI_LINT_VERSION}
	${GOPATH}/bin/golangci-lint --version
	${GOPATH}/bin/golangci-lint run pkg/... cmd/... --config .golangci.yml		

# @HELP run go unit tests
go-test: test-dependency
	$(GOCMD) test -gcflags=-l `go list $(PKG)/cmd/... $(PKG)/pkg/... | grep -v "/mocks" | grep -v "/test/"` -v -coverprofile=coverage.txt -covermode count
	${GOPATH}/bin/gocover-cobertura < coverage.txt > coverage.xml	
	#$(GOCMD) tool cover -html=cover.out -o cover.html
	#$(GOCMD) tool cover -func cover.out -o cover.function-coverage.log

# @HELP ensure that the required dependencies are in place
go-deps:
	go build -v ./...
	bash -c "diff -u <(echo -n) <(git diff go.mod)"
	bash -c "diff -u <(echo -n) <(git diff go.sum)"	

# @HELP build docker images
docker-build:
	go mod tidy
	go mod vendor
	docker build ${PLATFORM} . -f Dockerfile \
	$(DOCKER_BUILD_ARGS) \
	-t ${DOCKER_REPOSITORY}roc-adapter:${ROC-ADAPTER_VERSION}

# @HELP publish docker images to CI/CD system
docker-cicd: docker-build
	docker tag ${DOCKER_REPOSITORY}roc-adapter:${ROC-ADAPTER_VERSION} ${IMG_REGISTRY}
	docker push ${IMG_REGISTRY}	

# @HELP remove all the build artifacts
clean::
	rm -rf ./build/_output
	rm -rf ./vendor
	rm -rf ./cmd/roc-adapter/roc-adapter
	rm -rf ./stratum

# synonyms for prior Makefile targets
build: go-build
deps: go-deps
images: docker-build
linters: go-lint
unit-test: go-test