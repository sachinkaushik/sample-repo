# SPDX-FileCopyrightText: 2022-present Intel Corporation
#
# SPDX-License-Identifier: LicenseRef-Intel

DOCKER_TAG                      ?= latest
DOCKER_REPOSITORY               ?= sca_0_1_x/
CONTAINER_NAME                  ?= roc-rest-api
DOCKER_REGISTRY                 ?= ""
PLATFORM                        ?= --platform linux/x86_64

mod-update: ## Download the dependencies
	go mod tidy
	go mod vendor

build: ## build the Go binaries
build: mod-update
	go build -mod=vendor -o roc ./cmd/main.go

docker-build: # Build roc-api-server docker image
docker-build:
	@go mod tidy
	@go mod vendor
	docker build ${PLATFORM} \
          -t ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${CONTAINER_NAME}:${DOCKER_TAG} \
          -f build/Dockerfile .
	@rm -rf vendor

docker-push: # Push the docker image into repository
docker-push:
	docker push ${DOCKER_REGISTRY}${DOCKER_REPOSITORY}${CONTAINER_NAME}:${DOCKER_TAG}
