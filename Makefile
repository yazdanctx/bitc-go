VERSION ?= 1.0.0
LDFLAGS := -ldflags "-X main.version=$(VERSION)"
BINARY := bitc

.PHONY: build install clean

build:
	go build $(LDFLAGS) -o $(BINARY) ./cmd/bitc

install:
	cp $(BINARY) /usr/local/bin/$(BINARY)

clean:
	rm -f $(BINARY)

tag:
	git tag -a v$(VERSION) -m "v$(VERSION)"
	git push origin v$(VERSION)
