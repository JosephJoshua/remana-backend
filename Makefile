BINARY_PATH = build/webserver
MAIN_FILE = ./cmd/webserver

.PHONY: run/dev
run/dev: build
	CGO_ENABLED=1 REMANA_APP_ENV=development ./${BINARY_PATH}

.PHONY: build
build:
	go build -o ${BINARY_PATH} ${MAIN_FILE}

.PHONY: clean
clean:
	go clean
	rm -f ${BINARY_PATH}
	rm -f ${BINARY_PATH}-darwin
	rm -f ${BINARY_PATH}-linux
	rm -f ${BINARY_PATH}-windows

.PHONY: dep
dep:
	go mod download

.PHONY: vet
vet:
	go vet

.PHONY: tidy
tidy:
	go fmt ./...
	go mod tidy -v

.PHONY: vulncheck
vulncheck:
	govulncheck ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: lint-fix
lint-fix:
	golangci-lint run --fix

.PHONY: generate
generate:
	go generate ./...

.PHONY: test
test:
	go test -count=1 ./...
