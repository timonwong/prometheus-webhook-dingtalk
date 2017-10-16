GO    := GO15VENDOREXPERIMENT=1 go

REPO_PATH               ?= github.com/timonwong/prometheus-webhook-dingtalk
TESTARGS                ?= -race
VETARGS                 ?= -all
COVERARGS               ?= -coverprofile=profile.out -covermode=atomic
TEST                    ?= $$(go list ./... | grep -v '/vendor/')
GOFMT_FILES             ?= $$(find . -name '*.go' | grep -v vendor)

export REPO_PATH

all: format test

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

format:
	@echo ">> formatting code"
	@gofmt -w $(GOFMT_FILES)

fmtcheck:
	@echo ">> checking code style"
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

build:
	@sh -c "'$(CURDIR)/scripts/build.sh'"

.PHONY: all format test cover vet fmtcheck build
