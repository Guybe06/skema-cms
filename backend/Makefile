BINARY     = skema-api
CMD        = ./cmd/server
DOCS_DIR   = docs
SWAG       = $(shell go env GOPATH)/bin/swag

.PHONY: build run dev docs stop

## build : compile le binaire
build:
	go build -o $(BINARY).exe $(CMD)/

## docs : génère la documentation Swagger
docs:
	$(SWAG) init -g $(CMD)/main.go -o $(DOCS_DIR)

## run : compile, génère les docs et démarre le serveur
run: docs build
	./$(BINARY).exe

## dev : démarre le serveur sans recompiler (utilise le binaire existant)
dev:
	./$(BINARY).exe

## stop : arrête le serveur
stop:
	-taskkill /F /IM $(BINARY).exe 2>nul || true
