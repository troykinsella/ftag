
PACKAGE=github.com/troykinsella/ftag
BINARY=ftag
VERSION=1.0.0

LDFLAGS=-ldflags "-X main.AppVersion=${VERSION}"

build:
	go build ${LDFLAGS} ${PACKAGE}

install:
	go install ${LDFLAGS}

deps:
	go get -d -v ./...

dev-deps: deps
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega
	go get github.com/mitchellh/gox

test:
	go test ${PACKAGE}/...

coverage:
	go test -cover ${PACKAGE}/...

dist:
	gox ${LDFLAGS} \
		-arch="amd64" \
		-os="darwin linux windows" \
		-output="${BINARY}_{{.OS}}_{{.Arch}}" \
		${PACKAGE}

clean:
	test -f ${BINARY} && rm ${BINARY} || true
	rm ${BINARY}_* || true

.PHONY: build deps dev-deps install test coverage dist release clean
