VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := compressor

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/compressor

install:
	go install $(LDFLAGS) ./cmd/compressor

clean:
	rm -f $(BINARY)
