ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)

.PHONY: all build clean test serve test-serve


CMD_DIR := cmd
CMDS := $(wildcard $(CMD_DIR)/*)
BINARIES := $(patsubst $(CMD_DIR)/%,bin/%,$(CMDS))

all: @build

$(BINARIES): bin/%: $(CMD_DIR)/%
	@$(GO) build -o $@ ./$<

build: clean
	@$(MAKE) $(BINARIES)

clean:
	@$(GO) clean
	@rm -f bin/*

db-reset:
	@psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@migrate -path $(MIGRATIONS_DIR) -database postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up 2

dump:
	@echo "Environment Variables:"
	@echo "----------------------"
	@echo "DB_HOST: $(DB_HOST)"
	@echo "DB_USER: $(DB_USER)"
	@echo "DB_PORT: $(DB_PORT)"
	@echo "DB_NAME: $(DB_NAME)"
	@echo ""
	@echo "Test Variables:"
	@echo "----------------"
	@echo "TEST_FIRST_NAME: $(TEST_FIRST_NAME)"
	@echo "TEST_LAST_NAME: $(TEST_LAST_NAME)"
	@echo "TEST_EMAIL: $(TEST_EMAIL)"
	@echo "TEST_PASSWORD: $(TEST_PASSWORD)"
	@echo ""
	@echo "Current Branch:"
	@echo "---------------"
	@echo "CURRENT_BRANCH: $(CURRENT_BRANCH)"
	@echo ""
	@echo "Binary Targets:"
	@echo "---------------"
	@echo "$(BINARIES)"

dev-serve:
	@air cmd/serve/serve.go | humanlog

dev-serve-test:
	@air cmd/serve_test/serve.go | humanlog

test:
	@$(GO) test ./...

test-register:
	@curl -X POST http://localhost:3000/register \
	     -H "Content-Type: application/json" \
	     -d '{"first_name": "$(TEST_FIRST_NAME)", "last_name": "$(TEST_LAST_NAME)", "email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}'

test-login:
	@curl -X POST http://localhost:3000/login \
	     -H "Content-Type: application/json" \
	     -d '{"email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}'

serve: build
	@./bin/serve

hserve: build
	@./bin/serve | humanlog

serve-test: build
	@./bin/serve_test

hserve-test: build
	@./bin/serve_test | humanlog
