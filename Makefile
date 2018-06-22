GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
GOGET=$(GOCMD) get
GOCLEAN=$(GOCMD) clean
BINARY_NAME=cloud-config-builder
BINARY_UNIX=$(BINARY_NAME)_unix

all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./
run: build
	./$(BINARY_NAME)
test:
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)


# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v ./
docker-build:
	docker run --rm -it -e CGO_ENABLED=0 -e GOOS=linux -e GOARCH=amd64 -w /go/src/github.com/mixtore/cloud-config-builder -v "$(PWD)":/go/src/github.com/mixtore/cloud-config-builder golang:latest sh -c "$(GOGET) -v && $(GOBUILD) -o $(BINARY_UNIX) -v ./"
