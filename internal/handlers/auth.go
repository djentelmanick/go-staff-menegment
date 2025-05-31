package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"staff-management/internal/models"
	"staff-management/internal/config"
	"golang.org/x/crypto/bcrypt"
)

func LoginHandler(db *sql.DB, cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.LoginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		var user models.User
		var hashedPassword string
		err := db.QueryRow("SELECT id, login, password_hash FROM users WHERE login = $1", 
			req.Login).Scan(&user.ID, &user.Login, &hashedPassword)

		if err != nil {
			response := models.AuthResponse{Success: false, Message: "Неверный логин или пароль"}
			json.NewEncoder(w).Encode(response)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
			response := models.AuthResponse{Success: false, Message: "Неверный логин или пароль"}
			json.NewEncoder(w).Encode(response)
			return
		}

		token := fmt.Sprintf("%s%d_%d", cfg.Auth.TokenPrefix, user.ID, time.Now().Unix())
		response := models.AuthResponse{Success: true, Message: "Успешная авторизация", Token: token}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}