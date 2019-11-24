pkgName := usersd
hugoPort := 1313
godocPort := 6060

goAllFiles := $(filter-out ./vendor/%, $(shell find . -iname "*.go" -type f))
goSrcFiles := $(shell go list -f "{{ \$$path := .Dir }}{{ range .GoFiles }}{{ \$$path }}/{{ . }} {{ end }}" ./...)
goTestFiles := $(shell go list -f "{{ \$$path := .Dir }}{{ range .TestGoFiles }}{{ \$$path }}/{{ . }} {{ end }}" ./...)

.PHONY: all
all: build

.PHONY: build
build: dist/$(pkgName)-$(shell go env "GOOS")-$(shell go env "GOARCH")

.PHONY: build-all
build-all:
	$(MAKE) -s build-darwin-386
	$(MAKE) -s build-darwin-amd64
	$(MAKE) -s build-linux-386
	$(MAKE) -s build-linux-amd64
	$(MAKE) -s build-linux-arm
	$(MAKE) -s build-linux-arm64
	$(MAKE) -s build-windows-386
	$(MAKE) -s build-windows-amd64

.PHONY: build-%
build-%:
	@\
		GOOS="$(shell echo "$*" | cut -d "-" -f 1)" \
		GOARCH="$(shell echo "$*" | cut -sd "-" -f 2)" \
		$(MAKE) -s build

.PHONY: clean
clean:
	rm -rf dist/

.PHONY: doc
doc:
	@echo "http://localhost:$(hugoPort)/en/projects/$(pkgName)/"
	@echo "http://localhost:$(hugoPort)/es/projects/$(pkgName)/"
	@docker run --rm -it \
		-e PORT=$(hugoPort) \
		-p $(hugoPort):$(hugoPort) \
		-v "$$PWD/.ntweb":/site/content/projects/$(pkgName)/ \
		ntrrg/ntweb:editing --port $(hugoPort)

.PHONE: doc-go
doc-go:
	godoc -http :$(godocPort) -play

dist/%: $(goSrcFiles)
	CGO_ENABLED=0 go build -o "dist/$*" .

# Development

coverage_file := coverage.txt

.PHONY: benchmark
benchmark:
	go test -v -bench . -benchmem ./...

.PHONY: ca
ca:
	golangci-lint run

.PHONY: ci
ci: clean-dev test lint ca coverage benchmark build

.PHONY: ci-race
ci-race: clean-dev test-race lint ca coverage benchmark build

.PHONY: clean-dev
clean-dev: clean
	rm -rf $(coverage_file)

.PHONY: coverage
coverage: $(coverage_file)
	go tool cover -func $<

.PHONY: coverage-web
coverage-web: $(coverage_file)
	go tool cover -html $<

.PHONY: format
format:
	gofmt -s -w -l $(goAllFiles)

.PHONY: lint
lint:
	gofmt -d -e -s $(goAllFiles)

.PHONY: test
test:
	go test -v ./...

.PHONY: test-race
test-race:
	go test -race -v ./...

.PHONY: watch
watch:
	reflex -d "none" -r '\.go$$' -- $(MAKE) -s test lint

.PHONY: watch-race
watch-race:
	reflex -d "none" -r '\.go$$' -- $(MAKE) -s test-race lint

$(coverage_file): $(goSrcFiles) $(goTestFiles)
	go test -coverprofile $(coverage_file) ./...

