package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Структуры данных
type Staff struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Email    string `json:"email"`
	Address  string `json:"address"`
}

type User struct {
	ID       int    `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

// Глобальные переменные
var db *sql.DB
var config *Config

// Функции для работы с базой данных
func initDB() error {
	var err error
	db, err = sql.Open("postgres", config.GetDatabaseURL())
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	// Создание таблиц
	createTables()
	
	// Создание админа по умолчанию
	createDefaultAdmin()

	return nil
}

func createTables() {
	// Таблица пользователей
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL
	)`

	// Таблица сотрудников
	staffTable := `
	CREATE TABLE IF NOT EXISTS staff (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		email VARCHAR(100),
		address TEXT
	)`

	db.Exec(userTable)
	db.Exec(staffTable)
}

func createDefaultAdmin() {
	// Проверяем, есть ли уже админ
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1", config.Auth.DefaultLogin).Scan(&count)
	
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(config.Auth.DefaultPassword), bcrypt.DefaultCost)
		db.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", config.Auth.DefaultLogin, string(hashedPassword))
		log.Printf("Создан пользователь по умолчанию: %s/%s", config.Auth.DefaultLogin, config.Auth.DefaultPassword)
	}
}

// Функции аутентификации
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	var user User
	var hashedPassword string
	err := db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = $1", 
		req.Login).Scan(&user.ID, &user.Login, &hashedPassword)

	if err != nil {
		response := AuthResponse{Success: false, Message: "Неверный логин или пароль"}
		json.NewEncoder(w).Encode(response)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		response := AuthResponse{Success: false, Message: "Неверный логин или пароль"}
		json.NewEncoder(w).Encode(response)
		return
	}

	// Простой токен (в продакшене использовать JWT)
	token := fmt.Sprintf("%s%d_%d", config.Auth.TokenPrefix, user.ID, time.Now().Unix())
	
	response := AuthResponse{Success: true, Message: "Успешная авторизация", Token: token}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CRUD операции для сотрудников
func getStaffHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, full_name, phone, email, address FROM staff ORDER BY id")
	if err != nil {
		http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var staff []Staff
	for rows.Next() {
		var s Staff
		if err := rows.Scan(&s.ID, &s.FullName, &s.Phone, &s.Email, &s.Address); err != nil {
			continue
		}
		staff = append(staff, s)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func createStaffHandler(w http.ResponseWriter, r *http.Request) {
	var staff Staff
	if err := json.NewDecoder(r.Body).Decode(&staff); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	err := db.QueryRow(
		"INSERT INTO staff (full_name, phone, email, address) VALUES ($1, $2, $3, $4) RETURNING id",
		staff.FullName, staff.Phone, staff.Email, staff.Address).Scan(&staff.ID)

	if err != nil {
		http.Error(w, "Ошибка создания сотрудника", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func updateStaffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var staff Staff
	if err := json.NewDecoder(r.Body).Decode(&staff); err != nil {
		http.Error(w, "Неверный формат данных", http.StatusBadRequest)
		return
	}

	_, err = db.Exec(
		"UPDATE staff SET full_name = $1, phone = $2, email = $3, address = $4 WHERE id = $5",
		staff.FullName, staff.Phone, staff.Email, staff.Address, id)

	if err != nil {
		http.Error(w, "Ошибка обновления сотрудника", http.StatusInternalServerError)
		return
	}

	staff.ID = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(staff)
}

func deleteStaffHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("DELETE FROM staff WHERE id = $1", id)
	if err != nil {
		http.Error(w, "Ошибка удаления сотрудника", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Сотрудник удален"})
}

// Middleware для проверки авторизации
func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			http.Error(w, "Токен авторизации не предоставлен", http.StatusUnauthorized)
			return
		}
		
		// Простая проверка токена (в продакшене использовать JWT)
		if len(token) < 10 {
			http.Error(w, "Неверный токен", http.StatusUnauthorized)
			return
		}
		
		next(w, r)
	}
}

func main() {
	// Загрузка конфигурации
	config = LoadConfig()
	
	// Инициализация базы данных
	if err := initDB(); err != nil {
		log.Fatal("Ошибка подключения к базе данных:", err)
	}
	defer db.Close()

	// Настройка маршрутов
	r := mux.NewRouter()

	// API маршруты
	api := r.PathPrefix("/api").Subrouter()
	
	// Авторизация
	api.HandleFunc("/login", loginHandler).Methods("POST")
	
	// CRUD для сотрудников
	api.HandleFunc("/staff", authMiddleware(getStaffHandler)).Methods("GET")
	api.HandleFunc("/staff", authMiddleware(createStaffHandler)).Methods("POST")
	api.HandleFunc("/staff/{id}", authMiddleware(updateStaffHandler)).Methods("PUT")
	api.HandleFunc("/staff/{id}", authMiddleware(deleteStaffHandler)).Methods("DELETE")

	// Статические файлы
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(config.Server.StaticDir)))

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{config.Server.CORSOrigin},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	handler := c.Handler(r)

	log.Printf("Сервер запущен на порту %s", config.Server.Port)
	log.Printf("Пользователь по умолчанию: %s/%s", config.Auth.DefaultLogin, config.Auth.DefaultPassword)
	log.Fatal(http.ListenAndServe(":"+config.Server.Port, handler))
}
