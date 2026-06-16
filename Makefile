include .env
export

export PROJECT_ROOT=$(shell pwd)
COMPOSE=docker compose

.PHONY: dev run wait-db env-up env-down env-cleanup
dev: env-up wait-db migrate-up run

run: wait-db
	@export LOGGER_FOLDER=${PROJECT_ROOT}/out/logs && \
	export POSTGRES_HOST=localhost && \
	export POSTGRES_PORT=${POSTGRES_EXTERNAL_PORT} && \
	go mod tidy && \
	go run cmd/todoapp/main.go

.PHONY: fmt
fmt:
	gofmt -w .

env-up:
	$(COMPOSE) up -d todoapp-postgres port-forwarder

wait-db: env-up
	@until $(COMPOSE) exec -T todoapp-postgres pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; do \
		echo "Waiting for postgres..."; \
		sleep 1; \
	done


env-down:
	$(COMPOSE) down

env-cleanup:
	@read -p "очистить volumes? [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		$(COMPOSE) down todoapp-postgres port-forwarder && \
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

migrate-up: wait-db
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
	 	-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@todoapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
	 	"$(action)"

env-port-forward:
	@${COMPOSE} up -d port-forwarder

env-port-forward-down:
	@${COMPOSE} down port-forwarder
