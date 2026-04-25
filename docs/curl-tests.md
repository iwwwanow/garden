# curl-tests — Garden of the Goddess of Flowers

Полный сценарий: регистрация → логин → шаблоны → посадить → полить → тик → профиль → уведомления → лидерборд.

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
TOKEN=""        # заполняется после логина Alice
TOKEN2=""       # второй пользователь Bob
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
# Получить профиль (user + flowers + seeds)
curl -s $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" | jq .

# Обновить first_name
curl -s -X PUT $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":"Alicia"}' | jq .

# Пустое имя → 400
curl -s -X PUT $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"first_name":""}' | jq .
```

---

## 4. Шаблоны цветков

```bash
# Список всех шаблонов (что можно посадить)
curl -s $BASE/api/flowers \
  -H "Authorization: Bearer $TOKEN" | jq .

# Конкретный шаблон по ID
curl -s $BASE/api/flowers/1 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Несуществующий → 404
curl -s $BASE/api/flowers/999 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 5. Дать семена (dev-эндпоинт)

```bash
# Дать 5 семян flower_id=1 текущему пользователю
curl -s -X POST $BASE/api/dev/seeds \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1,"quantity":5}' | jq .
```

---

## 6. Посадить цветок

```bash
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')
echo "flower id: $FLOWER_ID"

# Несуществующий flower_id → 422
curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"flower_id":999}' | jq .
```

---

## 7. Полить цветок

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
```

---

## 8. Чужой сад

```bash
# Цветки пользователя 2 (Bob) — для страницы чужого сада
curl -s $BASE/api/flowers/user/2 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Публичный профиль пользователя 2
curl -s $BASE/api/users/2 \
  -H "Authorization: Bearer $TOKEN" | jq .

# Несуществующий пользователь → 404
curl -s $BASE/api/users/999 \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 9. Dev — тик (+24h)

```bash
# Симулировать полночь: пересохшие цветки → dried, политые → day++, FD++
curl -s -X POST $BASE/api/dev/tick \
  -H "Authorization: Bearer $TOKEN" | jq .

curl -s $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" | jq '{fd: .user.fd_balance, flowers: (.flowers | map({id, day, is_dried}))}'
```

---

## 10. Семена

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

## 11. Уведомления

```bash
# Получить список
curl -s $BASE/api/notifications \
  -H "Authorization: Bearer $TOKEN" | jq .

# Отметить все прочитанными
curl -s -X PATCH $BASE/api/notifications/read \
  -H "Authorization: Bearer $TOKEN" | jq .

# Убедиться, что is_read=true
curl -s $BASE/api/notifications \
  -H "Authorization: Bearer $TOKEN" | jq '[.[] | {id, type, is_read}]'
```

Возможные типы уведомлений:

| Тип              | Когда                                        |
| ---------------- | -------------------------------------------- |
| `flower_watered` | Кто-то (не вы) полил ваш цветок              |
| `seeds_received` | Вам поделились семенами (`/api/seeds/share`) |
| `flower_dried`   | Тик: цветок не был полит вчера               |
| `seed_earned`    | Тик: цветок достиг дня кратного 7            |

---

## 12. Лидерборд

```bash
# Топ-100 (по умолчанию)
curl -s $BASE/api/leaderboard \
  -H "Authorization: Bearer $TOKEN" | jq '[.[] | {username, fd_balance}]'

# Пагинация: первые 10
curl -s "$BASE/api/leaderboard?limit=10&offset=0" \
  -H "Authorization: Bearer $TOKEN" | jq .

# Следующая страница
curl -s "$BASE/api/leaderboard?limit=10&offset=10" \
  -H "Authorization: Bearer $TOKEN" | jq .
```

---

## 13. Dev — служебные

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

# Шаблоны цветков
echo "=== Flower templates ==="
curl -s $BASE/api/flowers -H "Authorization: Bearer $TOKEN" | jq '[.[] | {id, season, image_path}]'

# Дать семена через dev-эндпоинт
curl -s -X POST $BASE/api/dev/seeds \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"flower_id":1,"quantity":5}' > /dev/null

# Посадить
FLOWER_ID=$(curl -s -X POST $BASE/api/flowers/plant \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"flower_id":1}' | jq -r '.id')
echo "Planted flower: $FLOWER_ID"

# Обновить имя
curl -s -X PUT $BASE/api/me \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"first_name":"Alicia"}' | jq '{id, first_name}'

# Bob поливает цветок Alice
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN2" | jq .

# Alice пытается полить тот же цветок → 409
curl -s -X POST $BASE/api/flowers/$FLOWER_ID/water \
  -H "Authorization: Bearer $TOKEN" | jq .

# Alice делится семенами с Bob
curl -s -X POST $BASE/api/seeds/share \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json" \
  -d '{"to_user_id":2,"flower_id":1,"quantity":1}' > /dev/null

# Тик
curl -s -X POST $BASE/api/dev/tick -H "Authorization: Bearer $TOKEN" | jq .

# Отметить уведомления прочитанными
curl -s -X PATCH $BASE/api/notifications/read -H "Authorization: Bearer $TOKEN" | jq .

# Итог
echo "=== Alice after tick ==="
curl -s $BASE/api/me -H "Authorization: Bearer $TOKEN" | \
  jq '{fd: .user.fd_balance, name: .user.first_name, flowers: (.flowers | map({id, day, is_dried}))}'

echo "=== Alice notifications (all read) ==="
curl -s $BASE/api/notifications -H "Authorization: Bearer $TOKEN" | \
  jq '[.[] | {type, is_read}]'

echo "=== Bob public profile ==="
curl -s $BASE/api/users/2 -H "Authorization: Bearer $TOKEN" | jq '{id, username, fd_balance}'

echo "=== Leaderboard (top 10) ==="
curl -s "$BASE/api/leaderboard?limit=10&offset=0" -H "Authorization: Bearer $TOKEN" | \
  jq '[.[] | {username, fd_balance}]'
```
