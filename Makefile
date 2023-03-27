# Copyright © 2022 Alibaba Group Holding Ltd.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Help by default, even if it's not first
.DEFAULT_GOAL := help

.PHONY: all
all: tidy gen add-copyright format lint cover build

# ==============================================================================
# Build options

ROOT_PACKAGE=github.com/sealerio/sealer
VERSION_PACKAGE=github.com/sealerio/sealer/pkg/version

Dirs=$(shell ls)

# ==============================================================================
# Includes

include scripts/make-rules/common.mk # make sure include common.mk at the first include line
include scripts/make-rules/golang.mk
include scripts/make-rules/image.mk
include scripts/make-rules/copyright.mk
include scripts/make-rules/gen.mk
include scripts/make-rules/ca.mk
include scripts/make-rules/release.mk
include scripts/make-rules/swagger.mk
include scripts/make-rules/dependencies.mk
include scripts/make-rules/tools.mk

# ==============================================================================
# Usage

define USAGE_OPTIONS

Options:
  DEBUG            Whether to generate debug symbols. Default is 0.

  BINS             The binaries to build. Default is all of cmd.
                   This option is available when using: make build/build.multiarch
                   Example: make build BINS="iam-apiserver iam-authz-server"

  IMAGES           Backend images to make. Default is all of cmd starting with iam-.
                   This option is available when using: make image/image.multiarch/push/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-authz-server"

  REGISTRY_PREFIX  Docker registry prefix. Default is marmotedu. 
                   Example: make push REGISTRY_PREFIX=ccr.ccs.tencentyun.com/marmotedu VERSION=v1.6.2

  PLATFORMS        The multiple platforms to build. Default is linux_amd64 and linux_arm64.
                   This option is available when using: make build.multiarch/image.multiarch/push.multiarch
                   Example: make image.multiarch IMAGES="iam-apiserver iam-pump" PLATFORMS="linux_amd64 linux_arm64"

  VERSION          The version information compiled into binaries.
                   The default is obtained from gsemver or git.

  V                Set to 1 enable verbose build. Default is 0.
endef
export USAGE_OPTIONS

# ==============================================================================
# Targets
GIT_TAG := $(shell git describe --exact-match --tags --abbrev=0  2> /dev/null || echo untagged)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

TOOLS_DIR := hack/build.sh

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifneq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

## fmt: Run go fmt against code.
.PHONY: fmt
fmt:
	go fmt ./...

## vet: Run go vet against code.
.PHONY: vet
vet:
	go vet ./...

## lint: Run go lint against code.
.PHONY: lint
lint:
	golangci-lint run -v ./...

## style: code style -> fmt,vet,lint
.PHONY: style
style: fmt vet lint

## build: Build binaries by default
build: clean
	@echo "===========> build sealer and seautil bin"
	@scripts/build.sh

## linux-amd64: Build binaries for Linux (amd64)
linux-amd64: clean
	@echo "Building sealer and seautil binaries for Linux (amd64)"
	GOOS=linux GOARCH=amd64 $(TOOLS_DIR) $(GIT_TAG)

## linux-arm64: Build binaries for Linux (arm64)
linux-arm64: clean
	@echo "Building sealer and seautil binaries for Linux (arm64)"
	GOOS=linux GOARCH=arm64 $(TOOLS_DIR) $(GIT_TAG)

## build-in-docker: sealer should be compiled in linux platform, otherwise there will be GraphDriver problem.
build-in-docker:
	docker run --rm -v ${PWD}:/usr/src/sealer -w /usr/src/sealer registry.cn-qingdao.aliyuncs.com/sealer-io/sealer-build:v1 make linux

## clean: Remove all files that are created by building. 
.PHONY: clean
clean:
	@echo "===========> Cleaning all build output"
	@-rm -rf _output

## install-addlicense: check license if not exist install addlicense tools
install-addlicense:
ifeq (, $(shell which addlicense))
	@{ \
	set -e ;\
	LICENSE_TMP_DIR=$$(mktemp -d) ;\
	cd $$LICENSE_TMP_DIR ;\
	go mod init tmp ;\
	go get -v github.com/google/addlicense ;\
	rm -rf $$LICENSE_TMP_DIR ;\
	}
ADDLICENSE_BIN=$(GOBIN)/addlicense
else
ADDLICENSE_BIN=$(shell which addlicense)
endif

filelicense: SHELL:=/bin/bash
## filelicense: add license
filelicense:
	for file in ${Dirs} ; do \
		if [[  $$file != '_output' && $$file != 'docs' && $$file != 'vendor' && $$file != 'logger' && $$file != 'applications' ]]; then \
			$(ADDLICENSE_BIN)  -y $(shell date +"%Y") -c "Alibaba Group Holding Ltd." -f scripts/LICENSE_TEMPLATE ./$$file ; \
		fi \
    done


## install-gosec: check license if not exist install addlicense tools
install-gosec:
ifeq (, $(shell which gosec))
	@{ \
	set -e ;\
	curl -sfL https://raw.githubusercontent.com/securego/gosec/master/install.sh | sh -s -- -b $(GOBIN) v2.2.0 ;\
	}
GOSEC_BIN=$(GOBIN)/gosec
else
GOSEC_BIN=$(shell which gosec)
endif

gosec: install-gosec
	$(GOSEC_BIN) ./...


install-deepcopy-gen:
ifeq (, $(shell which deepcopy-gen))
	@{ \
	set -e ;\
	LICENSE_TMP_DIR=$$(mktemp -d) ;\
	cd $$LICENSE_TMP_DIR ;\
	go mod init tmp ;\
	go get -v k8s.io/code-generator/cmd/deepcopy-gen ;\
	rm -rf $$LICENSE_TMP_DIR ;\
	}
DEEPCOPY_BIN=$(GOBIN)/deepcopy-gen
else
DEEPCOPY_BIN=$(shell which deepcopy-gen)
endif

HEAD_FILE := scripts/boilerplate.go.txt
INPUT_DIR := github.com/sealerio/sealer/types/api
deepcopy:install-deepcopy-gen
	$(DEEPCOPY_BIN) \
      --input-dirs="$(INPUT_DIR)/v1" \
      -O zz_generated.deepcopy   \
      --go-header-file "$(HEAD_FILE)" \
      --output-base "${GOPATH}/src"
	$(DEEPCOPY_BIN) \
	  --input-dirs="$(INPUT_DIR)/v2" \
	  -O zz_generated.deepcopy   \
	  --go-header-file "$(HEAD_FILE)" \
	  --output-base "${GOPATH}/src"

## help: Show this help info.
.PHONY: help
help: Makefile
	@printf "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:\n"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'
	@echo "$$USAGE_OPTIONS"