package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"staff-management/internal/config"
	"staff-management/internal/database"
	"staff-management/internal/middleware"
)

func SetupRoutes(r *mux.Router, cfg *config.Config) {
	db := database.GetDB()
	
	// API маршруты
    api := r.PathPrefix("/api").Subrouter()
    
    // Авторизация
    api.HandleFunc("/login", LoginHandler(db, cfg)).Methods("POST")
    
    // CRUD для сотрудников
    api.HandleFunc("/staff", middleware.AuthMiddleware(GetStaffHandler(db))).Methods("GET")
    api.HandleFunc("/staff", middleware.AuthMiddleware(CreateStaffHandler(db))).Methods("POST")
    api.HandleFunc("/staff/{id}", middleware.AuthMiddleware(UpdateStaffHandler(db))).Methods("PUT")
    api.HandleFunc("/staff/{id}", middleware.AuthMiddleware(DeleteStaffHandler(db))).Methods("DELETE")

    // Обслуживание статических файлов
    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/assets"))))
    
    // Главная страница
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    })
}