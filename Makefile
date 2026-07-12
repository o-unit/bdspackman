.PHONY: build release clean

APP := bdspackman
VERSION ?= dev
OUTDIR := build/$(VERSION)

build:
	mkdir -p $(OUTDIR)

	go build -ldflags="-s -w" \
		-o $(OUTDIR)/$(APP) \
		./cmd/bdspackman

release:
	mkdir -p $(OUTDIR)

	GOOS=linux GOARCH=amd64 \
	go build -ldflags="-s -w" \
		-o $(OUTDIR)/$(APP)-linux-amd64 \
		./cmd/bdspackman

	GOOS=linux GOARCH=arm64 \
	go build -ldflags="-s -w" \
		-o $(OUTDIR)/$(APP)-linux-arm64 \
		./cmd/bdspackman

clean:
	rm -rf build