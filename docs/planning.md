# Planning — Garden of the Goddess of Flowers

> Логика игры: `docs/game-logic.md`  
> UI-спецификация (источник правды для интерфейса): `docs/ui-spec.md`  
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
- [x] `.env.example` — DATABASE*URL, JWT_SECRET, PORT, POSTGRES*\*
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

| Метод | Путь                    | Реализован | Тесты |
| ----- | ----------------------- | ---------- | ----- |
| POST  | /api/auth/register      | [x]        | [x]   |
| POST  | /api/auth/login         | [x]        | [x]   |
| GET   | /api/me                 | [x]        | [ ]   |
| PUT   | /api/me                 | [x]        | [ ]   |
| GET   | /api/flowers            | [x]        | [x]   |
| GET   | /api/flowers/:id        | [x]        | [x]   |
| POST  | /api/flowers/plant      | [x]        | [x]   |
| POST  | /api/flowers/:id/water  | [x]        | [x]   |
| GET   | /api/flowers/user/:id   | [x]        | [ ]   |
| GET   | /api/leaderboard        | [x]        | [ ]   |
| GET   | /api/seeds              | [x]        | [x]   |
| POST  | /api/seeds/share        | [x]        | [x]   |
| GET   | /api/notifications      | [x]        | [x]   |
| PATCH | /api/notifications/read | [x]        | [x]   |
| GET   | /api/users/:id          | [x]        | [ ]   |

### 1.4 Dev-эндпоинты (только при `APP_ENV=development`)

| Метод | Путь           | Описание                                    |
| ----- | -------------- | ------------------------------------------- |
| POST  | /api/dev/tick  | Запустить ежедневный тик вручную (+24h)     |
| GET   | /api/dev/users | Список всех пользователей для смены профиля |
| POST  | /api/dev/seeds | Выдать семена текущему пользователю         |
| POST  | /api/dev/reset | Очистить все пользовательские данные        |

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

## Доработки API ✅

- [x] `GET /api/flowers` — список шаблонов цветков
- [x] `PATCH /api/notifications/read` — отметить все прочитанными
- [x] `GET /api/flowers/:id` — цветок по ID
- [x] `PUT /api/me` — обновить first_name
- [x] Пагинация leaderboard (`?limit=&offset=`)
- [x] `GET /api/users/:id` — публичный профиль

---

## Фаза 2 — Web-клиент 🔄

**Цель:** полностью рабочий web-app + Telegram Mini App.

**Стек:** SvelteKit 2.57.1 + TypeScript 6.0.3 + pnpm + ui-kit (`iwwwanow.github.io/ui-kit/`)

**UI-спецификация:** `docs/ui-spec.md` — источник правды для всего интерфейса: дизайн-система, компоненты, страницы, поведение.

**UI-kit:** компоненты из кита — основа. Недостающие добавлять консистентно (токены, стиль, поведение).

### Архитектура SSR/CSR

**Текущее решение: все страницы CSR** (`ssr = false`).

Плановая таблица SSR не реализована — заблокирована тремя техническими причинами:

1. **JWT в localStorage** — SSR рендер происходит на Node.js, `localStorage` там недоступен. Для SSR с авторизацией нужен переход на httpOnly cookie.
2. **Go API требует `Authorization` на всех эндпоинтах** — даже «публичные» leaderboard и users/:id закрыты токеном. SSR-запрос с Node.js не может его передать без cookie.
3. **Vite proxy не работает на сервере** — `fetch('/api/...')` во время SSR уходит в SvelteKit, а не в Go. Нужна переменная `PUBLIC_API_BASE=http://localhost:8080`.

**Что нужно для SSR** (отложено в бэклог):
- Перейти на cookie-based JWT (httpOnly, `Set-Cookie` при логине)
- Добавить `PUBLIC_API_BASE` в env, использовать в SSR load-функциях
- Вынести публичные эндпоинты (leaderboard, users/:id, flowers/user/:id) из-за auth-middleware
- Переписать leaderboard и garden/[id] на `+page.server.ts`

### Структура `src/`

```
routes/
  +layout.svelte          # глобальный layout, подключение стилей
  +layout.ts              # глобальный load — проверка JWT, редирект
  +page.svelte            # редирект: авториз → /home, нет → /auth
  auth/
    +page.svelte          # логин / регистрация
    +page.ts
  home/
    +page.svelte          # свои цветки, полив, FD
    +page.ts              # ssr = false
  garden/
    [id]/
      +page.svelte        # публичный сад пользователя, полив
      +page.ts            # SSR load: GET /api/users/:id + /api/flowers/user/:id
  leaderboard/
    +page.svelte          # топ по FD
    +page.ts              # SSR load: GET /api/leaderboard
  seeds/
    +page.svelte          # инвентарь, поделиться, посадить
    +page.ts              # ssr = false
  herbarium/
    +page.svelte          # пересохшие цветки
    +page.ts              # ssr = false
  notifications/
    +page.svelte          # лента событий, кнопка «прочитать всё»
    +page.ts              # ssr = false
  profile/
    +page.svelte          # username, FD-баланс, edit first_name
    +page.ts              # ssr = false
lib/
  api.ts                  # типизированный fetch-wrapper, JWT
  stores/
    auth.ts               # user store (writable)
  components/
    DevPanel.svelte       # dev-only: +24h, change user, give seeds, reset
```

