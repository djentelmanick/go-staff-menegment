package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"staff-management/internal/models"

	"github.com/gorilla/mux"
)

func GetStaffStatuses(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`SELECT id, full_name, status FROM staff`)
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var staff_statuses []models.EmployeeStatus
		for rows.Next() {
			var s models.EmployeeStatus
			var bd_status sql.NullString
			if err := rows.Scan(&s.ID, &s.FullName, &bd_status); err != nil {
				http.Error(w, "Ошибка обработки данных", http.StatusInternalServerError)
				continue
			}
			s.Status = bd_status.String
			staff_statuses = append(staff_statuses, s)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff_statuses)
	}
}

func UpdateStaffStatus(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var status models.Status
		if err := json.NewDecoder(r.Body).Decode(&status); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"UPDATE staff SET status = $1 WHERE id = $2",
			status.Status, id)
		
		if err != nil {
			http.Error(w, "Ошибка обновления статуса", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}	
}
