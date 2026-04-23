# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Garden of the Goddess of Flowers** — social web game + Telegram Mini App where users grow virtual flowers.
Players water flowers daily, earn FD (Flowering Days) currency, collect seeds, and visit each other's gardens.

Full game logic: `docs/game-logic.md`  
Development plan: `docs/planning.md`

## Project Structure

```
server/       # Go REST API
client/       # SvelteKit web app + Telegram Mini App
client-3d/    # Zig + raylib + Sciter 3D client
docs/         # game logic, planning, agent plans
```

## Stack

| Service   | Stack                                          |
|-----------|------------------------------------------------|
| server    | Go, chi, pgx, goose, golang-jwt, air (hot-reload) |
| client    | SvelteKit + TypeScript + pnpm, ui-kit          |
| client-3d | Zig 0.13.0 + raylib 5.x + Sciter               |
| database  | PostgreSQL                                     |

## Versioning Rules

**Always pin exact versions everywhere — no `latest`, no `^`, no `~`.**

- Dockerfile base images: `golang:1.24.2-alpine`, `postgres:17.4`, etc.
- docker-compose service images: same — exact tag
- Go modules: exact versions in `go.mod` (default behaviour)
- pnpm packages: exact versions — `.npmrc` contains `save-exact=true`
- Always use the latest stable version available at the time of writing, then pin it

## Common Commands

```bash
# Start everything (dev)
make dev

# Go server (hot-reload via air)
cd server && air

# SvelteKit client
cd client && pnpm dev

# DB migrations
cd server && goose up
```

## Go Server Architecture (`server/`)

```
cmd/server/       # main.go — entry point
internal/
  handler/        # HTTP handlers (chi router)
  service/        # business logic
  repo/           # database queries (pgx, raw SQL)
  model/          # domain types
migrations/       # SQL migration files (goose)
```

- Auth: `POST /auth/register`, `POST /auth/login` → JWT (30d)
- All other routes require `Authorization: Bearer <token>`
- Daily tick: cron at midnight — dries unwatered flowers, awards FD, gives seeds, writes notifications

## SvelteKit Client Architecture (`client/`)

```
src/
  routes/         # file-based routing
  lib/
    api.ts        # typed fetch wrapper, JWT management
    components/   # UI components (from ui-kit + custom)
```

- UI-kit: `iwwwanow.github.io/ui-kit/` — use as primary component source
- Missing components: add new ones consistent with existing kit (tokens, style, behaviour)
- TMA: `@twa-dev/sdk` — optional Telegram integration
- API proxy: `/api/*` → `localhost:8080` in dev (vite config)

## 3D Client (`client-3d/`)

- Zig 0.13.0 + raylib 5.x + Sciter (C API via `@cImport`)
- Desktop native + WASM (Sciter replaced by DOM overlay in WASM build)
- Phase 0 plan: `docs/agents/3d-client-planning.md`

## Database Schema

Tables: `users`, `flowers`, `user_flowers`, `waterings`, `seeds`, `notifications`  
Details: `docs/game-logic.md`

## Game Rules (summary)

- FD formula: `FD(day) = day` (linear, infinite growth)
- Flower dies if not watered by anyone in a day
- One watering per flower per user per day; unlimited flowers can be watered per day
- Seeds: +1 every 7 days per flower
- Max 64 active flowers per user
