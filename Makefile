# +-------------------------------------------------------------------------
# | Copyright (C) 2025 Zenlayer, Inc.
# +-------------------------------------------------------------------------
# | Licensed under the Apache License, Version 2.0 (the "License");
# | you may not use this work except in compliance with the License.
# | You may obtain a copy of the License in the LICENSE file, or at:
# |
# | http://www.apache.org/licenses/LICENSE-2.0
# |
# | Unless required by applicable law or agreed to in writing, software
# | distributed under the License is distributed on an "AS IS" BASIS,
# | WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# | See the License for the specific language governing permissions and
# | limitations under the License.
# +-------------------------------------------------------------------------

.PHONY: bin image

DISK_VERSION=v1.0.0

DISK_IMAGE_NAME=docker.io/zenlayer297/zeccsi
CONTAINER_CMD=docker

GOARCH=$(shell go env GOARCH 2>/dev/null)
ifeq ($(GOARCH),amd64)
	BASE_IMAGE=docker.io/ubuntu:24.04
	BUILD_IMAGE=docker.io/library/golang:1.24
else 
	exit 1
endif


.container-cmd:
	@test -n "$(shell which $(CONTAINER_CMD) 2>/dev/null)" || { echo "Missing container support, install Docker"; exit 1; }

bin:
	if [ ! -d ./vendor ]; then (go mod tidy && go mod vendor); fi
	CGO_ENABLED=0 GOOS=linux go build -a -mod=vendor  -ldflags "-s -w" -o  _output/zeccsi ./cmd/

image: GOARCH?=$(shell go env GOARCH 2>/dev/null)
image: .container-cmd
	$(CONTAINER_CMD) build -t ${DISK_IMAGE_NAME}:${DISK_VERSION} -f deploy/docker/Dockerfile  .  --build-arg BASE_IMAGE=$(BASE_IMAGE) --build-arg BUILD_IMAGE=$(BUILD_IMAGE)

clean:
	go clean -mod=vendor -r -x
	rm -rf ./_output
	docker image prune --force
	docker rmi ${DISK_IMAGE_NAME}:${DISK_VERSION}
	
