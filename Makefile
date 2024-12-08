ROOT := $(shell git rev-parse --show-toplevel)

GIT_SHA := $(shell git rev-parse HEAD)
GIT_SHA_SHORT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION := $(shell git describe --tags)-$(GIT_SHA_SHORT)
LDFLAGS := -s -w \
        -X 'github.com/jlewi/ddctl/pkg/version.Date=$(DATE)' \
        -X 'github.com/jlewi/ddctl/pkg/version.Version=$(subst v,,$(VERSION))' \
        -X 'github.com/jlewi/ddctl/pkg/version.Commit=$(GIT_SHA)'

build: build-dir
	CGO_ENABLED=0 go build -o .build/ddctl -ldflags="$(LDFLAGS)" github.com/jlewi/ddctl

build-dir:
	mkdir -p .build

tidy:
	gofmt -s -w .
	goimports -w .
	

lint:
	# golangci-lint automatically searches up the root tree for configuration files.
	golangci-lint run

test:	
	GITHUB_ACTIONS=true go test -v ./...