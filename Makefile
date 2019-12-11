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

# Needs to be defined before including Makefile.common to auto-generate targets
DOCKER_ARCHS ?= amd64 armv7 arm64

REACT_APP_PATH = web/ui/react-app
REACT_APP_SOURCE_FILES = $(wildcard $(REACT_APP_PATH)/public/* $(REACT_APP_PATH)/src/* $(REACT_APP_PATH)/tsconfig.json $(shell find $(REACT_APP_PATH)/src -type f -print))
REACT_APP_OUTPUT_DIR = web/ui/static/react
REACT_APP_NODE_MODULES_PATH = $(REACT_APP_PATH)/node_modules

include Makefile.common

DOCKER_IMAGE_NAME       ?= prometheus-webhook-dingtalk

STATICCHECK_IGNORE =

$(REACT_APP_NODE_MODULES_PATH): $(REACT_APP_PATH)/package.json $(REACT_APP_PATH)/yarn.lock
	cd $(REACT_APP_PATH) && yarn --frozen-lockfile

$(REACT_APP_OUTPUT_DIR): $(REACT_APP_NODE_MODULES_PATH) $(REACT_APP_SOURCE_FILES)
	@echo ">> building React app"
	@./scripts/build_react_app.sh

.PHONY: build
build: assets common-build

.PHONY: assets
assets: $(REACT_APP_OUTPUT_DIR)
	@echo ">> writing assets"
	# Un-setting GOOS and GOARCH here because the generated Go code is always the same,
	# but the cached object code is incompatible between architectures and OSes (which
	# breaks cross-building for different combinations on CI in the same container).
	GO111MODULE=$(GO111MODULE) GOOS= GOARCH=  $(GO) generate $(GOOPTS) ./template
	GO111MODULE=$(GO111MODULE) GOOS= GOARCH=  $(GO) generate $(GOOPTS) ./web/ui
	@$(GOFMT) -w ./template
	@$(GOFMT) -w ./web/ui

.PHONY: react-app-lint
react-app-lint:
	@echo ">> running React app linting"
	cd $(REACT_APP_PATH) && yarn lint:ci

.PHONY: react-app-lint-fix
react-app-lint-fix:
	@echo ">> running React app linting and fixing errors where possibe"
	cd $(REACT_APP_PATH) && yarn lint

.PHONY: react-app-test
react-app-test: | $(REACT_APP_NODE_MODULES_PATH) react-app-lint
	@echo ">> running React app tests"
	cd $(REACT_APP_PATH) && yarn test --no-watch --coverage

.PHONY: test
#test: common-test react-app-test
test: common-test

.PHONY: clean
clean:
	- @rm -rf asset/assets_vfsdata.go
	- @rm -rf "$(REACT_APP_OUTPUT_DIR)"s
