package main

import (
	"os"
	"strconv"
)

// Config содержит конфигурацию приложения
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
}

// DatabaseConfig содержит настройки базы данных
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// ServerConfig содержит настройки сервера
type ServerConfig struct {
	Port       string
	StaticDir  string
	CORSOrigin string
}

// AuthConfig содержит настройки авторизации
type AuthConfig struct {
	DefaultLogin    string
	DefaultPassword string
	TokenPrefix     string
}

// LoadConfig загружает конфигурацию из переменных окружения или использует значения по умолчанию
func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			DBName:   getEnv("DB_NAME", "staff_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Server: ServerConfig{
			Port:       getEnv("SERVER_PORT", "8080"),
			StaticDir:  getEnv("STATIC_DIR", "./static/"),
			CORSOrigin: getEnv("CORS_ORIGIN", "*"),
		},
		Auth: AuthConfig{
			DefaultLogin:    getEnv("DEFAULT_LOGIN", "admin"),
			DefaultPassword: getEnv("DEFAULT_PASSWORD", "admin123"),
			TokenPrefix:     getEnv("TOKEN_PREFIX", "token_"),
		},
	}
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// getEnvAsInt получает значение переменной окружения как integer
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool получает значение переменной окружения как boolean
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// GetDatabaseURL возвращает строку подключения к базе данных
func (c *Config) GetDatabaseURL() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}
