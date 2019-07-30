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
CURRENT_DIR=$(shell pwd)

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-w -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.Branch=${BRANCH}"

# Build the project
all: clean test vet linux docker

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} . ; \

darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} . ; \

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH}.exe . ; \

#test:
#	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
#	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \

vet:
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \

fmt:
	go fmt $$(go list ./... | grep -v /vendor/) ; \

docker: linux
	docker build -t ${GITHUB_USERNAME}/${BINARY} .

tag: tag-latest tag-version

tag-version: docker
	docker tag ${GITHUB_USERNAME}/${BINARY} ${GITHUB_USERNAME}/${BINARY}:${VERSION}

tag-latest: docker
	docker tag ${GITHUB_USERNAME}/${BINARY} ${GITHUB_USERNAME}/${BINARY}:latest

publish: publish-latest publish-version

publish-version: tag
	docker push ${GITHUB_USERNAME}/${BINARY}:${VERSION}

publish-latest: tag
	docker push ${GITHUB_USERNAME}/${BINARY}:latest

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: linux darwin windows test vet fmt docker tag tag-version tag-latest publish publish-version publish-latest clean
