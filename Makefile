LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

GO := GO111MODULE=on go

.PHONY: init
init:
	go get -u golang.org/x/lint/golint
	go get -u golang.org/x/tools/cmd/goimports
	@echo "Install pre-commit hook"
	@ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit || true
	@chmod +x ./hack/check.sh

.PHONY: setup
setup: init
	git init

.PHONY: check
check:
	@./hack/check.sh ${scope}

.PHONY: ci
ci: init
	@$(GO) mod tidy && $(GO) mod vendor

.PHONY: build
build: check
	$(GO) build -ldflags '$(LDFLAGS)'

.PHONY: install
install: check
	@echo "Installing arun..."
	@$(GO) install -ldflags '$(LDFLAGS)'

.PHONY: release
release: check
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin/aun
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux/arun
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/windows/aun.exe

.PHONY: docker-image
docker-image:
	docker build -t ahuigo/arun:v0.1.1 -f ./Dockerfile .

.PHONY: push-docker-image
push-docker-image:
	docker push ahuigo/arun:v0.1.1
