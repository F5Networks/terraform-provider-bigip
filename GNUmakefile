TEST = ./bigip
TESTARGS = -v
PROJ = terraform-provider-bigip

ARCHS = amd64 386
OS = windows darwin linux

OUT_DIR = target
BIN_DIR = $(OUT_DIR)/bin
PKG_DIR = $(OUT_DIR)/pkg

PKGS = $(foreach arch,$(ARCHS),$(foreach os,$(OS),$(PKG_DIR)/$(PROJ)_$(os)_$(arch)$(PKG_SUFFIX)))
BINS = $(foreach arch,$(ARCHS),$(foreach os,$(OS),$(BIN_DIR)/$(os)_$(arch)/$(PROJ)))

default: bin

build:
	@go build ./...

bin: test
	@gox -help >/dev/null 2>&1 ; if [ $$? -ne 2 ]; then \
		go get github.com/mitchellh/gox; \
	fi
	@gox -output="$(BIN_DIR)/{{.OS}}_{{.Arch}}/terraform-{{.Dir}}" -arch="$(ARCHS)" -os="$(OS)" "github.com/f5devcentral/terraform-provider-bigip"

dist:
	@mkdir -p $(PKG_DIR) 2>/dev/null
	@for arch in $(ARCHS); do \
		for os in $(OS); do \
			echo "$(PKG_DIR)/$(PROJ)_$${os}_$${arch}.tar.gz"; \
			tar czf $(PKG_DIR)/$(PROJ)_$${os}_$${arch}.tar.gz -C $(BIN_DIR)/$${os}_$${arch} .; \
		done \
	done

fmt:
	@gofmt -l -w . bigip/

# vet runs the Go source code static analysis tool `vet` to find
# any common errors.
vet:
	@go tool vet 2>/dev/null ; if [ $$? -eq 3 ]; then \
		go get golang.org/x/tools/cmd/vet; \
	fi
	@echo "go tool vet $(VETARGS) ."
	@go tool vet $(VETARGS) $$(ls -d */ | grep -v vendor) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

test: build
	@TF_ACC=true go test $(TEST) $(TESTARGS) -timeout=300s -parallel=1

testacc: fmt build
	@if [[ "$(BIGIP_USER)" == "admin" || "$(BIGIP_HOST)" == "13.57.25.130" || "-z $(BIGIP_PASSWORD)" == "cisco123" ]]; then \
		echo "ERROR: BIGIP_USER, BIGIP_PASSWORD and BIGIP_HOST must be set."; \
		exit 1; \
	fi
	@TF_ACC=1 go test $(TEST) $(TESTARGS) -timeout 120m
	e2etest: generate
		TF_ACC=1 go test -v ./command/e2etest

test-compile: fmtcheck generate
		@if [ "$(TEST)" = "./..." ]; then \
			echo "ERROR: Set TEST to a specific package. For example,"; \
			echo "  make test-compile TEST=./builtin/providers/bigip"; \
			exit 1; \
		fi
		go test -c $(TEST) $(TESTARGS)

clean:
	@go clean
	@rm -rf target/
