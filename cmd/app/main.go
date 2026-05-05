package main

import (
	"database/sql"
	"employee-system/Development/Projects/Research/Backend/employee_manager_example/internal/middleware"
	"log"
	"net/http"
	"time"
)

func main() {
	setupMySQL()
	setupServerMux()
}

func setupMySQL() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/employee_db?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Can't open DB, error: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Can't connect to DB, error: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func setupServerMux() {
	mux := http.NewServeMux()
	loggedMux := middleware.LoggingMiddleware(mux)
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("Server error: %V", err)
	}
}
