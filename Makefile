LDFLAGS += -X "main.BuildTimestamp=$(shell date -u "+%Y-%m-%d %H:%M:%S")"
LDFLAGS += -X "main.Version=$(shell git rev-parse HEAD)"

msg?=

.PHONY: build
build: 
	$(GO) build -ldflags '$(LDFLAGS)'

.PHONY: install
install:
	@echo "Installing arun..."
	@$(GO) install -ldflags '$(LDFLAGS)'

.PHONY: release
release:
	GOOS=darwin GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/darwin/aun
	GOOS=linux GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/linux/arun
	GOOS=windows GOARCH=amd64 $(GO) build -ldflags '$(LDFLAGS)' -o bin/windows/aun.exe

.PHONY: docker-image
docker-image:
	docker build -t ahuigo/arun:v0.1.1 -f ./Dockerfile .

.PHONY: push-docker-image
push-docker-image:
	docker push ahuigo/arun:v0.1.1


pkg:
	{ hash newversion.py 2>/dev/null && newversion.py version;} ;  { echo version `cat version`; }
	git commit -am "$(msg)"

	#jfrog "rt" "go-publish" "go-pl" $$(cat version) "--url=$$GOPROXY_API" --user=$$GOPROXY_USER --apikey=$$GOPROXY_PASS
	v=`cat version` && git tag "$$v" && git push origin "$$v"

