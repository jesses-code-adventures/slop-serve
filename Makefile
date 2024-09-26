ifneq (,$(wildcard ./.env))
    include .env
    export $(shell sed 's/=.*//' .env)
endif

CURRENT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
CMD_DIR := cmd
CMDS := $(wildcard $(CMD_DIR)/*)
DB_CONNECTION_STRING := postgres://$(DB_USER)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable
BINARIES := $(patsubst $(CMD_DIR)/%,bin/%,$(CMDS))

.PHONY: all build clean dump help db-reset dev-serve dev-serve-test gen test-register test-login hserve hserve-test reset test serve test-serve


help:
	@echo "Available Commands:"
	@echo ""
	@echo "General:"
	@echo "  all                    - Build all binaries."
	@echo "  build                  - Clean and build all binaries."
	@echo "  clean                  - Clean up built binaries."
	@echo "  dump                   - Dump environment variables and current branch information."
	@echo ""
	@echo "Database:"
	@echo "  db-reset               - Drop and recreate the database, then run migrations."
	@echo ""
	@echo "Development:"
	@echo "  dev                    - Serve the application in development mode with live reload."
	@echo "  dev-test               - Serve the test application in development mode with live reload."
	@echo "  gen                    - Generate code from SQL schemas."
	@echo ""
	@echo "Testing:"
	@echo "  test                   - Run all tests."
	@echo "  test-register          - Test the register endpoint."
	@echo "  test-login             - Test the login endpoint."
	@echo "  test-login-fail	 - Test the login endpoint."
	@echo ""
	@echo "Serving:"
	@echo "  serve                  - Clean, build, and serve the application."
	@echo "  hserve                 - Clean, build, and serve the application with human-readable logs."
	@echo "  serve-test             - Clean, build, and serve the test application."
	@echo "  hserve-test            - Clean, build, and serve the test application with human-readable logs."
	@echo ""
	@echo "Utilities:"
	@echo "  reset                  - Clean, reset the database, generate code, and build binaries."


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
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_CONNECTION_STRING) up

dump:
	@echo "Environment Variables:"
	@echo "----------------------"
	@echo "DB_HOST: $(DB_HOST)"
	@echo "DB_USER: $(DB_USER)"
	@echo "DB_PORT: $(DB_PORT)"
	@echo "DB_NAME: $(DB_NAME)"
	@echo "DB_CONNECTION_STRING: $(DB_CONNECTION_STRING)"
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

dev:
	@air cmd/serve/serve.go | humanlog

dev-reset: reset dev

dev-test:
	@air cmd/serve_test/serve.go | humanlog

gen:
	@sqlc generate

test:
	@$(GO) test ./...

reset: db-reset clean gen

serve: clean build
	@./bin/serve

hserve: clean build
	@./bin/serve | humanlog

serve-test: clean build
	@./bin/serve_test

hserve-test: clean build
	@./bin/serve_test | humanlog

test-register:
	@curl -X POST http://localhost:3000/register \
	     -H "Content-Type: application/json" \
	     -d '{"first_name": "$(TEST_FIRST_NAME)", "last_name": "$(TEST_LAST_NAME)", "email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}' \
	     -s | jq -c

test-login:
	@curl -X POST http://localhost:3000/login \
	     -H "Content-Type: application/json" \
	     -d '{"email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}' \
	     -s | jq -c

test-login-fail:
	@curl -X POST http://localhost:3000/login \
	     -H "Content-Type: application/json" \
	     -d '{"email": "$(TEST_EMAIL)","password": "this-is-wrong"}' \
	     -s | jq -c

test-image-generate:
	@token=$(shell curl -X POST http://localhost:3000/login \
	    -H "Content-Type: application/json" \
	    -d '{"email": "$(TEST_EMAIL)", "password": "$(TEST_PASSWORD)"}' -s | jq -r .token); \
	curl -X POST http://localhost:3000/image \
	    --cookie "Authorization=$$token" \
	    -H "x-slop-user-id: $(TEST_USER_ID)" \
	    -F "image=@./test.jpeg" \
	    -F "character=dave" -s

