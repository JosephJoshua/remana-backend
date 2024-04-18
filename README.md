# Remana Backend

Backend API service for Remana.

## Requirements
- [Go v1.22.1](https://go.dev/dl/)
- [golangci-lint](https://github.com/golangci/golangci-lint)
	Set of linters for Go.
- [ogen](https://github.com/ogen-go/ogen)
	OpenAPI v3 code generator.
- [air](https://github.com/cosmtrek/air)
	Provides live reload for Go apps.
- [sqlc](https://github.com/sqlc-dev/sqlc/)
	SQL type-safe code generator for Go.

## Setup
1. Generate self-signed public and private keys.
	 ```
	 openssl genrsa -out server.key 2048
	 openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
	 ```
