package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"employee_manager_example/internal/handlers"
	"employee_manager_example/internal/middleware"
	"employee_manager_example/internal/repositories"
	"employee_manager_example/internal/services"
	. "employee_manager_example/internal/utils"
	"log"
	"net/http"
	"time"
)

func main() {
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

	setupServerMux(appHandlers)
	db.Close()
}

func setupMySQL() *sql.DB {
	dsn := "root:@tcp(127.0.0.1:3306)/employee_management?parseTime=true"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("%v, error: %v", ErrCanNotOpenDB, err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("%v, error: %v", ErrCanNotConnectDB, err)
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
		log.Fatalf("%v, error: %v", ErrServerError, err)
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
			http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
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
			http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/employees/search", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEmployees(w, r)
		default:
			http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
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
			http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})

	mux.HandleFunc("/departments/{id}/employees", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handler.GetEmployeesByDepartmentId(w, r)
		default:
			http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
			return
		}
	})
}
