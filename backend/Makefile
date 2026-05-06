BINARY     = skema-api
CMD        = ./cmd/server
MIGRATE    = ./cmd/migrate
DOCS_DIR   = docs
SWAG       = $(shell go env GOPATH)/bin/swag

.PHONY: build run dev docs stop migrate-up migrate-down migrate-step-up migrate-step-down migrate-force

## build : compile le binaire
build:
	go build -o $(BINARY).exe $(CMD)/

## docs : génère la documentation Swagger
docs:
	$(SWAG) init -g $(CMD)/main.go -d $(CMD) --parseDependency -o $(DOCS_DIR)

## run : compile, génère les docs et démarre le serveur
run: docs build
	./$(BINARY).exe

## dev : démarre le serveur sans recompiler (utilise le binaire existant)
dev:
	./$(BINARY).exe

## stop : arrête le serveur
stop:
	-taskkill /F /IM $(BINARY).exe 2>nul || true

## migrate-up : applique toutes les migrations
migrate-up:
	go run $(MIGRATE)/main.go up

## migrate-down : rollback toutes les migrations
migrate-down:
	go run $(MIGRATE)/main.go down

## migrate-step-up : applique N migrations (ex: make migrate-step-up N=1)
migrate-step-up:
	go run $(MIGRATE)/main.go -steps $(N) up

## migrate-step-down : rollback N migrations (ex: make migrate-step-down N=1)
migrate-step-down:
	go run $(MIGRATE)/main.go -steps $(N) down

## migrate-force : force version sans exécuter la migration (ex: make migrate-force N=6)
migrate-force:
	go run $(MIGRATE)/main.go -steps $(N) force
