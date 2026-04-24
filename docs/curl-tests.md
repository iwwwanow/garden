# curl-tests — Garden of the Goddess of Flowers

Полный сценарий: регистрация → логин → посадить → полить → лидерборд → семена → уведомления.

Предварительно запустите сервер:

```bash
make dev   # или: cd server && air
```

Для чистого старта:

```bash
make db-clean   # сбросить пользовательские данные (схема и шаблоны цветков остаются)
# или
make db-reset   # полный сброс + миграции
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

echo "Alice: $TOKEN"
echo "Bob:   $TOKEN2"

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

---

## 4. Дать семена (dev-эндпоинт)

```bash
# Дать 5 семян flower_id=1 текущему пользователю
curl -s -X POST $BASE/api/dev/seeds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1,"quantity":5}' | jq .
```

---

## 5. Посадить цветок

```bash
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')
echo "flower id: $FLOWER_ID"

# Недостаточно семян → 422
curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":999}' | jq .
```

---

## 6. Полить цветок

> **Правило:** каждый цветок можно полить ровно один раз в день — кем угодно.  
> Первый полив за день «закрывает» цветок. Любая следующая попытка → 409.

```bash
# Alice поливает свой цветок
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .
# → {"status":"ok"}

# Любая повторная попытка (хоть Alice, хоть Bob) → 409 Already watered
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN2" | jq .
# → {"error":"already watered today"}

# Bob поливает ДРУГОЙ цветок Alice (если у неё их несколько)
# curl -s -X POST $BASE/api/flowers/$FLOWER_ID2/water -H "Authorization: Bearer $TOKEN2" | jq .
```

---

## 7. Чужой сад

```bash
# Цветки пользователя с id=2 (для страницы чужого сада)
curl -s $BASE/api/flowers/user/2 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 8. Dev — тик (+24h)

```bash
# Симулировать полночь: пересохшие цветки → dried, политые → day++, FD++
curl -s -X POST $BASE/api/dev/tick \
  -H "Authorization: Bearer $TOKEN" | jq .

curl -s $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" | jq '{fd: .user.fd_balance, flowers: (.flowers | map({id, day, is_dried}))}'
```

---

## 9. Семена

```bash
curl -s $BASE/api/seeds \
  -H "Authorization: Bearer $TOKEN" | jq .

# Alice делится с Bob
curl -s -X POST $BASE/api/seeds/share \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"to_user_id":2,"flower_id":1,"quantity":1}' | jq .

# Нет столько → 422
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

Возможные типы уведомлений и условия их появления:

| Тип              | Когда                                          |
|------------------|------------------------------------------------|
| `flower_watered` | Кто-то (не вы) полил ваш цветок               |
| `seeds_received` | Вам поделились семенами (`/api/seeds/share`)   |
| `flower_dried`   | Тик: цветок не был полит вчера                 |
| `seed_earned`    | Тик: цветок достиг дня кратного 7             |

Чтобы увидеть все четыре типа в одном прогоне:
```bash
# 1. Зарегистрировать двух пользователей
# 2. Alice сажает цветок, Bob поливает → flower_watered у Alice
# 3. Alice делится семенами с Bob → seeds_received у Bob
# 4. Запустить тик, не полив какой-то цветок → flower_dried
# 5. Запустить 7 тиков с поливом → seed_earned
```

---

## 11. Лидерборд

```bash
curl -s $BASE/api/leaderboard \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 12. Dev — служебные

```bash
# Список всех пользователей
curl -s $BASE/api/dev/users \
  -H "Authorization: Bearer $TOKEN" | jq '[.[] | {id, username, fd_balance}]'

# Дать семена текущему пользователю
curl -s -X POST $BASE/api/dev/seeds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1,"quantity":10}' | jq .

# Сбросить все пользовательские данные (не трогает flower templates)
curl -s -X POST $BASE/api/dev/reset \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## Полный сценарий за один раз

```bash
BASE=http://localhost:8080

# Сброс данных
curl -s -X POST $BASE/api/dev/reset -H "Authorization: Bearer $TOKEN" > /dev/null 2>&1 || true

# Регистрация
curl -s -X POST $BASE/api/auth/register -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123","first_name":"Alice"}' > /dev/null
curl -s -X POST $BASE/api/auth/register -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pass456","first_name":"Bob"}' > /dev/null

TOKEN=$(curl -s -X POST $BASE/api/auth/login -H "Content-Type: application/json" \
  -d '{"username":"alice","password":"pass123"}' | jq -r '.token')
TOKEN2=$(curl -s -X POST $BASE/api/auth/login -H "Content-Type: application/json" \
  -d '{"username":"bob","password":"pass456"}' | jq -r '.token')

# Дать семена через dev-эндпоинт
curl -s -X POST $BASE/api/dev/seeds \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"flower_id":1,"quantity":5}' > /dev/null

# Посадить
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')
echo "Planted flower: $FLOWER_ID"

# Bob поливает цветок Alice (социальный полив)
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN2" | jq .

# Alice пытается полить тот же цветок → 409
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .

# Тик
curl -s -X POST $BASE/api/dev/tick \
  -H "Authorization: Bearer $TOKEN" | jq .

# Итог
echo "=== Alice after tick ==="
curl -s $BASE/api/me -H "Authorization: Bearer $TOKEN" | \
  jq '{fd: .user.fd_balance, flowers: (.flowers | map({id, day, is_dried}))}'

echo "=== Alice notifications ==="
curl -s $BASE/api/notifications -H "Authorization: Bearer $TOKEN" | jq .

echo "=== Leaderboard ==="
curl -s $BASE/api/leaderboard -H "Authorization: Bearer $TOKEN" | \
  jq '[.[] | {username, fd_balance}]'
```
