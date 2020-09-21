# Borrowed from:
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

PWD = $(shell pwd)
BINARY = $(shell basename ${PWD})
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION=$(shell cat VERSION)
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=chremoas
DEV_REGISTRY=docker.4amlunch.net

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

# Build the project
all: clean docker

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -mod=vendor ${LDFLAGS} -o ${BINARY}-linux-amd64 . ; \

illumos:
	CGO_ENABLED=0 GOOS=illumos GOARCH=${GOARCH} go build -mod=vendor ${LDFLAGS} -o ${BINARY}-illumos-${GOARCH} . ; \

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build -mod=vendor ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build -mod=vendor ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \

#test:
#	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
#	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \

vet:
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \

fmt:
	go fmt $$(go list ./... | grep -v /vendor/) ; \

docker:
	docker buildx build --build-arg VERSION=${VERSION} --build-arg COMMIT=${COMMIT} --build-arg BRANCH=${BRANCH} --build-arg BINARY=${BINARY} --platform=linux/amd64,linux/arm64 -t ${GITHUB_USERNAME}/${BINARY}:${VERSION} -t ${GITHUB_USERNAME}/${BINARY}:latest --push .

install-illumos: illumos
	cp ${BINARY}-illumos-${GOARCH} /usr/local/bin/${BINARY}
	svccfg import smf.xml

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: linux illumos darwin windows test fmt docker clean
