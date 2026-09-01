include .env
export

.PHONY: help up down run generate migrate migrate-down tidy build

help:
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## up: поднять PostgreSQL
up:
	docker compose up -d
	@echo "Ждём готовности PostgreSQL..."
	@until docker compose exec postgres pg_isready -U $(POSTGRES_USER) -d $(POSTGRES_DB) > /dev/null 2>&1; do sleep 1; done
	@echo "PostgreSQL готов!"

## down: остановить
down:
	docker compose down

## down-v: остановить и удалить данные
down-v:
	docker compose down -v

## run: запустить сервис
run:
	go run ./cmd/server/main.go

## generate: кодогенерация из OpenAPI
generate:
	go generate ./..

## migrate: применить миграции
migrate:
	goose -dir migrations postgres \
	"host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" up

## migrate-down: откатить последнюю миграцию
migrate-down:
	goose -dir migrations postgres \
	"host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" down

## migrate-status: статус миграций
migrate-status:
	goose -dir migrations postgres \
	"host=$(POSTGRES_HOST) port=$(POSTGRES_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=disable" status

## tidy: обновить зависимости
tidy:
	go mod tidy

## build: собрать бинарник
build:
	go build -o bin/server ./cmd/server/main.go