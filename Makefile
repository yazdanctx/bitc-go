VERSION ?= 0.2.0
LDFLAGS := -ldflags "-X github.com/yazdanctx/bitc-go/internal/version.Version=$(VERSION)"
BINARY := bitc

.PHONY: build install clean tag

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/bitc

install:
	cp $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)

tag:
	git tag -a v$(VERSION) -m "v$(VERSION)"
	git push origin v$(VERSION)
