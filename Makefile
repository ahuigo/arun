LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

msg?=
GO := GO111MODULE=on go


.PHONY: init
init:
	go install golang.org/x/lint/golint@latest
	go install golang.org/x/tools/cmd/goimports@latest
	@echo "Install pre-commit hook"
	@ln -s $(shell pwd)/hooks/pre-commit $(shell pwd)/.git/hooks/pre-commit || true
	@chmod +x ./hack/check.sh
	go mod tidy

.PHONY: build
build: 
	$(GO) build -ldflags '$(LDFLAGS)'

.PHONY: install
install: init
	@echo "Installing arun..."
	@$(GO) install -ldflags '$(LDFLAGS)'

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin/aun
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux/arun
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/windows/aun.exe

.PHONY: docker-image
docker-image:
	docker build -t ahuigo/arun:`cat version` -f ./Dockerfile .

.PHONY: push-docker-image
push-docker-image:
	docker push ahuigo/arun:`cat version`


pkg:
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"
	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v" && git push origin HEAD

release: 
	#goreleaser init
	goreleaser release --snapshot --rm-dist
