include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE=docker compose

.PHONY: env-up env-down env-cleanup
todoapp-run:
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	go mod tidy && \
	go run cmd/todoapp/main.go
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

migrate-up:
	make migrate-action action=up

migrate-down:
	make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует action"; \
		exit 1; \
	fi
	${COMPOSE} run --rm todoapp-postgres-migrations \
	 	-path /migrations \
	 	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432:/${POSTGRES_DB}?sslmode=disable \
	 	"$(action)"

env-port-forward:
	@${COMPOSE} up -d port-forwarder

env-port-forward-down:
	@${COMPOSE} down port-forwarder
