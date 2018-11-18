gofiles := $(filter-out ./vendor/%, $(shell find . -iname "*.go" -type f))
gosrcfiles := $(filter-out %_test.go, $(gofiles))
make_bin := /tmp/$(shell basename "$$PWD")-bin

.PHONY: all
all: deps build

.PHONY: build
build: dist/usersd

.PHONY: clean
clean:
	rm -rf $(make_bin)
	rm -f $(coverage_file)
	rm -rf dist/

.PHONY: deps
deps: $(make_bin)/dep
	PATH="$(make_bin):$$PATH" dep ensure -v -update

.PHONY: docs
docs:

.PHONY: docs-ref
docs-ref:
	@echo "See http://localhost:$${GODOC_PORT:-6060}/pkg/github.com/ntrrg/usersd"
	godoc -http ":$${GODOC_PORT:-6060}" -play

dist/usersd: $(gosrcfiles)
	CGO_ENABLED=0 go build -o $@

$(make_bin)/dep:
	mkdir -p $(make_bin)
	wget -cO $@ 'https://storage.nt.web.ve/_/software/linux/dep-0.5.0-linux-amd64' || wget -cO $@ 'https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64'
	chmod +x $@

# Development

coverage_file := coverage.txt

.PHONY: benchmark
benchmark:
	go test -race -bench . -benchmem ./...

.PHONY: build-docker
build-docker:
	docker build -t ntrrg/usersd .

.PHONY: ci
ci: test lint qa coverage benchmark

.PHONY: clean-dev
clean-dev: clean
	rm -rf $(coverage_file)
	rm -rf pkg/usersd/test-db
	docker image rm ntrrg/usersd || true

.PHONY: coverage
coverage: test
	go tool cover -func $(coverage_file)

.PHONY: coverage-web
coverage-web: test
	go tool cover -html $(coverage_file)

.PHONY: deps-dev
deps-dev: $(make_bin)/gometalinter
	go get -u -v golang.org/x/lint/golint

.PHONY: format
format:
	gofmt -s -w -l $(gofiles)

.PHONY: lint
lint:
	gofmt -d -e -s $(gofiles)
	golint ./ ./api/... ./pkg/...

.PHONY: lint-md
lint-md:
	@docker run --rm -itv "$$PWD":/files/ ntrrg/md-linter

.PHONY: qa
qa: $(make_bin)/gometalinter
	PATH="$(make_bin):$$PATH" CGO_ENABLED=0 gometalinter --tests ./ ./api/... ./pkg/...

.PHONY: test
test:
	go test -race -coverprofile $(coverage_file) -v ./...

$(make_bin)/gometalinter:
	mkdir -p $(make_bin)
	wget -cO /tmp/gometalinter.tar.gz 'https://storage.nt.web.ve/_/software/linux/gometalinter-2.0.11-linux-amd64.tar.gz' || wget -cO /tmp/gometalinter.tar.gz 'https://github.com/alecthomas/gometalinter/releases/download/v2.0.11/gometalinter-2.0.11-linux-amd64.tar.gz'
	tar -xf /tmp/gometalinter.tar.gz -C /tmp/
	cp -a $$(find /tmp/gometalinter-2.0.11-linux-amd64/ -type f -executable) $(make_bin)/

