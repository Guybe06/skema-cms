BACKEND_DIR = backend
FRONTEND_DIR = frontend

.PHONY: build build-api build-web run run-api run-web stop stop-api stop-web \
        test test-unit test-integration migrate-up migrate-down docs

# ────────────── BUILD ──────────────

## build : compile le backend et prépare le frontend
build: build-api build-web

## build-api : compile le binaire Go
build-api:
	$(MAKE) -C $(BACKEND_DIR) build

## build-web : installe les dépendances et compile le frontend Nuxt
build-web:
	cd $(FRONTEND_DIR) && npm install && npm run build

# ────────────── RUN ──────────────

## run : démarre le backend et le frontend en parallèle
run:
	$(MAKE) run-api & $(MAKE) run-web

## run-api : compile et démarre le serveur backend
run-api:
	$(MAKE) -C $(BACKEND_DIR) run

## run-web : démarre le serveur de développement Nuxt
run-web:
	cd $(FRONTEND_DIR) && npm run dev

# ────────────── STOP ──────────────

## stop : arrête les deux processus
stop: stop-api

## stop-api : arrête le serveur backend
stop-api:
	$(MAKE) -C $(BACKEND_DIR) stop

# ────────────── TESTS ──────────────

## test : lance tous les tests (unitaires + intégration)
test:
	$(MAKE) -C $(BACKEND_DIR) test

## test-unit : lance uniquement les tests unitaires
test-unit:
	$(MAKE) -C $(BACKEND_DIR) test-unit

## test-integration : lance uniquement les tests d'intégration
test-integration:
	$(MAKE) -C $(BACKEND_DIR) test-integration

# ────────────── MIGRATIONS ──────────────

## migrate-up : applique toutes les migrations
migrate-up:
	$(MAKE) -C $(BACKEND_DIR) migrate-up

## migrate-down : rollback toutes les migrations
migrate-down:
	$(MAKE) -C $(BACKEND_DIR) migrate-down

# ────────────── DOCS ──────────────

## docs : génère la documentation Swagger
docs:
	$(MAKE) -C $(BACKEND_DIR) docs
