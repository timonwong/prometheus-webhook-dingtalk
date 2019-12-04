GO    := go

REPO_PATH               ?= github.com/timonwong/prometheus-webhook-dingtalk
TESTARGS                ?= -v -race
COVERARGS               ?= -coverprofile=coverage.txt -covermode=atomic
TEST                    ?= $(shell go list ./... | grep -v '/vendor/')
TESTPKGS                ?= $(shell go list ./... | grep -v '/cmd/')
GOFMT_FILES             ?= $(shell find . -name '*.go' | grep -v vendor | xargs)
FIRST_GOPATH            ?= $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))
PROMU                   ?= $(FIRST_GOPATH)/bin/promu
GOLANGCI_LINT           ?= $(FIRST_GOPATH)/bin/golangci-lint
GOCOV                   ?= $(FIRST_GOPATH)/bin/gocov
GOCOV_HTML              ?= $(FIRST_GOPATH)/bin/gocov-html
GOIMPORTS               ?= $(FIRST_GOPATH)/bin/goimports
GO_BINDATA              ?= $(FIRST_GOPATH)/bin/go-bindata

PREFIX                  ?= $(shell pwd)
BIN_DIR                 ?= $(shell pwd)
DOCKER_IMAGE_NAME       ?= prometheus-webhook-dingtalk
DOCKER_IMAGE_TAG        ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

export GOFLAGS=-mod=vendor
export REPO_PATH

_comma := ,
_space :=
_space +=


.PHONY: all
all: format test


$(GOLANGCI_LINT):
	@echo ">> installing linter"
	@$(GO) install "github.com/golangci/golangci-lint/cmd/golangci-lint"


$(GOCOV):
	@echo ">> installing gocov tool"
	@$(GO) install "github.com/axw/gocov/gocov"


$(GOCOV_HTML):
	@echo ">> installing gocov-html tool"
	@$(GO) install "github.com/matm/gocov-html"


$(GOIMPORTS):
	@echo ">> installing goimports tool"
	@$(GO) install "golang.org/x/tools/cmd/goimports"


$(GO_BINDATA):
	@echo ">> installing go-bindata"
	@$(GO) install github.com/go-bindata/go-bindata/...


.PHONY: dep
dep:
	@$(GO) mod vendor


.PHONY: test
test:
	@echo ">> running tests"
	@$(GO) test $(TESTARGS) $(TEST)


.PHONY: cover
cover: $(GOCOV) $(GOCOV_HTML)
	@echo ">> running test coverage"
	@rm -f coverage.txt
	@$(GO) test $(TESTARGS) $(COVERARGS) -coverpkg "$(subst $(_space),$(_comma),$(TESTPKGS))" $(TEST) && \
		$(GOCOV) convert coverage.txt >coverage.json && \
		$(GOCOV) report coverage.json && \
		$(GOCOV_HTML) coverage.json >coverage.html


.PHONY: lint
lint: $(GOLANGCI_LINT)
	@echo ">> linting code"
	@$(GOLANGCI_LINT) run \
		--deadline=10m \
		--disable-all \
		--enable=gofmt \
		--enable=goimports \
		--enable=govet \
		--enable=typecheck \
		--enable=varcheck \
		--enable=errcheck \
		--enable=staticcheck \
		--enable=gas \
		--enable=ineffassign \
		--enable=gosimple \
		--enable=maligned \
		--enable=golint \
		./...


.PHONY: format
format: $(GOIMPORTS)
	@echo ">> formatting code"
	@$(GOIMPORTS) -local "git.meiqia.com" -w $(GOFMT_FILES)


.PHONY: assets
assets: $(GO_BINDATA) template/internal/deftmpl/bindata.go

template/internal/deftmpl/bindata.go: template/default.tmpl
	@$(GO_BINDATA) $(bindata_flags) -mode 420 -modtime 1 -pkg deftmpl -o template/internal/deftmpl/bindata.go template/default.tmpl


.PHONY: build
build: promu
	@echo ">> building binaries"
	@$(PROMU) build --prefix $(PREFIX)


# Will build both the front-end as well as the back-end
.PHONY: build-all
build-all: assets build


.PHONY: tarball
tarball: promu
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(PREFIX) $(BIN_DIR)


.PHONY: docker
docker:
	@echo ">> building docker image"
	@docker build -t "$(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)" .


.PHONY: promu
promu:
	@echo ">> installing promu"
	@$(GO) install github.com/prometheus/promu
