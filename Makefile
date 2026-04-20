include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE=docker compose

.PHONY: env-up env-down env-cleanup

env-up:
	$(COMPOSE) up -d todoapp-postgres


env-down:
	$(COMPOSE) down

env-cleanup:
	@read -p "очистить volumes? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(COMPOSE) down && \
		rm -rf out/pgdata && \
		echo "Удалено"; \
	else \
		echo "Очистка отменена"; \
	fi

seq ?= init

migration-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует seq"; \
		exit 1; \
	fi
	$(COMPOSE) run --rm todoapp-postgres-migrations \
		create \
		-ext sql \
		-dir /migrations \
		-seq \
		"$(seq)"
