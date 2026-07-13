.PHONY: build release clean

APP := bdspackman
VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD)
OUTDIR := build/$(VERSION)

build:
	mkdir -p $(OUTDIR)

	go build -ldflags="-s -w -X main.version=$(VERSION)" \
		-o $(OUTDIR)/$(APP) \
		./cmd/bdspackman

release:
	mkdir -p $(OUTDIR)/linux-amd64
	mkdir -p $(OUTDIR)/linux-arm64

	GOOS=linux GOARCH=amd64 \
	go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)" \
		-o $(OUTDIR)/linux-amd64/$(APP) \
		./cmd/bdspackman

	GOOS=linux GOARCH=arm64 \
	go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=$(COMMIT)" \
		-o $(OUTDIR)/linux-arm64/$(APP) \
		./cmd/bdspackman

clean:
	rm -rf build