package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"fmt"
	"strings"

	"staff-management/internal/models"

	"github.com/gorilla/mux"
)

func GetStaffGroups(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query(`
			SELECT 
				sg.id,
				sg.name,
				sg.description,
				COUNT(stg.staff_id) AS member_count
			FROM 
				staff_groups sg
			LEFT JOIN 
				staff_to_groups stg ON sg.id = stg.group_id
			GROUP BY 
				sg.id, sg.name, sg.description
			ORDER BY 
				sg.id`)
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var staff_groups []models.StaffGroup
		for rows.Next() {
			var s models.StaffGroup
			if err := rows.Scan(&s.ID, &s.Name, &s.Description, &s.MemberCount); err != nil {
				continue
			}
			staff_groups = append(staff_groups, s)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff_groups)
	}
}

func CreateStaffGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var staff_group models.StaffGroup
		if err := json.NewDecoder(r.Body).Decode(&staff_group); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		err := db.QueryRow(
			"INSERT INTO staff_groups (name, description) VALUES ($1, $2) RETURNING id",
			staff_group.Name, staff_group.Description).Scan(&staff_group.ID)

		if err != nil {
			http.Error(w, "Ошибка создания группы", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff_group)
	}
}

func UpdateStaffGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var staff_group models.StaffGroup
		if err := json.NewDecoder(r.Body).Decode(&staff_group); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		_, err = db.Exec(
			"UPDATE staff_groups SET name = $1, description = $2 WHERE id = $3",
			staff_group.Name, staff_group.Description, id)
		
		if err != nil {
			http.Error(w, "Ошибка обновления группы", http.StatusInternalServerError)
			return
		}

		staff_group.ID = id
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff_group)
	}
}

func DeleteStaffGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM staff_groups WHERE id = $1", id)
		if err != nil {
			http.Error(w, "Ошибка удаления группы", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Группа удален"})
	}
}

func GetGroupMemebers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group_id, err := strconv.Atoi(vars["group_id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
			SELECT s.id, s.full_name
			FROM staff s
			JOIN staff_to_groups sg ON s.id = sg.staff_id
			WHERE sg.group_id = $1`, group_id)
		if err != nil {
			http.Error(w, "Ошибка получения данных", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var staff_short_info []models.StaffShortInfo
		for rows.Next() {
			var s models.StaffShortInfo
			if err := rows.Scan(&s.ID, &s.FullName); err != nil {
				continue
			}
			staff_short_info = append(staff_short_info, s)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(staff_short_info)
	}
}

func AddMemebersToGroup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group_id, err := strconv.Atoi(vars["group_id"])
		if err != nil {
			http.Error(w, "Неверный ID", http.StatusBadRequest)
			return
		}

		var staff_ids models.StaffIds
		if err := json.NewDecoder(r.Body).Decode(&staff_ids); err != nil {
			http.Error(w, "Неверный формат данных", http.StatusBadRequest)
			return
		}

		stmt, err := db.Prepare("INSERT INTO staff_to_groups (staff_id, group_id) VALUES ($1, $2)")
		if err != nil {
			http.Error(w, "Ошибка подготовки запроса", http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		for _, staff_id := range staff_ids.StaffIDs {
			_, err := stmt.Exec(staff_id, group_id)
			if err != nil {
				if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
					continue
				}
				http.Error(w, fmt.Sprintf("Ошибка добавления сотрудника %d в группу", staff_id), http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "Сотрудники успешно добавлены в группу"})
	}
}

func DeleteGroupMember(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		group_id, err := strconv.Atoi(vars["group_id"])
		if err != nil {
			http.Error(w, "Неверный ID группы", http.StatusBadRequest)
			return
		}

		staff_id, err := strconv.Atoi(vars["staff_id"])
		if err != nil {
			http.Error(w, "Неверный ID сотрудника", http.StatusBadRequest)
			return
		}

		result, err := db.Exec(
			"DELETE FROM staff_to_groups WHERE staff_id = $1 AND group_id = $2",
			staff_id, group_id)
		if err != nil {
			http.Error(w, "Ошибка удаления сотрудника из группы", http.StatusInternalServerError)
			return
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			http.Error(w, "Ошибка проверки результата удаления", http.StatusInternalServerError)
			return
		}

		if rowsAffected == 0 {
			http.Error(w, "Сотрудник не найден в группе", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
