DEV = docker compose -f docker-compose.yml -f docker-compose.dev.yml
PROD = docker compose -f docker-compose.yml

.PHONY: dev dev-d dev-down dev-logs dev-migrate dev-seed dev-shell \
        prod prod-down prod-logs prod-migrate prod-seed

# DEV
dev:
	$(DEV) up --build

dev-d:
	$(DEV) up -d --build

dev-down:
	$(DEV) down

dev-logs:
	$(DEV) logs -f app

dev-migrate:
	$(DEV) run --rm migration

dev-seed:
	$(DEV) run --rm seeder

dev-shell:
	$(DEV) exec app sh

# PROD
prod:
	$(PROD) up -d --build

prod-down:
	$(PROD) down

prod-logs:
	$(PROD) logs -f app

prod-migrate:
	$(PROD) run --rm migration

prod-seed:
	$(PROD) run --rm seeder
