# URL-shortener

## Описание:
REST API сервис на Go для превращения длинных и сложных интернет-адресов в короткую и удобную ссылку
Сервис позволяет:
- создавать короткие ссылки;

- перенаправлять пользователя по короткой ссылке на исходный URL;

- изменять сохранённые URL;

- удалять ссылки.

---
## Технологии:

- Go
- PostgreSQL
- Docker/Docker Compose

---
## Установка
1. Клонировать репозиторий:

`git clone https://github.com/Nonlight/url-shortener.git`

2. Перейти в директорию проекта:

`cd url-shortener`

3. Скопировать файл-шаблон `.env.example`:

`cp .env.example .env`

4. Запустить контейнеры:

`docker compose up`

После запуска сервис доступен по адресу `http://localhost:8082`

---
## HTTP-Методы
````
POST   /url          — создать короткую ссылку
GET    /{alias}      — открыть исходную ссылку по её alias
PUT    /url          — изменить URL
DELETE /url/{alias}  — удалить ссылку
````
---
## Тесты
В проекте реализованы:
- unit-тесты для HTTP-handlers (`save`, `redirect`, `update`, `delete`);
- интеграционные тесты для проверки работы API.

Перед запуском интеграционных тестов сервис должен быть запущен:
`docker compose up`

Для запуска всех тестов:
`go test ./...`

---
## Автор
[Nonlight](https://github.com/Nonlight)