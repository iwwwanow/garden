# Planning — Garden of the Goddess of Flowers

> Логика игры: `docs/game-logic.md`  
> 3D клиент фаза 0: `docs/agents/3d-client-planning.md`

---

## Фаза 0 — Инфраструктура

**Цель:** все три сервиса запускаются локально одной командой.

- [ ] Удалить старый `server/` (Node/Express/Drizzle)
- [ ] Удалить старый `client/` (React/Svelte — всё)
- [ ] Создать `server/` — Go-модуль (`go mod init`)
- [ ] Создать `client/` — Svelte + Vite + TypeScript (`npm create vite`)
- [ ] Создать `client-3d/` — Zig-проект со структурой из `docs/agents/3d-client-planning.md`
- [ ] `docker-compose.yml` — PostgreSQL, Go-сервер, Svelte-клиент (с hot-reload в dev)
- [ ] `.env.example` — DATABASE_URL, JWT_SECRET, PORT
- [ ] Один `make dev` поднимает всё

---

## Фаза 1 — Go-сервер

**Цель:** полностью рабочий REST API, покрытый тестами, проверенный curl-ами.

### 1.1 Основа

- [ ] Структура проекта: `cmd/`, `internal/handler/`, `internal/service/`, `internal/repo/`, `internal/model/`
- [ ] HTTP-роутер: `chi` или `net/http` стандартный
- [ ] PostgreSQL: `pgx` (без ORM, сырые запросы)
- [ ] Миграции: `goose` или `migrate` — SQL-файлы в `migrations/`
- [ ] JWT middleware (`golang-jwt/jwt`)
- [ ] Конфиг из `.env` (`joho/godotenv`)

### 1.2 Схема БД

Таблицы согласно `docs/game-logic.md`:
- `users` — id, username, password_hash, first_name, telegram_id?, fd_balance, created_at
- `flowers` — id, season, image_path, created_at
- `user_flowers` — id, user_id, flower_id, day, last_watered_at, is_dried, created_at
- `waterings` — id, user_flower_id, watered_by_user_id, watered_date
- `seeds` — id, user_id, flower_id, quantity
- `notifications` — id, user_id, type, payload (JSONB), read, created_at

### 1.3 Эндпоинты

| Метод | Путь                  | Тесты |
|-------|-----------------------|-------|
| POST  | /auth/register        | [ ]   |
| POST  | /auth/login           | [ ]   |
| GET   | /me                   | [ ]   |
| POST  | /flowers/:id/water    | [ ]   |
| GET   | /flowers/user/:userId | [ ]   |
| POST  | /flowers/plant        | [ ]   |
| GET   | /leaderboard          | [ ]   |
| GET   | /seeds                | [ ]   |
| POST  | /seeds/share          | [ ]   |
| GET   | /notifications        | [ ]   |

### 1.4 Бизнес-логика

- [ ] Ежедневный тик (cron или встроенный scheduler): проверка поливов, начисление FD, выдача семян, запись уведомлений
- [ ] Лимит: не более 64 активных цветков на пользователя
- [ ] Один полив на цветок в день от одного пользователя (constraint в БД + проверка в сервисе)

### 1.5 Тестирование

- [ ] Unit-тесты: сервисный слой (бизнес-логика изолирована от БД)
- [ ] Integration-тесты: тестовая БД (docker или `testcontainers-go`)
- [ ] curl-сценарий в `docs/curl-tests.md`: регистрация → логин → посадить цветок → полить → лидерборд

---

## Фаза 2 — Web-клиент

**Цель:** полностью рабочий web-app / TMA.

**Стек:** Svelte 5 + Vite + TypeScript + ui-kit (`iwwwanow.github.io/ui-kit/`)

### Экраны (уточняются при наличии набросков)

- [ ] Auth — регистрация / вход
- [ ] Home — свои цветки, полив, день/FD
- [ ] Garden — посетить профиль друга, полить чужой цветок
- [ ] Leaderboard — топ по FD
- [ ] Seeds — инвентарь семян, поделиться, посадить
- [ ] Herbarium — пересохшие цветки
- [ ] Notifications — лента событий
- [ ] Profile — аватар, username, FD-баланс

### Задачи

- [ ] Scaffold: Svelte 5 + TypeScript + Vite
- [ ] Подключить ui-kit
- [ ] `api.ts` — типизированный fetch-wrapper с JWT
- [ ] Роутинг (`svelte-routing` или file-based)
- [ ] TMA: инициализация `@twa-dev/sdk`, опциональная привязка Telegram
- [ ] Деплой: Docker, статика через nginx или `vite preview`

---

## Фаза 3 — 3D-клиент (singleplayer)

**Цель:** рабочий нативный клиент для своей фермы.

**Стек:** Zig 0.13.0 + raylib 5.x + Sciter

Детальный план фазы 0: `docs/agents/3d-client-planning.md`

### Этапы

- [ ] Фаза 0 (tech POC): куб + Sciter UI-оверлей + WASM-сборка — по `3d-client-planning.md`
- [ ] Фаза 1 (singleplayer):
  - [ ] Подключение к Go API (HTTP-клиент в Zig)
  - [ ] Отрисовка своего цветка (3D-модель по дню роста)
  - [ ] Полив из интерфейса
  - [ ] FD-баланс и семена
  - [ ] Своя территория / «полка» с цветками
- [ ] Фаза 2 (multiplayer): визиты к друзьям, полив чужих цветков
- [ ] Фаза N: платформа для кастомных игр (отдельный разговор)

---

## Текущий приоритет

```
Фаза 0 (инфра) → Фаза 1 (Go-сервер) → Фаза 2 (Web-клиент) → Фаза 3 (3D)
```
