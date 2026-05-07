package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"

	"employee_manager_example/internal/handlers"
	"employee_manager_example/internal/middleware"
	"employee_manager_example/internal/repositories"
	"employee_manager_example/internal/services"
	"employee_manager_example/internal/texts"
	"log"
	"net/http"
	"time"
)

func loadEnv(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		os.Setenv(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
	}
	return scanner.Err()
}

func main() {
	if err := loadEnv(".env"); err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	db := setupMySQL()
	employeeRepo := repositories.NewEmployeeRepository(db)
	employeeService := services.NewEmployeeService(employeeRepo)
	employeeHander := handlers.NewEmployeeHandler(employeeService)

	departmentRepo := repositories.NewDepartmentRepository(db)
	departmentService := services.NewDepartmentService(departmentRepo)
	departmentHandler := handlers.NewDepartmentHandler(departmentService)

	appHandlers := &handlers.AppHandler{
		EmployeeHandler:   employeeHander,
		DepartmentHandler: departmentHandler,
	}

	err := employeeService.ExportData(context.Background())
	if err != nil {
		fmt.Println("Export failure:", err)
	}

	defer db.Close()
	setupServerMux(appHandlers)
}

func setupMySQL() *sql.DB {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbName == "" {
		log.Fatal("Missing required database environment variables")
	}

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("%v, error: %v", texts.CanNotOpenDB, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("%v, error: %v", texts.CanNotConnectDB, err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	return db
}

func setupServerMux(appHandler *handlers.AppHandler) {
	mux := http.NewServeMux()
	registerEmployeeRoutes(mux, appHandler.EmployeeHandler)
	registerDepartmentRoutes(mux, appHandler.DepartmentHandler)

	loggedMux := middleware.LoggingMiddleware(mux)
	err := http.ListenAndServe(":8080", loggedMux)
	if err != nil {
		log.Fatalf("%v, error: %v", texts.ServerError, err)
	}
}

func registerEmployeeRoutes(mux *http.ServeMux, handler *handlers.EmployeeHandler) {
	mux.HandleFunc("/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.CreateEmployee(w, r)
		case http.MethodGet:
			handler.GetEmployees(w, r)
		default:
			http.Error(w, texts.MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/employees/search", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEmployees(w, r)
		default:
			http.Error(w, texts.MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/employees/{id}", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEmployeeById(w, r)
		case http.MethodPut:
			handler.UpdateEmployee(w, r)
		case http.MethodDelete:
			handler.DeleteEmployee(w, r)
		default:
			http.Error(w, texts.MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})
}

func registerDepartmentRoutes(mux *http.ServeMux, handler *handlers.DepartmentHandler) {
	mux.HandleFunc("/departments", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetDepartments(w, r)
		case http.MethodPost:
			handler.CreateDepartment(w, r)
		default:
			http.Error(w, texts.MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/departments/{id}/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEmployeesByDepartmentId(w, r)
		default:
			http.Error(w, texts.MethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})
}
