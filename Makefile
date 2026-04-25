.PHONY: dev dev-db server client migrate-up migrate-down db-reset db-clean

# Start full stack (postgres + server + client)
dev:
	docker compose up

# Start only postgres (useful for local server dev)
dev-db:
	docker compose up postgres

# Run Go server locally with hot-reload
server:
	cd server && air

# Run SvelteKit client locally
client:
	cd client && pnpm dev

# DB migrations (requires running stack: make dev)
migrate-up:
	docker compose exec server sh -c 'goose -dir migrations postgres "$$DATABASE_URL" up'

migrate-down:
	docker compose exec server sh -c 'goose -dir migrations postgres "$$DATABASE_URL" down'

# Full DB reset: drop all tables and re-run migrations (requires running stack: make dev)
db-reset:
	docker compose exec postgres psql -U $${POSTGRES_USER:-garden} -d $${POSTGRES_DB:-garden} \
	  -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	docker compose exec server sh -c 'goose -dir migrations postgres "$$DATABASE_URL" up'

# Wipe user data only (keep flower templates and schema)
db-clean:
	docker compose exec postgres psql -U $${POSTGRES_USER:-garden} -d $${POSTGRES_DB:-garden} \
	  -c "TRUNCATE notifications, waterings, seeds, user_flowers, users RESTART IDENTITY CASCADE;"
