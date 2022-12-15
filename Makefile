GO							?= @go

PACKAGES					?= ./...
TEST_PACKAGES				?= $(PACKAGES)
COVER_PACKAGES				?= $(shell echo $(TEST_PACKAGES) | tr " " ",")
COVER_PROFILE				?= coverage.out

GOLANGCI_LINT				?= @golangci-lint

DOCKER_IMAGE_REGISTRY		?= docker.io/library
DOCKER_IMAGE_REPOSITORY		?= lindenhoney/linden-honey-scraper-go
DOCKER_IMAGE_TAG			?= latest
DOCKER_IMAGE				?= $(DOCKER_IMAGE_REGISTRY)/$(DOCKER_IMAGE_REPOSITORY):$(DOCKER_IMAGE_TAG)

.PHONY: all
all: build test

.PNONY: fmt
fmt:
	$(GO) fmt $(PACKAGES)

.PHONY: mod/download
mod/download:
	$(GO) mod download

.PHONY: prepare
prepare: mod/download fmt

.PHONY: run
run: prepare
	$(GO) run -v ./cmd/server

.PHONY: build
build: prepare
	$(GO) build -v $(PACKAGES)

.PHONY: install
install: prepare
	$(GO) install -v $(PACKAGES)

.PHONY: test
test: prepare
	$(GO) test -v -race -coverpkg $(COVER_PACKAGES) -coverprofile $(COVER_PROFILE) $(TEST_PACKAGES)

.PHONY: coverage
coverage: test
	$(GO) tool cover -func $(COVER_PROFILE) -o coverage.txt
	$(GO) tool cover -html $(COVER_PROFILE) -o coverage.html

.PHONY: lint
lint: prepare
	$(GOLANGCI_LINT) run -v

.PHONY: docker/build
docker/build:
	@docker build -t $(DOCKER_IMAGE) .

.PHONY: docker/push
docker/push:
	@docker push $(DOCKER_IMAGE)
