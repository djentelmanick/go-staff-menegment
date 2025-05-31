package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"staff-management/internal/config"
	"staff-management/internal/database"
	"staff-management/internal/handlers"
)

func main() {
	// Загрузка конфигурации
	cfg := config.LoadConfig()
	
	// Инициализация базы данных
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer database.CloseDB()

	// Настройка маршрутов
	r := mux.NewRouter()
	handlers.SetupRoutes(r, cfg)

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{cfg.Server.CORSOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	log.Printf("Сервер запущен на порту %s", cfg.Server.Port)
	log.Printf("Пользователь по умолчанию: %s/%s", cfg.Auth.DefaultLogin, cfg.Auth.DefaultPassword)
	log.Fatal(http.ListenAndServe(":"+cfg.Server.Port, handler))
}