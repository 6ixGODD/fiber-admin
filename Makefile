# Go parameters
GO_CMD=go
GO_BUILD=$(GOCMD) build
GO_CLEAN=$(GOCMD) clean
GO_TEST=$(GOCMD) test
GO_GET=$(GOCMD) get
BINARY_NAME=fiber-admin
BINARY_UNIX=$(BINARY_NAME)_unix

# Wire parameters
WIRE_CMD=wire
WIRE_DIR=internal/wire

# Swagger parameters
SWAGGER_CMD=swag
SWAGGER_DIR=./docs

# Docker parameters
DOCKER_CMD=docker
DOCKER_IMAGE_NAME=fiber-admin
DOCKER_IMAGE_TAG=latest

# Build the project
build:
	$(GO_BUILD) -o $(BINARY_NAME) -v

# Build the project for linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO_BUILD) -o $(BINARY_UNIX) -v

# Clean the project
clean:
	$(GO_CLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(WIRE_DIR)/wire_gen.go
	rm -f $(SWAGGER_DIR)/docs.go
	rm -f $(SWAGGER_DIR)/swagger.json
	rm -f $(SWAGGER_DIR)/swagger.yaml

# Run the project
run:
	$(GO_BUILD) -o $(BINARY_NAME) -v
	./$(BINARY_NAME)

# Generate wire
wire:
	$(WIRE_CMD) ./$(WIRE_DIR)

# Generate swagger
swagger:
	$(SWAGGER_CMD) init -g ./$(SWAGGER_DIR)/docs.go

# Generate swagger json
swagger-json:
	$(SWAGGER_CMD) generate spec -o ./$(SWAGGER_DIR)/swagger.json

# Generate swagger yaml
swagger-yaml:
	$(SWAGGER_CMD) generate spec -o ./$(SWAGGER_DIR)/swagger.yaml

# Build docker image
docker-build:
	$(DOCKER_CMD) build -t $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG) .

# Run docker container
docker-run:
	$(DOCKER_CMD) run -p 8080:8080 $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Push docker image
docker-push:
	$(DOCKER_CMD) push $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Remove docker image
docker-rm:
	$(DOCKER_CMD) rmi $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Remove docker container
docker-stop:
	$(DOCKER_CMD) stop $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

# Remove docker container
docker-rm:
	$(DOCKER_CMD) rm $(DOCKER_IMAGE_NAME):$(DOCKER_IMAGE_TAG)

