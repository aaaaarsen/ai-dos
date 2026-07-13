# AI Dos — Персональный журнал мыслей с ИИ-наставником

AI Dos — это мобильное приложение для саморефлексии. Пользователь ведёт дневник мыслей в формате диалога с ИИ-ассистентом, который помогает замечать эмоциональные паттерны и лучше понимать себя.

> Проект разработан в рамках ESG Internship Summer 2026 в команде из четырёх человек: бэкенд (Go) + iOS (SwiftUI).

---

## Возможности

- Регистрация и авторизация через email + пароль (JWT)
- Чат с ИИ-наставником на базе Llama 3.3 70B (Groq API)
- Автоматическая генерация сводок каждые 10 сообщений — ИИ "запоминает" прошлые темы
- Персональные инсайты — анализ эмоциональных паттернов по всем чатам
- Трекер настроения по дням
- Статистика активности за неделю для визуального графика
- Полное удаление аккаунта и всех данных
- Swagger UI — интерактивная документация API

---

## Стек

### Backend
- **Go 1.25** — основной язык
- **Gin** — HTTP-роутер
- **PostgreSQL 16** — база данных
- **pgx/v5** — драйвер и пул соединений к Postgres
- **golang-migrate** — версионирование схемы БД
- **golang-jwt** — генерация и валидация JWT-токенов
- **bcrypt** — хеширование паролей
- **Groq API** (Llama 3.3 70B) — ИИ-ответы и генерация сводок
- **Docker + Docker Compose + nginx** — контейнеризация и reverse proxy
- **swaggo/swag** — генерация Swagger документации

### iOS
- **SwiftUI** — UI-фреймворк
- **MVVM** — архитектурный паттерн
- **URLSession + Codable** — сетевой слой

---

## Деплой

Бэкенд задеплоен на [Render](https://render.com):
- **API Base URL:** `https://ai-dos-api.onrender.com/api`
- **Swagger UI:** `https://ai-dos-api.onrender.com/swagger/index.html`
- Auto-deploy при push в ветку `main`

> ⚠️ Render Free tier: первый запрос после простоя может занять 30–50 секунд (холодный старт). Последующие запросы работают в штатном режиме.

---

## API

Базовый URL: `https://ai-dos-api.onrender.com/api`

Все защищённые эндпоинты (✓) требуют заголовок:
```
Authorization: Bearer <token>
```

| Метод | URL | Описание | Auth |
|-------|-----|----------|------|
| GET | /health | Проверка сервера | — |
| POST | /auth/register | Регистрация | — |
| POST | /auth/login | Вход | — |
| GET | /users/me | Профиль | ✓ |
| DELETE | /users/me | Удалить аккаунт | ✓ |
| GET | /users/me/insights | Анализ паттернов (ИИ) | ✓ |
| GET | /users/me/stats | Статистика за 7 дней | ✓ |
| POST | /mood | Сохранить настроение | ✓ |
| GET | /mood/today | Настроение сегодня | ✓ |
| POST | /chats | Создать чат | ✓ |
| GET | /chats | Список чатов | ✓ |
| DELETE | /chats/:id | Удалить чат | ✓ |
| POST | /chats/:id/messages | Отправить сообщение → ИИ отвечает | ✓ |
| GET | /chats/:id/messages | История сообщений | ✓ |
| GET | /chats/:id/summaries | Сводки чата | ✓ |

Полная интерактивная документация: **https://ai-dos-api.onrender.com/swagger/index.html**

---

## Архитектура бэкенда

```
cmd/
  api/
    main.go           # точка входа: миграции, пул соединений, роутер
docs/                 # сгенерированная Swagger документация
internal/
  ai/
    groq.go           # HTTP-клиент для Groq API
    prompts.go        # системные промпты (наставник + генерация сводок)
  auth/
    jwt.go            # генерация и валидация JWT
    password.go       # хеширование bcrypt
  db/
    db.go             # инициализация pgxpool
    migrate.go        # автоматическое применение миграций при старте
    users.go          # CRUD пользователей
    chats.go          # CRUD чатов
    messages.go       # CRUD сообщений + статистика
    summaries.go      # создание и получение сводок
    moods.go          # трекер настроения
  handlers/
    auth.go           # register, login, profile, delete account
    chats.go          # create/get/delete chat, summaries
    messages.go       # send message (→ AI reply), get history
    mood.go           # save/get mood
    insights.go       # weekly stats, AI insights
  middleware/
    auth.go           # JWT middleware для защищённых маршрутов
  models/
    user.go
    chat.go
    message.go
    summary.go
    mood.go
    stat.go
migrations/
  000001_create_users_table.up.sql
  000002_create_chats_table.up.sql
  000003_create_messages_table.up.sql
  000004_create_summaries_table.up.sql
  000005_add_name_to_users.up.sql
  000006_create_moods_table.up.sql
```

---

## Схема базы данных

```
users
  id, name, email, password_hash, created_at

chats
  id, user_id → users(id), title, created_at

messages
  id, chat_id → chats(id), role (user/assistant/system), content, created_at

summaries
  id, chat_id → chats(id), content, created_at

moods
  id, user_id → users(id), emoji, date, created_at
  UNIQUE(user_id, date)
```

Каскадное удаление: удаление пользователя удаляет все его чаты → сообщения → сводки.

---

## Как работает ИИ-память

1. Каждое сообщение пользователя сохраняется в базу с ролью `user`
2. Сервер собирает контекст: системный промпт + последние 3 сводки + последние 10 сообщений
3. Контекст отправляется в Groq API → модель генерирует ответ
4. Ответ сохраняется с ролью `assistant` и возвращается клиенту
5. Каждые 10 сообщений: сервер автоматически генерирует сводку диалога и сохраняет её в `summaries` — это "долгосрочная память" ассистента

---

## Локальный запуск

### Требования
- Go 1.25+
- Docker + Docker Compose

### Через Docker (рекомендуется)

```bash
git clone https://github.com/aaaaarsen/ai-dos.git
cd ai-dos/ai-dos
cp .env.example .env   # заполни переменные
docker compose up --build
```

### Переменные окружения (.env)

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=ai_dos_dev
DB_SSLMODE=disable

SERVER_PORT=8080
JWT_SECRET=your_jwt_secret

GROQ_API_KEY=your_groq_api_key
GROQ_MODEL=llama-3.3-70b-versatile
```

### Обновление Swagger документации

```bash
swag init -g cmd/api/main.go
docker compose up --build
```

---

## Команда

- **Мукажанов Арлан** — Team Lead
- **Нурдинов Ильхан** — UI/UX Designer
- **Райымбек Арсен** — Backend (Go) · [github.com/aaaaarsen](https://github.com/aaaaarsen)
- **Мырзахметов Асан** — iOS (SwiftUI) · [github.com/asaaaanmyrza](https://github.com/asaaaanmyrza)