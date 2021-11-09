BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)

LDFLAGS = -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

VAULT_ADDR?='http://127.0.0.1:8200'
DATABASE_ADDRESS?=bolt://localhost:7687

# The following variables are used to generate a vault token for the app. The reason for declaring variables, is that
# its difficult to move the token code in a Makefile action. Doing so makes the Makefile more difficult to
# read and starts introduction if/else statements.
VAULT_POLICY:="$(shell vault policy write -address=$(VAULT_ADDR) write-psk policy.hcl)"
TOKEN_INFO:="$(shell vault token create -address=$(VAULT_ADDR) -policy=write-psk -period=24h -display-name=dp-cantabular-metadata-exporter)"
APP_TOKEN:="$(shell echo $(TOKEN_INFO) | awk '{print $$6}')"

.PHONY: all
all: audit test build

.PHONY: audit
audit:
	go list -json -m all | nancy sleuth --exclude-vulnerability-file ./.nancy-ignore

.PHONY: build
build:
	go build -tags 'production' $(LDFLAGS) -o $(BINPATH)/dp-cantabular-metadata-exporter

.PHONY: debug
debug:
	go build -tags 'debug' $(LDFLAGS) -o $(BINPATH)/dp-cantabular-metadata-exporter
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-cantabular-metadata-exporter
	HUMAN_LOG=1 VAULT_TOKEN=$(APP_TOKEN) VAULT_ADDR=$(VAULT_ADDR) DEBUG=1 $(BINPATH)/dp-cantabular-metadata-exporter

PHONY: debug-run
debug-run:
	VAULT_TOKEN=$(APP_TOKEN) HUMAN_LOG=1 DEBUG=1 go run -tags 'debug' $(LDFLAGS) main.go 
	
.PHONY: lint
lint:
	exit

.PHONY: test
test:
	go test -race -cover ./...

.PHONY: convey
convey:
	goconvey ./...

.PHONY: test-component
test-component:
	go test -cover -coverpkg=github.com/ONSdigital/dp-cantabular-metadata-exporter/... -component

	.PHONY: vault
vault:
	@echo "$(VAULT_POLICY)"
	@echo "$(TOKEN_INFO)"
	@echo "$(APP_TOKEN)"
