BINARY_NAME = arb-updater
VERSION     = 1.0.0
BUILD_DATE  = $(shell date +%F)
BUILD_FLAGS = "-X main.Version=${VERSION} -X main.Build=${BUILD_DATE}"

.PHONY: build
build:
	go mod tidy
	GOARCH=amd64 CGO_ENABLED=0 \
	go build -mod mod -buildvcs=false -ldflags ${BUILD_FLAGS} \
	-o bin/${BINARY_NAME}

run:
	go mod tidy
	go run -race -mod mod -buildvcs=false -ldflags ${BUILD_FLAGS} . \
	-t ${ARB_TEMPLATE_FILE} -l ${ARB_LOCALE_FILE}

lint:
	golangci-lint run --verbose

benchmark:
	go mod tidy
	go test -bench=.

.PHONY: test
test:
	go mod tidy
	go test -race -v ./...

clean:
	go clean
	go clean -testcache
	@if [ -d bin ] ; then rm -rf bin ; fi

help:
	@echo "make build VERSION=1.0.0 - compile the binary file with golang codes"
	@echo "make clean - clean cache remove the binary file in the bin directory"
	@echo "make lint - check golang syntax"
	@echo "make benchmark - run benchmark"
	@echo "make test - run test with -race parameter"
	@echo "make run ARB_TEMPLATE_FILE=path/to/xx.arb ARB_LOCALE_FILE=path/to/xx.arb - run the service with arb files"
