.PHONY: build up down ps api-shell web-shell postgres

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
