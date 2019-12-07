GO ?= go

server:
	$(GO) run cmd/server/main.go

processor:
	$(GO) run cmd/processor/main.go