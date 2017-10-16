# Copyright 2015 The Prometheus Authors
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

GO    := GO15VENDOREXPERIMENT=1 go
PROMU := $(GOPATH)/bin/promu

PREFIX                  ?= $$(pwd)
BIN_DIR                 ?= $$(pwd)
DOCKER_IMAGE_NAME       ?= uwsgi-exporter
DOCKER_IMAGE_TAG        ?= $(subst /,-,$$(git rev-parse --abbrev-ref HEAD))

TESTARGS                ?= -race -v
VETARGS                 ?= -all
COVERARGS               ?= -coverprofile=profile.out -covermode=atomic
TEST                    ?= $$(go list ./... | grep -v '/vendor/')
GOFMT_FILES             ?= $$(find . -name '*.go' | grep -v vendor)

all: format build test

test: fmtcheck
	@echo ">> running tests"
	@$(GO) test $(TEST) $(TESTARGS)

cover: fmtcheck
	@echo ">> running test coverage"
	rm -f coverage.txt
	@for d in $(TEST); do \
		go test $(TESTARGS) $(COVERARGS) $$d; \
		if [ -f profile.out ]; then \
			cat profile.out >> coverage.txt; \
			rm profile.out; \
		fi \
	done

vet:
	@echo ">> vetting code"
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)

tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)

docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .

promu:
	@GOOS=$(shell uname -s | tr A-Z a-z) \
		GOARCH=$(subst x86_64,amd64,$(patsubst i%86,386,$(shell uname -m))) \
		$(GO) get -u github.com/prometheus/promu

format:
	@echo ">> formatting code"
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@echo ">> checking code style"
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

.PHONY: all format build test cover vet tarball docker promu fmtcheck
