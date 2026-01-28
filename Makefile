.PHONY: all deps verify build run test lint clean release-local release-check install uninstall

BINARY_NAME=agen
PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

all: deps verify lint test build

deps:
	go mod download

verify:
	go mod verify

build:
	go build -v -o $(BINARY_NAME) ./cmd/agen

run: build
	./$(BINARY_NAME)

test:
	go test -v ./...

lint:
	go vet ./...

clean:
	go clean
	rm -f $(BINARY_NAME)
	rm -rf dist/

release-local:
	goreleaser release --snapshot --clean

release-check:
	goreleaser check

install: build
	install -d $(DESTDIR)$(BINDIR)
	install -m 755 $(BINARY_NAME) $(DESTDIR)$(BINDIR)/$(BINARY_NAME)

uninstall:
	rm -f $(DESTDIR)$(BINDIR)/$(BINARY_NAME)
