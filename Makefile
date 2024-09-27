TIMESTAMP=tmp-$(shell date +%s )
ARGOBOT_VERSION ?= $(shell cat version)
ARGOBOT_COMMIT ?= $(shell git rev-parse HEAD)

expose:
	ngrok http 8080

build:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X 'main.version=$(ARGOBOT_VERSION)' -X 'main.commit=$(ARGOBOT_COMMIT)'" -v -o argobot ./cmd/argobot

.PHONY: version
version:
	echo $(ARGOBOT_VERSION)

.PHONY: test
test:
	go test ./pkg/...

test-helm:
	go clean -testcache;
	go test ./test/integration/

build-argocli:
	docker build -t ghcr.io/corymurphy/argocli --target argocli .

.PHONY: dev-install
dev-install:
	go clean -testcache
	go test -run Test_LocalDevelopmentInstall ./test/dev
