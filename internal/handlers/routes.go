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

    // CRUD для групп сотрудников
    api.HandleFunc("/groups", middleware.AuthMiddleware(GetStaffGroups(db))).Methods("GET")
    api.HandleFunc("/groups", middleware.AuthMiddleware(CreateStaffGroup(db))).Methods("POST")
    api.HandleFunc("/groups/{id}", middleware.AuthMiddleware(UpdateStaffGroup(db))).Methods("PUT")
    api.HandleFunc("/groups/{id}", middleware.AuthMiddleware(DeleteStaffGroup(db))).Methods("DELETE")
    api.HandleFunc("/groups/{group_id}/members", middleware.AuthMiddleware(GetGroupMemebers(db))).Methods("GET")
    api.HandleFunc("/groups/{group_id}/members", middleware.AuthMiddleware(AddMemebersToGroup(db))).Methods("POST")
    api.HandleFunc("/groups/{group_id}/members/{staff_id}", middleware.AuthMiddleware(DeleteGroupMember(db))).Methods("DELETE")

    //CRUD для статусов сотрудников
    api.HandleFunc("/employees/statuses", middleware.AuthMiddleware(GetStaffStatuses(db))).Methods("GET")
    api.HandleFunc("/employees/{id}/status", middleware.AuthMiddleware(UpdateStaffStatus(db))).Methods("PUT")

    // Обслуживание статических файлов
    r.PathPrefix("/css/").Handler(http.StripPrefix("/css/", http.FileServer(http.Dir("./static/css"))))
    r.PathPrefix("/js/").Handler(http.StripPrefix("/js/", http.FileServer(http.Dir("./static/js"))))
    r.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./static/assets"))))
    r.PathPrefix("/pages/").Handler(http.StripPrefix("/pages/", http.FileServer(http.Dir("./static/pages"))))
    
    // Страницы 
    r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/login.html")
    }).Methods("GET")

    r.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/dashboard.html")
    }).Methods("GET")

    r.HandleFunc("/statuses", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/statuses.html")
    }).Methods("GET")

    r.HandleFunc("/groups", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/groups.html")
    }).Methods("GET")

    r.HandleFunc("/benefits", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/benefits.html")
    }).Methods("GET")

    r.HandleFunc("/vacation", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/pages/vacation.html")
    }).Methods("GET")

    // Главная страница
    r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./static/index.html")
    })
}