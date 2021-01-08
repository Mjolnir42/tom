# vim: set ft=make ffs=unix fenc=utf8:
# vim: set noet ts=4 sw=4 tw=72 list:
#
TOMVER != git describe --tags --abbrev=0
BRANCH != git rev-parse --symbolic-full-name --abbrev-ref HEAD
GITHASH != git rev-parse --short HEAD

all: install

install: install_freebsd

install_all: install_freebsd install_linux

install_freebsd: generate
	@echo "Building FreeBSD ...."
	@env GOOS=freebsd GOARCH=amd64 go install -ldflags "-X main.tomVersion=$(TOMVER)-$(GITHASH)/$(BRANCH)" ./...

install_linux: generate
	@echo "Building Linux ...."
	@env GOOS=linux GOARCH=amd64 go install -ldflags "-X main.tomVersion=$(TOMVER)-$(GITHASH)/$(BRANCH)" ./...

generate: sanitize
	@echo "Generating ...."
	@go generate ./cmd/...

sanitize: build check

check: vet

build:
	@echo "Compiling ...."
	@go build ./...

vet:
	@echo "Running 'go vet' ...."
	@go vet ./cmd/tom/
	@go vet ./cmd/tomd/
	@go vet ./internal/config/
	@go vet ./internal/core/
	@go vet ./internal/handler/
	@go vet ./internal/model/asset/
	@go vet ./internal/model/meta/
	@go vet ./internal/msg/
	@go vet ./internal/rest/
	@go vet ./internal/stmt/
	@go vet ./pkg/proto/
