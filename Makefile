GO := go
PROJECT_NAME := slop

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

dev-serve:
	@air cmd/serve/serve.go | humanlog

dev-serve-test:
	@air cmd/serve_test/serve.go | humanlog

test:
	@$(GO) test ./...

serve: build
	@./bin/serve

hserve: build
	@./bin/serve | humanlog

serve-test: build
	@./bin/serve_test

hserve-test: build
	@./bin/serve_test | humanlog
