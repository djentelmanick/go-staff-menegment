package config

import (
	"os"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	Auth     AuthConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type ServerConfig struct {
	Port       string
	StaticDir  string
	CORSOrigin string
}

type AuthConfig struct {
	DefaultLogin    string
	DefaultPassword string
	TokenPrefix     string
}

func LoadConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "db"),
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

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (c *Config) GetDatabaseURL() string {
	return "host=" + c.Database.Host +
		" port=" + c.Database.Port +
		" user=" + c.Database.User +
		" password=" + c.Database.Password +
		" dbname=" + c.Database.DBName +
		" sslmode=" + c.Database.SSLMode
}