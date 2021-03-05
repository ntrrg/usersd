module := $(shell go list -m)
PACKAGE ?= $(notdir $(module))
HUGO_PORT ?= 1313
GODOC_PORT ?= 6060

goAllFiles := $(shell find . -iname "*.go" -type f)
goFiles := $(filter-out ./vendor/%, $(goAllFiles))
goSrcFiles := $(shell go list -f '{{ range .GoFiles }}{{ $$.Dir }}/{{ . }} {{ end }}' ./...)
goTestFiles := $(shell go list -f "{{ range .TestGoFiles }}{{ $$.Dir }}/{{ . }} {{ end }}{{ range .XTestGoFiles }}{{ $$.Dir }}/{{ . }} {{ end }}" ./...)

.PHONY: all
all: build

.PHONY: build
build:
	go build ./...

.PHONY: clean
clean: clean-dev
	rm -rf dist/

.PHONY: doc
doc:
	@echo "Go to http://localhost:$(HUGO_PORT)/en/projects/$(PACKAGE)/"
	@docker run --rm -it \
		-e PORT=$(HUGO_PORT) \
		-p $(HUGO_PORT):$(HUGO_PORT) \
		-v "$$PWD/.ntweb":/site/content/projects/$(PACKAGE)/ \
		ntrrg/ntweb:editing --port $(HUGO_PORT)

.PHONE: godoc
godoc:
	@echo "Go to http://localhost:$(GODOC_PORT)/pkg/$(module)/"
	godoc -http :$(GODOC_PORT) -play

# Development

COVERAGE_FILE ?= coverage.txt
TARGET_FUNC ?= .
TARGET_PKG ?= ./...

.PHONY: benchmark
benchmark:
	go test -run none -bench "$(TARGET_FUNC)" -benchmem -v $(TARGET_PKG)

.PHONY: ca
ca:
	golangci-lint run

.PHONY: ca-fast
ca-fast:
	golangci-lint run --fast

.PHONY: ci
ci: test lint ca coverage build

.PHONY: ci-race
ci-race: test-race lint ca coverage build

.PHONY: clean-dev
clean-dev: clean
	rm -rf $(COVERAGE_FILE)

.PHONY: coverage
coverage:
	go tool cover -func $(COVERAGE_FILE)

.PHONY: coverage-web
coverage-web:
	go tool cover -html $(COVERAGE_FILE)

.PHONY: format
format:
	gofmt -s -w -l $(goFiles)

.PHONY: lint
lint:
	gofmt -d -e -s $(goFiles)

.PHONY: test
test:
	go test \
		-run "$(TARGET_FUNC)" \
		-coverprofile $(COVERAGE_FILE) \
		-v $(TARGET_PKG)

.PHONY: test-race
test-race:
	go test \
		-run "$(TARGET_FUNC)" \
		-coverprofile $(COVERAGE_FILE) \
		-race -v $(TARGET_PKG)

.PHONY: watch
watch:
	reflex -d "none" -r '\.go$$' -- $(MAKE) -s build test lint

.PHONY: watch-race
watch-race:
	reflex -d "none" -r '\.go$$' -- $(MAKE) -s build test-race lint
