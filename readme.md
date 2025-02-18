# Avito merch
Решение тестового задания для Avito на позицию Intern Backend разработчик.

# Инструкция по запуску
Запускаем докер контейнер с проектом и бд
```bash
docker-compose up --build
```

/////////////////////////////

# API
Реализованы методы для работы с магазином и системой монет.

## Получение информации о монетах, инвентаре и истории транзакций
**Эндпоинт:** GET "/api/info"

Возвращает информацию о балансе монет, инвентаре и истории транзакций.

### Пример запроса:
```bash
curl --request GET \
  --url http://localhost:8080/api/info \
  --header "Authorization: Bearer <TOKEN>"
```

## Отправка монет пользователю
**Эндпоинт:** POST "/api/sendCoin"

Позволяет отправить монеты другому пользователю.

### Пример запроса:
```bash
curl --request POST \
  --url http://localhost:8080/api/sendCoin \
  --header "Authorization: Bearer <TOKEN>" \
  --header "Content-Type: application/json" \
  --data '{
    "toUser": "user123",
    "amount": 10
  }'
```

## Покупка предмета за монеты
**Эндпоинт:** GET "/api/buy/{item}"

Позволяет приобрести предмет за монеты.

### Пример запроса:
```bash
curl --request GET \
  --url http://localhost:8080/api/buy/book \
  --header "Authorization: Bearer <TOKEN>"
```
### Cписок доступных в базе наименований и их цены.

| Название     | Цена |
|--------------|------|
| t-shirt      | 80   |
| cup          | 20   |
| book         | 50   |
| pen          | 10   |
| powerbank    | 200  |
| hoody        | 300  |
| umbrella     | 200  |
| socks        | 10   |
| wallet       | 50   |
| pink-hoody   | 500  |

Предполагается, что в магазине бесконечный запас каждого вида мерча.

## Аутентификация и получение JWT-токена
**Эндпоинт:** POST "/api/auth"

Позволяет получить JWT-токен для доступа к API.

### Пример запроса:
```bash
curl --request POST \
  --url http://localhost:8080/api/auth \
  --header "Content-Type: application/json" \
  --data '{
    "username": "testuser",
    "password": "password123"
  }'
```

### Пример успешного ответа:
```json
{
  "token": "eyJhbGciOiJIUzI1..."
}
```

Для всех защищенных эндпоинтов необходимо передавать токен в заголовке:
```bash
Authorization: Bearer <TOKEN>
```

# Тестирование
Реализованы интеграционные и E2E тесты:
- Тест на сценарий покупки мерча
- Тест на сценарий передачи монеток другим сотрудникам

Файл с тестами: `./test/e2e_test.go`

Также покрыты unit-тестами слои:
- Service
- Repository