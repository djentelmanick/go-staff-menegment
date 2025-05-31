package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"

	"staff-management/internal/models"
)

func GetStaffHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, full_name, phone, email, address FROM staff ORDER BY id")
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var staff []models.Staff
		for rows.Next() {
			var s models.Staff
			if err := rows.Scan(&s.ID, &s.FullName, &s.Phone, &s.Email, &s.Address); err != nil {
				continue
			}
			staff = append(staff, s)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff)
	}
}

func CreateStaffHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var staff models.Staff
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
}

func UpdateStaffHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var staff models.Staff
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
}

func DeleteStaffHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}