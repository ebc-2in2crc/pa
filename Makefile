.DEFAULT_GOAL := help

GOCMD := env GO111MODULE=on go
GOMOD := $(GOCMD) mod
GOBUILD := $(GOCMD) build
GOINSTALL := $(GOCMD) install
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
NAME := pa
CURRENT := $(shell pwd)
BUILDDIR := ./build
BINDIR := $(BUILDDIR)/bin
PKGDIR := $(BUILDDIR)/pkg
DISTDIR := $(BUILDDIR)/dist

GOFLAGS := -trimpath
VERSION := $(shell git describe --tags --abbrev=0)
LDFLAGS := -X 'github.com/ebc-2in2crc/$(NAME)/cmd.version=$(VERSION)'
GOXOSARCH := "darwin/amd64 darwin/arm64 windows/386 windows/amd64 linux/386 linux/amd64"
GOXOUTPUT := "$(PKGDIR)/$(NAME)_{{.OS}}_{{.Arch}}/{{.Dir}}"

export GO111MODULE=on

.PHONY: deps
## Install dependencies
deps:
	$(GOMOD) download

.PHONY: devel-deps
## Install dependencies for develop
devel-deps: deps
	$(GOINSTALL) golang.org/x/tools/cmd/goimports@latest
	$(GOINSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GOINSTALL) github.com/Songmu/make2help/cmd/make2help@latest
	$(GOINSTALL) github.com/mitchellh/gox@latest
	$(GOINSTALL) github.com/tcnksm/ghr@latest

.PHONY: build
## Build binaries
build: deps
	$(GOBUILD) -trimpath -ldflags "$(LDFLAGS)" -o $(BINDIR)/$(NAME) ./cmd/$(NAME)

.PHONY: cross-build
## Cross build binaries
cross-build:
	rm -rf $(PKGDIR)
	env GOFLAGS="$(GOFLAGS)" gox -osarch=$(GOXOSARCH) -ldflags "$(LDFLAGS)" -output=$(GOXOUTPUT) ./cmd/$(NAME)

.PHONY: package
## Make package
package: cross-build
	rm -rf $(DISTDIR)
	mkdir $(DISTDIR)
	pushd $(PKGDIR) > /dev/null && \
		for P in `ls | xargs basename`; do zip -r $(CURRENT)/$(DISTDIR)/$$P.zip $$P; done && \
		popd > /dev/null

.PHONY: release
## Release package to Github
release: package
	ghr $(VERSION) $(DISTDIR)

.PHONY: zsh_completion
## Generate ZSH completion
zsh_completion: build
	$(BINDIR)/$(NAME) completion zsh > ~/.zsh/completions/_$(NAME)
	. ~/.zsh/completions/_$(NAME)

.PHONY: test
## Run tests
test: deps
	$(GOTEST) -v ./...

.PHONY: lint
## Lint
lint: devel-deps
	golangci-lint run ./...

.PHONY: fmt
## Format source codes
fmt: deps
	find . -name "*.go" -not -path "./vendor/*" | xargs goimports -w

.PHONY: clean
clean:
	$(GOCLEAN)
	rm -rf $(BUILDDIR)

.PHONY: help
## Show help
help:
	@make2help $(MAKEFILE_LIST)
