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

# HELP - will output the help for each task in the Makefile
# In sorted order.
# The width of the first column can be determined by the `width` value passed to awk
#
# thanks to https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html for the initial version.
#
help: ## This help.
	@grep -E -h "^[a-zA-Z_-]+:.*?## " $(MAKEFILE_LIST) \
	  | sort \
	  | awk -v width=36 'BEGIN {FS = ":.*?## "} {printf "\033[36m%-*s\033[0m %s\n", width, $$1, $$2}'

all: @build ## Build all binaries.

$(BINARIES): bin/%: $(CMD_DIR)/%
	@$(GO) build -o $@ ./$<

build: clean ## Clean and build all binaries.
	@$(MAKE) $(BINARIES)

clean: ## Clean up built binaries.
	@$(GO) clean
	@rm -f bin/*


db-reset: ## Drop and recreate the database, then run migrations.
	@psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "DROP DATABASE IF EXISTS $(DB_NAME);"
	@psql -h $(DB_HOST) -U $(DB_USER) -p $(DB_PORT) -d postgres -c "CREATE DATABASE $(DB_NAME);"
	@migrate -path $(MIGRATIONS_DIR) -database $(DB_CONNECTION_STRING) up

dump: ## Dump environment variables and current branch information.
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

dev: ## Serve the application in development mode with live reload.
	@air cmd/serve/serve.go | humanlog

dev-reset: reset dev ## Reset the db, clean binaries, run codegen and run the command in .air.toml with live reload.

dev-test: ## Serve the test application in development mode with live reload.
	@air cmd/serve_test/serve.go | humanlog

gen: ## Generate code from SQL schemas.
	@sqlc generate

test: ## Run all tests.
	@$(GO) test ./...

reset: db-reset clean gen ## Clean, reset the database, generate code, and build binaries.

serve: clean build ## Clean, build, and serve the application.
	@./bin/serve

hserve: clean build ## Clean, build, and serve the application with human-readable logs.
	@./bin/serve | humanlog

serve-test: clean build ## Clean, build, and serve the test application.
	@./bin/serve_test

hserve-test: clean build ## Clean, build, and serve the test application with human##readable logs.
	@./bin/serve_test | humanlog

test-register: ## Test the POST /register endpoint.
	@curl -X POST http://localhost:3000/register \
	     -H "Content-Type: application/json" \
	     -d '{"first_name": "$(TEST_FIRST_NAME)", "last_name": "$(TEST_LAST_NAME)", "email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}' \
	     -s | jq -c

test-login: ## Test the POST /login endpoint.
	@curl -X POST http://localhost:3000/login \
	     -H "Content-Type: application/json" \
	     -d '{"email": "$(TEST_EMAIL)","password": "$(TEST_PASSWORD)"}' \
	     -s | jq -c

test-login-fail: ## Test POST /login that should fail.
	@curl -X POST http://localhost:3000/login \
	     -H "Content-Type: application/json" \
	     -d '{"email": "$(TEST_EMAIL)","password": "this-is-wrong"}' \
	     -s | jq -c

test-image-post: ## Test POST /image
	@token=$(shell curl -X POST http://localhost:3000/login \
	    -H "Content-Type: application/json" \
	    -d '{"email": "$(TEST_EMAIL)", "password": "$(TEST_PASSWORD)"}' -s | jq -r .token); \
	curl -X POST http://localhost:3000/image \
	    --cookie "Authorization=$$token" \
	    -H "x-slop-user-id: $(TEST_USER_ID)" \
	    -F "image=@./test.jpeg" \
	    -F "url=https://google.com/" -s | jq -c

