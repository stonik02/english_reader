COMPOSE := docker compose -f docker-compose.yml
TEST_COMPOSE := docker compose -f docker-compose.test.yml

ifneq (,$(wildcard .env))
include .env
export
endif

.DEFAULT_GOAL := help

.PHONY: help infra-up infra-down app-up app-down logs ps run build test test-race test-docker test-db-up test-db-down test-repository dictionary-convert dictionary-import book-reprocess mocks fmt vet tidy migrate-up migrate-down proto proto-lint

help: ## Show available commands
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z0-9_-]+:.*##/ {printf "%-16s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

infra-up: ## Start PostgreSQL and LibreTranslate
	$(COMPOSE) up -d postgres libretranslate

infra-down: ## Stop local infrastructure without deleting data
	$(COMPOSE) stop postgres libretranslate

app-up: ## Build and start API and Envoy with infrastructure
	$(COMPOSE) --profile app up --build -d

app-down: ## Stop all Compose services without deleting volumes
	$(COMPOSE) --profile app down

logs: ## Follow API logs from Compose
	$(COMPOSE) logs -f api

ps: ## Show Compose service status
	$(COMPOSE) ps

run: ## Run the API on the host; start infra-up first
	go run ./cmd/server

build: ## Build server and migration binaries into bin/
	mkdir -p bin
	go build -o bin/reader-api ./cmd/server
	go build -o bin/reader-migrate ./cmd/migrate

test: ## Run unit and integration-safe Go tests
	go test ./...

test-race: ## Run Go tests with the race detector
	go test -race ./...

test-docker: ## Run tests in the Go 1.25 API image
	$(COMPOSE) --profile app run --rm --entrypoint go api test ./...

test-db-up: ## Start the isolated PostgreSQL database for repository tests
	$(TEST_COMPOSE) up -d --wait

test-db-down: ## Stop the isolated repository-test database
	$(TEST_COMPOSE) down -v

test-repository: test-db-up ## Run PostgreSQL repository integration tests
	TEST_DATABASE_URL="postgres://reader_test:reader_test@localhost:$${TEST_POSTGRES_PORT:-5434}/reader_test?sslmode=disable" go test ./test/integration/... -count=1

dictionary-import: ## Import JSONL; set DATABASE_URL (or .env), DICTIONARY_FILE and DICTIONARY_VERSION
	go run ./cmd/dictionary-import -file "$(DICTIONARY_FILE)" -version "$(DICTIONARY_VERSION)" -source "$(or $(DICTIONARY_SOURCE),wiktionary)"

book-reprocess: ## Reparse stored EPUBs; set DATABASE_URL and REPROCESS_ALL=true or REPROCESS_BOOK_ID
	go run ./cmd/book-reprocess $(if $(REPROCESS_ALL),-all,-book-id "$(REPROCESS_BOOK_ID)")

dictionary-convert: ## Convert Kaikki English dump; set KAIKKI_FILE, DICTIONARY_FILE and DICTIONARY_VERSION
	python3 tools/convert_kaikki_en_ru.py --input "$(KAIKKI_FILE)" --output "$(DICTIONARY_FILE)" --version "$(DICTIONARY_VERSION)"

mocks: ## Generate GoMock mocks from every dependencies.go
	@command -v mockgen >/dev/null || (echo 'install mockgen: go install go.uber.org/mock/mockgen@v0.6.0' && exit 1)
	@for dependency in internal/usecase/*/*/dependencies.go internal/handler/grpc/*/*/dependencies.go internal/handler/http/*/*/dependencies.go internal/worker/*/dependencies.go; do \
		[ -f "$$dependency" ] || continue; \
		package_dir=$${dependency%/dependencies.go}; \
		package_name=$${package_dir##*/}; \
		mockgen -source="$$dependency" -destination="$$package_dir/mock_dependencies_test.go" -package="$$package_name"; \
	done

fmt: ## Format all Go source files
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

vet: ## Run go vet
	go vet ./...

tidy: ## Synchronize module dependencies
	go mod tidy

migrate-up: ## Apply PostgreSQL migrations using current DATABASE_URL
	go run ./cmd/migrate -direction up

migrate-down: ## Roll back the latest PostgreSQL migration
	go run ./cmd/migrate -direction down

proto: ## Generate Go and gRPC code from all protobuf contracts
	protoc -I ../protos \
		--go_out=paths=source_relative:./gen \
		--go-grpc_out=paths=source_relative:./gen \
		../protos/reader/v1/*.proto

proto-lint: ## Validate protobuf contracts when buf is installed
	buf lint ../protos
