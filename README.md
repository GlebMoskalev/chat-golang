# Chat API

REST API для управления чатами и сообщениями на Go с использованием PostgreSQL.

## Технологии

- **Go 1.24** - язык программирования
- **net/http + gorilla/mux** - HTTP сервер и роутинг
- **GORM** - ORM для работы с БД
- **PostgreSQL 16** - база данных
- **goose** - миграции
- **Docker & Docker Compose** - контейнеризация

## Архитектура

Проект следует чистой архитектуре с разделением на слои:
- `handler` - HTTP обработчики
- `service` - бизнес-логика
- `repository` - работа с БД
- `models` - модели данных

## Быстрый старт

### Требования
- Docker
- Docker Compose

### Запуск

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd chat-golang
```

2. Создайте `.env` файл (или используйте `.env.example`):
```bash
cp .env.example .env
```

3. Запустите проект:
```bash
docker-compose up --build
```

API будет доступен по адресу: `http://localhost:8080`

### Остановка

```bash
docker-compose down
```

Для удаления данных:
```bash
docker-compose down -v
```

## API Endpoints

### 1. Создать чат

```bash
POST /chats/
Content-Type: application/json

{
  "title": "Мой чат"
}
```

**Response (201):**
```json
{
  "id": 1,
  "title": "Мой чат",
  "created_at": "2026-01-28T10:30:00Z"
}
```

**Валидация:**
- `title` обязателен
- Длина: 1-200 символов
- Пробелы по краям удаляются автоматически

### 2. Получить чат с сообщениями

```bash
GET /chats/{id}?limit=20
```

**Query параметры:**
- `limit` - количество последних сообщений (по умолчанию 20, максимум 100)

**Response (200):**
```json
{
  "id": 1,
  "title": "Мой чат",
  "created_at": "2026-01-28T10:30:00Z",
  "messages": [
    {
      "id": 2,
      "chat_id": 1,
      "text": "Второе сообщение",
      "created_at": "2026-01-28T10:31:00Z"
    },
    {
      "id": 1,
      "chat_id": 1,
      "text": "Первое сообщение",
      "created_at": "2026-01-28T10:30:30Z"
    }
  ]
}
```

**Примечание:** Сообщения отсортированы по `created_at` DESC (новые первыми)

### 3. Отправить сообщение в чат

```bash
POST /chats/{id}/messages/
Content-Type: application/json

{
  "text": "Привет, мир!"
}
```

**Response (201):**
```json
{
  "id": 1,
  "chat_id": 1,
  "text": "Привет, мир!",
  "created_at": "2026-01-28T10:30:30Z"
}
```

**Валидация:**
- `text` обязателен
- Длина: 1-5000 символов
- Пробелы по краям удаляются автоматически
- Чат должен существовать (иначе 404)

### 4. Удалить чат

```bash
DELETE /chats/{id}
```

**Response (204):** No Content

**Примечание:** Все сообщения чата удаляются каскадно

## Примеры использования

### Создание чата и отправка сообщений

```bash
# Создать чат
curl -X POST http://localhost:8080/chats/ \
  -H "Content-Type: application/json" \
  -d '{"title":"Тестовый чат"}'

# Отправить сообщение
curl -X POST http://localhost:8080/chats/1/messages/ \
  -H "Content-Type: application/json" \
  -d '{"text":"Привет!"}'

# Получить чат с сообщениями
curl http://localhost:8080/chats/1?limit=10

# Удалить чат
curl -X DELETE http://localhost:8080/chats/1
```

## Разработка

### Запуск тестов

```bash
go test ./...
```

### Запуск с покрытием

```bash
go test -cover ./...
```

### Создание новой миграции

```bash
goose -dir migrations create migration_name sql
```

### Применение миграций вручную

```bash
goose -dir migrations postgres "host=localhost user=postgres password=postgres dbname=chat port=5432 sslmode=disable" up
```

## Структура проекта

```
.
├── cmd/
│   └── app/
│       └── main.go           # Точка входа
├── internal/
│   ├── handler/              # HTTP обработчики
│   ├── service/              # Бизнес-логика
│   ├── repository/           # Работа с БД
│   └── models/               # Модели данных
├── migrations/               # SQL миграции
├── docker-compose.yml        # Docker Compose конфигурация
├── Dockerfile                # Dockerfile для приложения
├── .env                      # Переменные окружения
└── README.md                 # Документация
```

## Переменные окружения

```env
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_DB=chat

DB_HOST=postgres
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=chat
DB_PORT=5432
```

## Особенности реализации

- **Каскадное удаление**: При удалении чата все сообщения удаляются автоматически через `ON DELETE CASCADE`
- **Валидация**: Все входные данные валидируются на уровне сервиса
- **Trim**: Пробелы по краям `title` и `text` удаляются автоматически
- **Индексы**: Добавлены индексы для оптимизации запросов по `chat_id` и сортировке
- **Health check**: PostgreSQL проверяется перед запуском миграций и приложения

## Лицензия

MIT
