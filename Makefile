VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := bitc

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/bitc

install:
	go install $(LDFLAGS) ./cmd/bitc

clean:
	rm -f $(BINARY)
