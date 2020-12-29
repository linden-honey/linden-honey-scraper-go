GO						:= go

PACKAGES				:= ./...
GO_COVER_PROFILE		:= coverage.out
GOLANGCI_LINT_VERSION	:= v1.29.0

.PHONY: all
all: build test

.PNONY: fmt
fmt:
	${GO} fmt $(PACKAGES)

.PHONY: deps
deps:
	${GO} mod tidy -v

.PHONY: prepare
prepare: deps fmt

.PHONY: build
build: prepare
	${GO} build -v $(PACKAGES)

.PHONY: install
install: prepare
	${GO} install -v $(PACKAGES)

.PHONY: test
test: prepare
	${GO} test -v -race -coverprofile=$(GO_COVER_PROFILE) $(PACKAGES)

.PHONY: coverage
coverage: test
	${GO} tool cover -func=$(GO_COVER_PROFILE) -o coverage.txt
	${GO} tool cover -html=$(GO_COVER_PROFILE) -o coverage.html

.PHONY: lint
lint: prepare
	docker run --rm -v $(CURDIR):/app -w /app golangci/golangci-lint:${GOLANGCI_LINT_VERSION} golangci-lint run -v