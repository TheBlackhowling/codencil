.PHONY: build up down ps api-shell web-shell postgres migrate-up migrate-down migrate-reset

build:
	docker compose build

up:
	docker compose up -d postgres

down:
	docker compose down

ps:
	docker compose ps

api-shell:
	docker compose run --rm api sh

web-shell:
	docker compose run --rm web sh

postgres:
	docker compose up -d postgres
	docker compose ps postgres

migrate-up:
	docker compose up -d postgres
	docker compose run --rm migrate -path /migrations -database "$$DATABASE_URL" up

migrate-down:
	docker compose run --rm migrate -path /migrations -database "$$DATABASE_URL" down 1

migrate-reset: migrate-down migrate-up
