.PHONY: build test docker-build run lint

BINARY_NAME := server
CMD_PATH    := ./cmd
IMAGE_NAME := app

# build compiles the service binary
build:
	go build -o $(BINARY_NAME) $(CMD_PATH)

# test runs the full unit test suite
test:
	go test ./...

# docker-build builds the Docker image
docker-build:
	docker build -t ${IMAGE_NAME} .

# run starts the service using go run (for local development)
run:
	go run $(CMD_PATH)

# lint runs golangci-lint
lint:
	golangci-lint run ./...