### Задачи

- [x] Scaffold SvelteKit, adapter-node, vite proxy `/api → localhost:8080`
- [x] `src/lib/api.ts` — типизированный fetch-wrapper с JWT (localStorage)
- [x] `src/lib/stores/auth.ts` — user store
- [x] Глобальный layout: защита авторизованных маршрутов, редирект
- [x] `src/app.css` — дизайн-токены и сброс стилей (по ui-spec.md)
- [ ] Подключить ui-kit
- [ ] TMA: `@twa-dev/sdk`, опциональная инициализация
- [x] Dev-панель (`src/lib/components/DevPanel.svelte`)

### Компоненты

- [x] `ActionButton.svelte` — variant: primary / confirm / muted
- [x] `FlowerCard.svelte` — type: flower / seed / herbarium
- [x] `BottomNav.svelte` — 3 таба, фиксированный
- [x] `DevPanel.svelte` — drag, tick, users, seeds, reset

### Экраны

- [x] Auth — регистрация / вход
- [x] Home — hero (цветок с макс. day, FD-баланс), preview семян и гербария
- [x] Garden — публичный сад: hero + сетка, полить чужой цветок
- [x] Leaderboard — топ по FD (Дней подряд / Семена — табы-заглушки)
- [x] Seeds — инвентарь сетка, empty state
- [x] Seed detail — посадить, поделиться (share form)
- [x] Flower detail — полить, запросить полив, статистика
- [x] Herbarium — витрина + сетка
- [x] Herbarium detail — скрыть/выставить (toggle, visibility API — бэклог)
- [x] Notifications — лента, отметить прочитанными
- [x] Profile — просмотр, редактировать имя, выход

### Dev-панель (только в dev-режиме)

- [x] Перемещается по экрану (drag), сворачивается
- [x] `+24h` — `POST /api/dev/tick`
- [x] `Change user` — переключение через `GET /api/dev/users`
- [x] `Give seeds` — `POST /api/dev/seeds`
- [x] `Reset DB` — `POST /api/dev/reset`

### Архитектурные решения

- Все страницы CSR (`ssr = false`) — SSR для leaderboard/garden/[id] в бэклоге (требует настройки PUBLIC_API_BASE и cookie-based auth)
- Маршрут-группа `(app)/` — добавляет BottomNav ко всем страницам кроме `/auth`
- Защита роутов: каждый `+page.ts` проверяет `localStorage.getItem('token')`, редирект на `/auth`
- Изображения цветков: эмодзи-плейсхолдеры до реализации static-сервинга на Go

### Интеграционные тесты клиента

**Стек:** Playwright + отдельный `docker-compose.test.yml`

- [ ] Установить Playwright в `client/` (`pnpm add -D @playwright/test`)
- [ ] `docker-compose.test.yml` — postgres + Go server (`APP_ENV=test`) + SvelteKit (built)
- [ ] `playwright.config.ts` — baseURL, запуск стека перед тестами (`webServer`)
- [ ] Тест: регистрация → логин → редирект на `/home`
- [ ] Тест: полить цветок → FD увеличился
- [ ] Тест: leaderboard рендерится, строка текущего пользователя выделена
- [ ] Тест: `/auth` редиректит на `/home` если уже авторизован
- [ ] Тест: полив на чужом профиле — user A сажает цветок, user B заходит на `/garden/:id`, нажимает полить, успех
- [ ] Тест: поделиться семенами — user A имеет семена, вводит username user B, передаёт, у user B появляются семена
- [ ] `make test-client` — команда в Makefile

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

## Бэклог (отложено)

- **SSR для leaderboard и garden/[id]:** SEO, og:tags, шаринг. Требует: cookie-based JWT, `PUBLIC_API_BASE` env, публичные эндпоинты на Go, переписать на `+page.server.ts`. Пока все страницы CSR.
- **Поливы-рейтинг ("Верда"):** таб в Leaderboard — топ пользователей по количеству поливов чужих цветков. Механика полива уже в API (`waterings`), нужен только запрос агрегации + UI-таб. Отложено до финализации основного Leaderboard.
- **Feed / Лента:** лента активности (запросы полива, события). Таб зарезервирован в nav. Отложено до backend-поддержки.
- **Herbarium visibility:** эндпоинт `PATCH /api/herbarium/:id/visibility` — скрыть/показать цветок в витрине. Нужен для страницы `/herbarium/[id]`.
- **FD-формула на бэке:** сейчас `total_fd = day*(day+1)/2` считается на клиенте. Бизнес-логика должна жить в сервисном слое Go. Вынести в `internal/service/` и отдавать готовое значение через API (поле `total_fd` в ответе `GET /api/flowers/user/:id` и `GET /api/flowers/:id`).

---

## Текущий приоритет

```
✅ Фаза 0 → ✅ Фаза 1 → ✅ Доработки API → 🔄 Фаза 2 (Web-клиент) → Фаза 3 (3D)
```
