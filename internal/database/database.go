package database

import (
	"database/sql"
	"log"

	"staff-management/internal/config"
	"golang.org/x/crypto/bcrypt"
	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB(cfg *config.Config) error {
	var err error
	db, err = sql.Open("postgres", cfg.GetDatabaseURL())
	if err != nil {
		return err
	}

	if err = db.Ping(); err != nil {
		return err
	}

	createTables()
	createDefaultAdmin(cfg)
	
	return nil
}

func createTables() {
	userTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		login VARCHAR(50) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL
	)`

	staffTable := `
	CREATE TABLE IF NOT EXISTS staff (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		phone VARCHAR(20),
		email VARCHAR(100),
		address TEXT
	)`

	staff_groupsTable := `
	CREATE TABLE IF NOT EXISTS staff_groups (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		description TEXT
	)`

	staff_to_groupsTable := `
	CREATE TABLE IF NOT EXISTS staff_to_groups (
		staff_id INT NOT NULL,
		group_id INT NOT NULL,
		PRIMARY KEY (staff_id, group_id),
		FOREIGN KEY (staff_id) REFERENCES staff(id) ON DELETE CASCADE,
		FOREIGN KEY (group_id) REFERENCES staff_groups(id) ON DELETE CASCADE
	)`

	add_fiels_status_to_staff := `ALTER TABLE staff ADD COLUMN IF NOT EXISTS status VARCHAR(50)`

	db.Exec(userTable)
	db.Exec(staffTable)
	db.Exec(staff_groupsTable)
	db.Exec(staff_to_groupsTable)
	db.Exec(add_fiels_status_to_staff)
}

func createDefaultAdmin(cfg *config.Config) {
	var count int
	db.QueryRow("SELECT COUNT(*) FROM users WHERE login = $1", cfg.Auth.DefaultLogin).Scan(&count)
	
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(cfg.Auth.DefaultPassword), bcrypt.DefaultCost)
		db.Exec("INSERT INTO users (login, password_hash) VALUES ($1, $2)", cfg.Auth.DefaultLogin, string(hashedPassword))
		log.Printf("Создан пользователь по умолчанию: %s/%s", cfg.Auth.DefaultLogin, cfg.Auth.DefaultPassword)
	}
}

func CloseDB() {
	db.Close()
}

func GetDB() *sql.DB {
	return db
}