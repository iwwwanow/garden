# Planning — Garden of the Goddess of Flowers

> Логика игры: `docs/game-logic.md`  
> 3D клиент tech validation: `docs/agents/3d-client-planning.md`

---

## Фаза 0 — Инфраструктура ✅

**Цель:** все три сервиса запускаются локально одной командой.

- [x] Удалить старый `server/` (Node/Express/Drizzle)
- [x] Удалить старый `client/` (React/Svelte — всё)
- [x] Создать `server/` — Go 1.26.2, chi v5.2.5, pgx v5.9.2, goose v3.27.0, golang-jwt v5.3.1
- [x] Создать `client/` — SvelteKit 2.57.1, Svelte 5.55.4, TypeScript 6.0.3, adapter-node
- [x] Создать `client-3d/` — Zig 0.16.0 stub, компилируется
- [x] `docker-compose.yml` — postgres:17.9, golang:1.26.2, node:22.16.0-alpine
- [x] Hot-reload: Go → `air v1.62.0`; SvelteKit → Vite HMR
- [x] `.env.example` — DATABASE_URL, JWT_SECRET, PORT, POSTGRES_*
- [x] `Makefile` — `make dev`, `make server`, `make client`, `make migrate-up/down`, `make db-reset`, `make db-clean`
- [x] Миграция `00001_init.sql` — все таблицы схемы

---

## Фаза 1 — Go-сервер ✅

**Цель:** полностью рабочий REST API, покрытый тестами, проверенный curl-ами.

### 1.1 Основа

- [x] Структура: `cmd/server/`, `internal/handler/`, `internal/service/`, `internal/repo/`, `internal/model/`, `config/`, `migrations/`
- [x] HTTP-роутер: chi v5.2.5
- [x] PostgreSQL: pgx v5.9.2 (без ORM, сырые запросы)
- [x] Миграции: goose v3.27.0, автозапуск при старте сервера
- [x] Конфиг: `config/config.go`, читает из `.env` / env vars
- [x] DB-пул: инициализация `pgxpool` в `main.go`, передача в handlers
- [x] JWT middleware: `internal/middleware/auth.go`

### 1.2 Схема БД ✅

- `users`, `flowers`, `user_flowers`, `waterings`, `seeds`, `notifications`
- `00002_seed_flowers.sql` — 3 шаблона цветков сезона 1
- `00003_watering_unique_per_day.sql` — один полив на цветок в день (любым пользователем)

### 1.3 Эндпоинты

| Метод | Путь                     | Реализован | Тесты |
|-------|--------------------------|-----------|-------|
| POST  | /api/auth/register       | [x]       | [x]   |
| POST  | /api/auth/login          | [x]       | [x]   |
| GET   | /api/me                  | [x]       | [ ]   |
| POST  | /api/flowers/:id/water   | [x]       | [x]   |
| GET   | /api/flowers/user/:id    | [x]       | [ ]   |
| POST  | /api/flowers/plant       | [x]       | [x]   |
| GET   | /api/leaderboard         | [x]       | [ ]   |
| GET   | /api/seeds               | [x]       | [ ]   |
| POST  | /api/seeds/share         | [x]       | [ ]   |
| GET   | /api/notifications       | [x]       | [ ]   |

### 1.4 Dev-эндпоинты (только при `APP_ENV=development`)

| Метод | Путь             | Описание                                      |
|-------|------------------|-----------------------------------------------|
| POST  | /api/dev/tick    | Запустить ежедневный тик вручную (+24h)        |
| GET   | /api/dev/users   | Список всех пользователей для смены профиля   |
| POST  | /api/dev/seeds   | Выдать семена текущему пользователю           |
| POST  | /api/dev/reset   | Очистить все пользовательские данные          |

### 1.5 Бизнес-логика

