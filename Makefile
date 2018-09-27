# Borrowed from: 
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY = lucky
VET_REPORT = vet.report
TEST_REPORT = tests.xml
GOARCH = amd64

VERSION?=?
COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)

# Symlink into GOPATH
GITHUB_USERNAME=blocktop
BUILD_DIR=${GOPATH}/src/github.com/${GITHUB_USERNAME}/go-lucky
CURRENT_DIR=$(shell pwd)
BUILD_DIR_LINK=$(shell readlink ${BUILD_DIR})

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X main.VERSION=${VERSION} -X main.COMMIT=${COMMIT} -X main.BRANCH=${BRANCH}"

# Build the project
all: link clean deps test vet linux darwin windows

link:
	BUILD_DIR=${BUILD_DIR}; \
	BUILD_DIR_LINK=${BUILD_DIR_LINK}; \
	CURRENT_DIR=${CURRENT_DIR}; \
	if [ "$${BUILD_DIR_LINK}" != "$${CURRENT_DIR}" ]; then \
	    echo "Fixing symlinks for build"; \
	    rm -f $${BUILD_DIR}; \
	    ln -s $${CURRENT_DIR} $${BUILD_DIR}; \
	fi

linux: deps
	cd ${BUILD_DIR}; \
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -a -o ${BINARY}-linux-${GOARCH} . ; \
	cd - >/dev/null

darwin: deps
	cd ${BUILD_DIR}; \
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -a -o ${BINARY}-darwin-${GOARCH} . ; \
	cd - >/dev/null

windows: deps
	cd ${BUILD_DIR}; \
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -a -o ${BINARY}-windows-${GOARCH}.exe . ; \
	cd - >/dev/null

test:
	if ! hash go2xunit 2>/dev/null; then go install github.com/tebeka/go2xunit; fi
	cd ${BUILD_DIR}; \
	godep go test -v ./... 2>&1 | go2xunit -output ${TEST_REPORT} ; \
	cd - >/dev/null

vet:
	-cd ${BUILD_DIR}; \
	godep go vet ./... > ${VET_REPORT} 2>&1 ; \
	cd - >/dev/null

deps:
	-cd ${BUILD_DIR} ; \
  glide install ; \
	#curl -L -o go-multiaddr-master.zip https://github.com/multiformats/go-multiaddr/archive/master.zip ; \
	#unzip -o go-multiaddr-master.zip ; \
	#rm -rf ./vendor/github.com/multiformats/go-multiaddr ; \
	#mv go-multiaddr-master ./vendor/github.com/multiformats/go-multiaddr ; \
	#rm go-multiaddr-master.zip ; 
	cd - >/dev/null

fmt:
	cd ${BUILD_DIR}; \
	go fmt $$(go list ./... | grep -v /vendor/) ; \
	cd - >/dev/null

clean:
	-rm -f ${TEST_REPORT}
	-rm -f ${VET_REPORT}
	-rm -f ${BINARY}-*

.PHONY: link linux darwin windows test vet fmt clean