BINARY_PATH = build/webserver
MAIN_FILE = ./cmd/webserver

.PHONY: migrate
migrate:
	REMANA_APP_ENV=development go run ./cmd/migrate

.PHONY: run/dev
run/dev:
	air

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
	sqlc generate

.PHONY: test/unit
test/unit:
	go test -count=1 -tags unit ./...

.PHONY: test/integration
test/integration:
	go test -count=1 -tags integration ./...

.PHONY: test/e2e
test/e2e:
	go test -count=1 -tags e2e ./...

.PHONY: test/all
test/all:
	go test -count=1 ./...
