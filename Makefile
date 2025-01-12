BIN_DIR=bin
BINARY_NAME=go-index
APP_SOURCE_CODE=cmd/${BINARY_NAME}/main.go

ARCH_ARM64=arm64
ARCH_AMD64=amd64
OS_DARWIN=darwin
OS_LINUX=linux
OS_WINDOWS=windows

.PHONY: build build-all run clean test test-coverage dep lint

build: ${BIN_DIR}/${BINARY_NAME}
	
${BIN_DIR}/${BINARY_NAME}: 
	@if [ "$$(uname)" = "Darwin" ]; then \
		if [ "$$(uname -m)" = "arm64" ]; then \
			GOARCH=$(ARCH_ARM64) GOOS=$(OS_DARWIN) go build -o ${BIN_DIR}/${BINARY_NAME}-darwin-arm64 ${APP_SOURCE_CODE}; \
		else \
			GOARCH=$(ARCH_AMD64) GOOS=$(OS_DARWIN) go build -o ${BIN_DIR}/${BINARY_NAME}-darwin-amd64 ${APP_SOURCE_CODE}; \
		fi \
	elif [ "$$(uname)" = "Linux" ]; then \
		GOARCH=$(ARCH_AMD64) GOOS=$(OS_LINUX) go build -o ${BIN_DIR}/${BINARY_NAME}-linux-amd64 ${APP_SOURCE_CODE}; \
	elif [ "$$(uname)" = "MINGW64_NT-10.0" ]; then \
		GOARCH=$(ARCH_AMD64) GOOS=$(OS_WINDOWS) go build -o ${BIN_DIR}/${BINARY_NAME}-windows-amd64.exe ${APP_SOURCE_CODE}; \
	else \
		echo "Unsupported OS"; \
	fi

build-all:
	GOARCH=$(ARCH_ARM64) GOOS=$(OS_DARWIN) go build -o ${BIN_DIR}/${BINARY_NAME}-darwin-arm64 ${APP_SOURCE_CODE}
	GOARCH=$(ARCH_AMD64) GOOS=$(OS_DARWIN) go build -o ${BIN_DIR}/${BINARY_NAME}-darwin-amd64 ${APP_SOURCE_CODE}
	GOARCH=$(ARCH_AMD64) GOOS=$(OS_LINUX) go build -o ${BIN_DIR}/${BINARY_NAME}-linux-amd64 ${APP_SOURCE_CODE}
	GOARCH=$(ARCH_AMD64) GOOS=$(OS_WINDOWS) go build -o ${BIN_DIR}/${BINARY_NAME}-windows-amd64.exe ${APP_SOURCE_CODE}

run: build
	@if [ "$$(uname)" = "Darwin" ]; then \
		if [ "$$(uname -m)" = "arm64" ]; then \
			./${BIN_DIR}/${BINARY_NAME}-darwin-arm64; \
		else \
			./${BIN_DIR}/${BINARY_NAME}-darwin-amd64; \
		fi \
	elif [ "$$(uname)" = "Linux" ]; then \
		./${BIN_DIR}/${BINARY_NAME}-linux-amd64; \
	elif [ "$$(uname)" = "MINGW64_NT-10.0" ]; then \
		./${BIN_DIR}/${BINARY_NAME}-windows-amd64.exe; \
	else \
		echo "Unsupported OS"; \
	fi

test:
	go test ./...

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

dep:
	go mod download

lint:
	golangci-lint run --enable-all

fmt:
	go fmt ./...

clean:
	go clean
	rm -f ${BIN_DIR}/${BINARY_NAME}-darwin-arm64
	rm -f ${BIN_DIR}/${BINARY_NAME}-darwin-amd64
	rm -f ${BIN_DIR}/${BINARY_NAME}-linux-amd64
	rm -f ${BIN_DIR}/${BINARY_NAME}-windows-amd64.exe
	rm -f coverage.out
	rm -f coverage.html
