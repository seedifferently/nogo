VERSION := 1.0.0-beta.1
BINARY := nogo
SHELL := /bin/bash
LDFLAGS := "-X main.version=$(VERSION) -X main.build=`git rev-parse --verify --short HEAD`"
GOX_OSARCH := "darwin/amd64 freebsd/386 freebsd/amd64 freebsd/arm linux/386 linux/amd64 linux/arm64 netbsd/386 netbsd/amd64 netbsd/arm openbsd/386 openbsd/amd64 windows/386 windows/amd64"
GOX_OUTPUT := "build/$(BINARY)_v$(VERSION)_{{.OS}}_{{.Arch}}/$(BINARY)"

.DEFAULT_GOAL := $(BINARY)
$(BINARY):
	go build -ldflags $(LDFLAGS) -o $(BINARY)

build: # https://github.com/mitchellh/gox
	CGO_ENABLED=0 gox -ldflags $(LDFLAGS) -osarch $(GOX_OSARCH) -output $(GOX_OUTPUT)
	# gox doesn't support fine-tuning arm (GOARM defaults to 6), so build linux/arm individually
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 go build -ldflags $(LDFLAGS) -o build/$(BINARY)_v$(VERSION)_linux_armv5/$(BINARY)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=6 go build -ldflags $(LDFLAGS) -o build/$(BINARY)_v$(VERSION)_linux_armv6/$(BINARY)
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -ldflags $(LDFLAGS) -o build/$(BINARY)_v$(VERSION)_linux_armv7/$(BINARY)

dist: build
	@mkdir dist
	$(eval BUILDS := $(shell find build/ -type f -printf '%P\n'))
	@for f in $(BUILDS); do \
		echo "Archiving $${f%/*}..."; \
		if [[ $$f =~ darwin|windows ]]; then \
			zip -j dist/$${f%/*}.zip build/$$f README.md; \
		else \
			tar -cvzf dist/$${f%/*}.tar.gz README.md -C build/$${f%/*} $${f#*/}; \
		fi; \
	done

.PHONY: deps
deps:
	go get github.com/miekg/dns
	go get github.com/boltdb/bolt
	go get github.com/pressly/chi

.PHONY: clean
clean:
	rm -rf build/
	rm -rf dist/
	go clean