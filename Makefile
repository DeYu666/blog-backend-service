Go=go

TARGETS := blog-backend-service
BUILD := $(shell git rev-parse --short HEAD)
VERSION := $(shell git rev-parse --abbrev-ref HEAD)-$(BUILD)

IMAGE_NAME := blog-backend-service
IMAGE_VERSION := $(IMAGE_NAME):$(VERSION)

REGISTRY_ADDRESS ?= registry.hub.docker.com
IMAGE_FULLNAME := $(REGISTRY_ADDRESS)/deyu666/$(IMAGE_VERSION)

LDFLAGS += -X "$(project)/version.BuildTS=$(shell date -u '+%Y-%m-%d %I:%M:%S')"
LDFLAGS += -X "$(project)/version.GitHash=$(shell git rev-parse HEAD)"
LDFLAGS += -X "$(project)/version.GitBranch=$(shell git rev-parse --abbrev-ref HEAD)"

build:
	go build -o $(TARGETS) main.go

image: $(TARGETS)
	cp -f Dockerfile.j2 Dockerfile
	sed -i'' -e "s/{{REGISTRY_ADDRESS}}/$(REGISTRY_ADDRESS)/g" Dockerfile
	sed -i'' -e "s/{{IMAGE_NAME}}/$(IMAGE_NAME)/g" Dockerfile
	sed -i'' -e "s/{{EXE_NAME}}/$(TARGETS)/g" Dockerfile
	docker build -t $(IMAGE_FULLNAME) .

image-push:
	docker tag $(IMAGE_FULLNAME) deyu666/$(IMAGE_NAME):latest
	docker push deyu666/$(IMAGE_NAME):latest

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags '$(LDFLAGS)' $(project)

lint:
	#@gometalinter --disable-all --enable=gosec --exclude="Errors unhandled"  ./...
	@golangci-lint run --deadline=5m

packages = $(shell go list ./...|grep -v /vendor/)
test: check
	$(GO) test ${packages}
