# URL-SHORTENER
URL Shortener Service
![Go](https://img.shields.io/badge/Go-1.24+-blue) ![Redis](https://img.shields.io/badge/Redis-7.0+-red) ![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue)


## Сервис для сокращения URL-cсылок, построенный на Go с использованием PostgreSQL  и Redis

## Структура проекта

url-shortener-go/
    cmd/main.go                           # Точка входа
    internal/
        config/config.go                  # Конфигурация
        handlers/url_handlers.go          # HTTP обработчики
        models/models.go                  # Модели данных
        repositories/repository.go        # Работа с БД
        services/service.go               # Бизнес-логика
        pkg/random/random.go              # Генрация коротких ссылок
Dockerfile
docker-compose.yml
init.sql                                  # Инициализация БД
go.mod
README.md

## Возможности

- Сокращение длинных URL в короткие ссылки
- Настраиваемое время жизни ссылок (по умолчанию 30 дней)
- Кэширование в Redis для быстрого доступа
- Подсчет переходов по ссылкам
- Безопасная генерация коротких кодов (криптостойкий random)
- Docker-контейнеризация

## Технологии

- Go 1.24 - основной язык разработки
- PostgreSQL 15 - основное хранилище данных
- Redis 7 - кэширование и быстрый доступ
- Gin - HTTP фреймворк
- Docker - контейнеризация
- pgx - драйвер для PostgreSQL
- go-redis - клиент для Redis

## Требования

- Go 1.24 или выше
- Docker и Docker Compose
- PostgreSQL 15 (если не используется Docker)
- Redis 7 (если не используется Docker)

## Быстрый старт

## Запуск через Docker

1. Клонируйте репозиторий:


```bash
git clone https://github.com/EvKrpv/url-shortener-go.git
cd url-shortener-go
```
Для более простой настройки можете использовать файл `docker-compose.yml`:

2. Запустите сервис: 

```bash
docker-compose up -d
```
При этом запустятся приложение Go, база данных PostgreSQL и Redis.

3. Проверьте работу 

```bash
curl http://localhost:8080/health
```

## Локальный запуск без Docker

1. Установите зависимости:

```bash
go mod download
```

2. Создайте файл .env:

```.env
DB_URL=postgres://user:password@localhost:5432/urls?sslmode=disable
REDIS_URL=redis://localhost:6379/0
PORT=8080
ENVIRONMENT=development
```

3. Создайте таблицу в PostgreSQL:

```sql
CREATE TABLE IF NOT EXISTS urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    user_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '30 days'),
    click_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT TRUE
);
```

4. Запустите приложение

```bash
go run cmd/main.go
```


## API Документация

- POST Создание короткой ссылки

```json
{
    "url": "https://www.example.com/very/long/url",
    "expires_in": "720h"  // опционально, по умолчанию 720h (30 дней)
}
```
Успешный ответ

```json
{
    "short_url": "http://localhost:8080/abc123defg"
}
```
Пример: 

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'
```

- GET редирект на оригинальный URL

```bash
curl -v http://localhost:8080/abc123defg
# Ответ: 301 Moved Permanently с редиректом
```

- GET информация о ссылке

```bash
curl http://localhost:8080/info/abc123defg
```
- Ответ

```json
{
    "short_code": "abc123defg",
    "original_url": "https://www.google.com",
    "created_at": "2024-01-01T12:00:00Z",
    "expires_at": "2024-01-31T12:00:00Z",
    "click_count": 42
}
```

## Ручное тестирование

```bash
# Создание ссылки
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://www.google.com"}'

# Переход по ссылке (замените CODE на полученный)
curl -v http://localhost:8080/CODE

# Информация о ссылке
curl http://localhost:8080/info/CODE
```

## Производительность
- Кэширование Redis: все запросы сначала проверяются в Redis (латентность <1мс)
- Пул соединений: pgxpool с настраиваемым количеством соединений к PostgreSQL
- Генерация кодов: криптостойкий random (63^10 ≈ 9.85×10^17 комбинаций)
- Масштабирование: горизонтальное масштабирование через добавление экземпляров приложения

## Безопасность
- Генерация: криптостойкий crypto/rand для непредсказуемых кодов
- Валидация: проверка URL перед сохранением в БД
- SQL: параметризованные запросы (защита от инъекций)
- Очистка: автоматическое удаление просроченных ссылок

## Планы на будущее

- Веб-интерфейс (удобный сайт вместо curl)
- Панель мониторинга для анализа ссылок.
- Настраиваемые алиасы.
- API для интеграций (чтобы другие сервисы могли использовать)
- Интеграция с Prometheus/Grafana для мониторинга в режиме реального времени.
