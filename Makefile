.PHONY: build-server build-client run-server run-client test-server test-client test-all test-unit test-integration

build-server:
	@cd server && go build -o bin/server src/main.go

build-client:
	@cd client && go build -o bin/client src/main.go

run-server:
	@cd server && go run src/main.go

run-client:
	@cd client && go run src/main.go

test-server-unit:
	@cd server && go test -v ./src/gateways ./src/handlers ./src/repositories

test-server-integration:
	@cd server && go test -v ./src/tests/integration

test-client-unit:
	@cd client && go test -v ./src/usecases

test-client-integration:
	@cd client && go test -v ./src/tests/integration

test-unit: test-server-unit test-client-unit

test-integration: test-server-integration test-client-integration

test-all: test-unit test-integration 