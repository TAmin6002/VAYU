.PHONY: up down reset logs sqlc-generate

## Build and start everything (database + app) in Docker.
up:
	docker compose up --build

## Same, but in the background.
up-d:
	docker compose up --build -d

## Stop the containers (data is kept).
down:
	docker compose down

## Stop the containers AND delete the database data (schema + default admin are recreated on next `up`).
reset:
	docker compose down -v

logs:
	docker compose logs -f app

## Requires the sqlc CLI: https://docs.sqlc.dev/en/latest/overview/install.html
sqlc-generate:
	sqlc generate
