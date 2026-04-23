.PHONY: dev dev-db server client migrate-up migrate-down

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

# DB migrations
migrate-up:
	cd server && goose -dir migrations postgres "$(DATABASE_URL)" up

migrate-down:
	cd server && goose -dir migrations postgres "$(DATABASE_URL)" down
