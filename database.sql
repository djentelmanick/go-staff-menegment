-- Создание базы данных
CREATE DATABASE staff_db;

-- Подключение к базе данных
\c staff_db;

-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание таблицы сотрудников
CREATE TABLE IF NOT EXISTS staff (
    id SERIAL PRIMARY KEY,
    full_name VARCHAR(255) NOT NULL,
    phone VARCHAR(20),
    email VARCHAR(100),
    address TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создание триггера для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_staff_updated_at 
    BEFORE UPDATE ON staff 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

-- Создание индексов для оптимизации
CREATE INDEX IF NOT EXISTS idx_staff_full_name ON staff(full_name);
CREATE INDEX IF NOT EXISTS idx_staff_email ON staff(email);
CREATE INDEX IF NOT EXISTS idx_users_login ON users(login);

-- Вставка администратора по умолчанию (пароль: admin123)
-- Хеш создан с помощью bcrypt
INSERT INTO users (login, password_hash) 
VALUES ('admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi')
ON CONFLICT (login) DO NOTHING;

-- Пример данных сотрудников для тестирования
INSERT INTO staff (full_name, phone, email, address) VALUES
('Иванов Иван Иванович', '+7 (999) 123-45-67', 'ivanov@company.com', 'г. Москва, ул. Тверская, д. 1'),
('Петрова Анна Сергеевна', '+7 (999) 234-56-78', 'petrova@company.com', 'г. Санкт-Петербург, Невский пр., д. 100'),
('Сидоров Петр Александрович', '+7 (999) 345-67-89', 'sidorov@company.com', 'г. Екатеринбург, ул. Ленина, д. 50')
ON CONFLICT DO NOTHING;
