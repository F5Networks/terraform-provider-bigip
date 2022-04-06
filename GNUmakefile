#Original work from https://github.com/DealerDotCom/terraform-provider-bigip
#Modifications Copyright 2019 F5 Networks Inc.
#This Source Code Form is subject to the terms of the Mozilla Public License, v. 2.0.
#If a copy of the MPL was not distributed with this file,You can obtain one at https://mozilla.org/MPL/2.0/.

TEST?=./...
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
PKG_NAME=bigip
WEBSITE_REPO=github.com/hashicorp/terraform-website

default: build

build: fmtcheck
	go install

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(TEST) -v -sweep=$(SWEEP) $(SWEEPARGS)

test: fmtcheck
	go test $(TEST) -timeout=30s -parallel=4

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

fmt:
	@echo "==> Fixing source code with gofmt..."
	gofmt -s -w ./$(PKG_NAME)

# Currently required by tf-deploy compile
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	golangci-lint run ./...

tools:
	@echo "==> installing required tooling..."
	go install github.com/katbyte/terrafmt@d117dc888b365322b5c23e5c41cff1e27c336f7f
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH || $$GOPATH)/bin v1.41.1

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

website:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

website-lint:
	@echo "==> Checking website against linters..."
	@misspell -error -source=text website/

website-test:
ifeq (,$(wildcard $(GOPATH)/src/$(WEBSITE_REPO)))
	echo "$(WEBSITE_REPO) not found in your GOPATH (necessary for layouts and assets), get-ting..."
	git clone https://$(WEBSITE_REPO) $(GOPATH)/src/$(WEBSITE_REPO)
endif
	@$(MAKE) -C $(GOPATH)/src/$(WEBSITE_REPO) website-provider-test PROVIDER_PATH=$(shell pwd) PROVIDER_NAME=$(PKG_NAME)

tflint:
	./scripts/run-tflint.sh

tffmtfix:
	@echo "==> Fixing acceptance test terraform blocks code with terrafmt..."
	@find ./bigip -type f -name "*_test.go" | sort | while read f; do terrafmt fmt -f $$f; done
	@echo "==> Fixing docs terraform blocks code with terrafmt..."
	@find ./docs -type f -name "*.md" | sort | while read f; do terrafmt fmt $$f; done

.PHONY: build sweep test testacc fmt fmtcheck lint tools test-compile website website-lint website-test tflint tffmtfix