- [x] Репозиторный слой (`internal/repo/`) — CRUD для каждой таблицы
- [x] Сервисный слой (`internal/service/`) — бизнес-правила, изолирован от БД
- [x] Ежедневный тик: cron-планировщик, проверка поливов, начисление FD, семена, уведомления
- [x] Лимит 64 активных цветков на пользователя
- [x] Один полив на цветок в день (любым пользователем) — UNIQUE (user_flower_id, watered_date)

### 1.6 Тестирование

- [x] Unit-тесты сервисного слоя (моки репозитория через интерфейсы)
- [x] Integration-тесты с реальной БД (`testcontainers-go`)
- [x] `docs/curl-tests.md` — полный сценарий с dev-эндпоинтами

---

## Доработки API (приоритет перед Фазой 2)

Без этого клиент не сможет нормально работать:

| Приоритет | Задача | Почему нужно |
|-----------|--------|--------------|
| 🔴 Высокий | `GET /api/flowers` — список шаблонов цветков | Клиент должен показывать, что можно посадить |
| 🔴 Высокий | `PATCH /api/notifications/read` — отметить прочитанными | Иначе бейдж уведомлений никогда не пропадёт |
| 🟡 Средний | `GET /api/flowers/:id` — цветок по ID | Нужен для страницы детали цветка |
| 🟡 Средний | `PUT /api/me` — обновить first_name | Экран Profile |
| 🟢 Низкий | Пагинация leaderboard (`?limit=&offset=`) | Когда пользователей станет много |
| 🟢 Низкий | `GET /api/users/:id` — публичный профиль | Для страницы Garden чужого пользователя |

---

## Фаза 2 — Web-клиент

**Цель:** полностью рабочий web-app / TMA.

**Стек:** SvelteKit 2.57.1 + TypeScript 6.0.3 + pnpm + ui-kit (`iwwwanow.github.io/ui-kit/`)

**UI-kit:** компоненты из кита — основа. Недостающие добавлять консистентно (токены, стиль, поведение).

### Dev-панель (только в dev-режиме)

Плавающее окно поверх интерфейса для тестирования:
- [ ] Перемещается по экрану (drag), сворачивается
- [ ] `+24h` — вызывает `POST /api/dev/tick`
- [ ] `Change user` — переключение через `GET /api/dev/users` без выхода/входа
- [ ] `Give seeds` — вызывает `POST /api/dev/seeds`
- [ ] `Reset DB` — вызывает `POST /api/dev/reset`
- [ ] Рендерится только при `dev: true` в SvelteKit (не попадает в prod-сборку)

### Задачи

- [x] Scaffold SvelteKit, adapter-node, vite proxy `/api → localhost:8080`
- [ ] Подключить ui-kit
- [ ] `src/lib/api.ts` — типизированный fetch-wrapper с JWT
- [ ] TMA: `@twa-dev/sdk`, опциональная привязка Telegram
- [ ] Dev-панель (`src/lib/components/DevPanel.svelte`)

### Экраны

- [ ] Auth — регистрация / вход
- [ ] Home — свои цветки, полив, день/FD
- [ ] Garden — профиль друга, полить чужой цветок
- [ ] Leaderboard — топ по FD
- [ ] Seeds — инвентарь, поделиться, посадить
- [ ] Herbarium — пересохшие цветки
- [ ] Notifications — лента событий
- [ ] Profile — username, FD-баланс

---

## Фаза 3 — 3D-клиент (singleplayer)

**Цель:** нативный клиент для своей фермы.

**Стек:** Zig 0.16.0 + raylib 5.x + Sciter

- [x] Scaffold `client-3d/` — Zig проект, компилируется
- [ ] Tech validation: куб + Sciter UI-оверлей + WASM — по `docs/agents/3d-client-planning.md`
- [ ] Singleplayer: HTTP-клиент → Go API, цветок, полив, семена, полка
- [ ] Multiplayer: визиты, полив чужих цветков
- [ ] Платформа для кастомных игр (отдельный разговор)

---

## Текущий приоритет

```
✅ Фаза 0 → ✅ Фаза 1 → 🔄 Доработки API → Фаза 2 (Web-клиент) → Фаза 3 (3D)
```
