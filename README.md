# Система управления сотрудниками

Веб-приложение для управления базой данных сотрудников с авторизацией, CRUD операциями и поиском.

## Технологии

- **Backend**: Go (Golang) с функциональной архитектурой
- **Database**: PostgreSQL
- **Frontend**: HTML, CSS, JavaScript (ванильный JS)
- **Authentication**: Простая авторизация с токенами

## Функциональность

- ✅ Авторизация по логину и паролю
- ✅ Просмотр всех сотрудников
- ✅ Добавление новых сотрудников
- ✅ Редактирование данных сотрудников
- ✅ Удаление сотрудников
- ✅ Поиск по всем полям
- ✅ Адаптивный дизайн

## Структура проекта

```
staff-management
│   .env
│   .gitignore
│   database.sql
│   docker-compose.yml
│   Dockerfile
│   go.mod
│   go.sum
│   README.md
│
├───cmd
│   └───server
│           main.go
│
├───internal
│   ├───config
│   │       config.go
│   │
│   ├───database
│   │   │   database.go
│   │   │
│   │   └───migrations
│   ├───handlers
│   │       auth.go
│   │       routes.go
│   │       staff.go
│   │       staff_groups.go
│   │       staff_statuses.go
│   │
│   ├───middleware
│   │       auth.go
│   │
│   └───models
│           staff.go
│           staff_groups.go
│           staff_statuses.go
│           user.go
│
├───pkg
│   └───utils
│           helpers.go
│
└───static
    │   index.html
    │
    ├───css
    │       styles.css
    │
    ├───js
    │       auth.js
    │       dashboard.js
    │       groups.js
    │       logout.js
    │       statuses.js
    │
    └───pages
            benefits.html
            dashboard.html
            groups.html
            login.html
            statuses.html
            vacation.html
```

## Быстрый старт

### Вариант 1: С Docker (рекомендуется)

1. Убедитесь, что у вас установлен Docker и Docker Compose

2. Клонируйте проект и перейдите в директорию:
```bash
cd staff-management
```

3. Добавьте .env в корень проекта
Пример .env
```
# Настройки базы данных
DB_HOST=db
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=staff_db
DB_SSLMODE=disable

# Настройки сервера
SERVER_PORT=8080
STATIC_DIR=./static/
CORS_ORIGIN=*

# Настройки авторизации
DEFAULT_LOGIN=admin
DEFAULT_PASSWORD=admin123
TOKEN_PREFIX=token_
```

4. Запустите приложение:
```bash
docker-compose up --build -d --force-recreate
```

5. Приложение будет доступно по адресу: http://localhost:8080

### Вариант 2: Локальная установка (для первой версии)

#### Требования:
- Go 1.21+
- PostgreSQL 12+

#### Установка:

1. Установите PostgreSQL и создайте базу данных:
```sql
createdb staff_db
```

2. Выполните SQL скрипт для создания таблиц:
```bash
psql staff_db < database.sql
```

3. Установите зависимости Go:
```bash
go mod tidy
```

4. Создайте папку для статических файлов:
```bash
mkdir static
```

5. Поместите файл `index.html` в папку `static/`

6. Запустите приложение:
```bash
go run main.go
```

7. Откройте браузер и перейдите по адресу: http://localhost:8080

## Настройка базы данных

По умолчанию приложение подключается к PostgreSQL со следующими параметрами:
- Host: localhost
- Port: 5432
- Database: staff_db
- User: postgres
- Password: password

Для изменения параметров подключения отредактируйте строку `connStr` в файле `main.go`:

```go
connStr := "user=postgres password=password dbname=staff_db sslmode=disable"
```

## Авторизация

По умолчанию создается пользователь:
- **Логин**: admin
- **Пароль**: admin123

## API Endpoints (для первой версии)

### Авторизация
- `POST /api/login` - Авторизация пользователя

### Сотрудники
- `GET /api/staff` - Получить всех сотрудников
- `POST /api/staff` - Создать нового сотрудника
- `PUT /api/staff/{id}` - Обновить данные сотрудника
- `DELETE /api/staff/{id}` - Удалить сотрудника

Все API endpoints (кроме login) требуют заголовок `Authorization` с токеном.

## Структура данных

### Сотрудник (Staff)
```json
{
  "id": 1,
  "full_name": "Иванов Иван Иванович",
  "phone": "+7 (999) 123-45-67",
  "email": "ivanov@company.com",
  "address": "г. Москва, ул. Тверская, д. 1"
}
```

### Пользователь (User)
```json
{
  "login": "admin",
  "password": "admin123"
}
```

## Безопасность

- Пароли хешируются с использованием bcrypt
- Простая токенная авторизация (в продакшене рекомендуется JWT)
- Валидация входных данных
- Защита от SQL инъекций через prepared statements

## Развитие проекта

Для расширения функциональности можно добавить:

1. **JWT токены** вместо простых токенов
2. **Роли пользователей** (администратор, менеджер, сотрудник)
3. **Пагинацию** для больших списков сотрудников
4. **Экспорт данных** в Excel/CSV
5. **Загрузку файлов** (фото сотрудников, документы)
6. **Логирование** операций
7. **Unit тесты**
8. **Валидацию email** и телефонных номеров
9. **Фильтрацию и сортировку**
10. **Историю изменений**
11. **Доделать страницы с отпуском и надбавками** 

## Отладка

Если возникают проблемы:

1. Проверьте логи Docker:
```bash
docker-compose logs
```

2. Убедитесь, что PostgreSQL запущен:
```bash
docker-compose ps
```

3. Проверьте подключение к БД:
```bash
docker-compose exec postgres psql -U postgres -d staff_db -c "\dt"
```

## Лицензия

MIT License