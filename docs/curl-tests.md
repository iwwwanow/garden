# curl-tests — Garden of the Goddess of Flowers

Полный сценарий: регистрация → логин → посадить → полить → лидерборд → семена → уведомления.

Предварительно запустите сервер:
```bash
make dev   # или: cd server && air
```

---

## 0. Переменные

```bash
BASE=http://localhost:8080
TOKEN=""        # заполняется после логина
TOKEN2=""       # второй пользователь
```

---

## 1. Регистрация

```bash
# Пользователь 1 — Alice
curl -s -X POST $BASE/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123","first_name":"Alice"}' | jq .

# Пользователь 2 — Bob
curl -s -X POST $BASE/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pass456","first_name":"Bob"}' | jq .

# Дубликат → 409 Conflict
curl -s -X POST $BASE/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"x","first_name":"x"}' | jq .
```

---

## 2. Логин

```bash
TOKEN=$(curl -s -X POST $BASE/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')

TOKEN2=$(curl -s -X POST $BASE/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pass456"}' | jq -r '.token')

echo "Alice token: $TOKEN"
echo "Bob token:   $TOKEN2"

# Неверный пароль → 401
curl -s -X POST $BASE/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"wrong"}' | jq .
```

---

## 3. Профиль (/me)

```bash
curl -s $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Ожидаемый ответ:
```json
{
  "user":    { "id": 1, "username": "alice", "fd_balance": 0, ... },
  "flowers": [],
  "seeds":   []
}
```

---

## 4. Дать семена (dev)

Через SQL или `/api/dev/tick` семена появляются после 7 дней. Для быстрого теста — вставьте напрямую:

```bash
# Вставить семя через psql (если запущен docker compose)
docker compose exec db psql -U garden -d garden \
  -c "INSERT INTO seeds (user_id, flower_id, quantity) VALUES (1, 1, 5) ON CONFLICT (user_id, flower_id) DO UPDATE SET quantity = seeds.quantity + 5;"
```

---

## 5. Посадить цветок

```bash
# Посадить flower_id=1 (шаблон из миграции 00002)
curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq .

# Недостаточно семян → 422
curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":999}' | jq .
```

Запомнить `id` посаженного цветка:
```bash
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')
echo "flower id: $FLOWER_ID"
```

---

## 6. Полить цветок

```bash
# Alice поливает свой цветок
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .

# Повторный полив → 409 Already watered
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .

# Bob поливает цветок Alice
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN2" | jq .
```

---

## 7. Цветки другого пользователя

```bash
# Получить цветки по user_id (для страницы чужого сада)
curl -s $BASE/api/flowers/user/1 \
  -H "Authorization: Bearer $TOKEN2" | jq .
```

---

## 8. Dev-тик (+24h)

```bash
# Продвинуть время на сутки вперёд (APP_ENV=development)
curl -s -X POST $BASE/api/dev/tick \
  -H "Authorization: Bearer $TOKEN" | jq .

# Посмотреть обновлённый профиль (FD должен вырасти)
curl -s $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" | jq '.user.fd_balance'
```

---

## 9. Семена

```bash
curl -s $BASE/api/seeds \
  -H "Authorization: Bearer $TOKEN" | jq .

# Поделиться: Alice → Bob, 1 семя flower_id=1
curl -s -X POST $BASE/api/seeds/share \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to_user_id":2,"flower_id":1,"quantity":1}' | jq .

# Недостаточно семян → 422
curl -s -X POST $BASE/api/seeds/share \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to_user_id":2,"flower_id":1,"quantity":9999}' | jq .
```

---

## 10. Уведомления

```bash
curl -s $BASE/api/notifications \
  -H "Authorization: Bearer $TOKEN" | jq .
```

Ожидаемые типы:
- `flower_watered` — кто-то полил ваш цветок
- `seeds_received` — вам поделились семенами
- `flower_dried` — цветок пересох
- `seed_earned` — цветок дал семя (день % 7 == 0)

---

## 11. Лидерборд

```bash
curl -s $BASE/api/leaderboard \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 12. Dev — список пользователей

```bash
curl -s $BASE/api/dev/users \
  -H "Authorization: Bearer $TOKEN" | jq '[.[] | {id, username, fd_balance}]'
```

---

## Полный сценарий за один раз

```bash
BASE=http://localhost:8080

# Регистрация
curl -s -X POST $BASE/api/auth/register -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123","first_name":"Alice"}' > /dev/null

TOKEN=$(curl -s -X POST $BASE/api/auth/login -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')

# Дать семена напрямую в БД
docker compose exec db psql -U garden -d garden \
  -c "INSERT INTO seeds (user_id, flower_id, quantity) VALUES (1, 1, 5) ON CONFLICT DO NOTHING;"

# Посадить
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')

echo "Planted flower: $FLOWER_ID"

# Полить
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .

# Тик
curl -s -X POST $BASE/api/dev/tick \
  -H "Authorization: Bearer $TOKEN" | jq .

# Итог
echo "=== Profile after tick ==="
curl -s $BASE/api/me -H "Authorization: Bearer $TOKEN" | jq '{fd: .user.fd_balance, flowers: (.flowers | length)}'

echo "=== Leaderboard ==="
curl -s $BASE/api/leaderboard -H "Authorization: Bearer $TOKEN" | jq '[.[] | {username, fd_balance}]'
```
