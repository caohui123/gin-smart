# This how we want to name the binary output
BINARY=gin-smart
BIN_PATH=bin
MAIN_FILE=cmd/app/main.go
# These are the values we want to pass for VERSION  and BUILD
VERSION=1.1.0
BUILD=`date +%Y-%m-%d^%H:%M:%S`
# Setup the -Idflags options for go build here,interpolate the variable values
LDFLAGS=-ldflags "-X main.BuildVersion=${VERSION} -X main.BuildAt=${BUILD}"
# Builds the project
build:
	go build ${LDFLAGS} -o ${BIN_PATH}/${BINARY} ${MAIN_FILE}
# build windows
build-win:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build  ${LDFLAGS} -o ${BINPATH}/${BINARY}.exe ${MAIN_FILE}
# build linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o ${BINPATH}/${BINARY} ${MAIN_FILE}
# docker build (see docker-build.md)
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

test: ## Run unittests
	@export CGO_CFLAGS_ALLOW="-maes" && go test -v  -parallel 1  $(PKG_LIST)

fmt: ## go fmt
	@go fmt ./...
	@gofumpt -s -w ./..
	@#gofumports -s -w ./..


.PHONY:  clean install


checkgofmt: fmt  ## get all go files and run go fmt on them
	@files=$$(git status -suno);if [ -n "$$files" ]; then \
		  echo "Error: 'make fmt' needs to be run on:"; \
		  echo "$${files}"; \
		  exit 1; \

