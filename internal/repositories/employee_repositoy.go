package repositories

import (
	"context"
	"database/sql"
	. "employee_manager_example/internal/models"
	. "employee_manager_example/internal/utils"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *Employee) error
}

type employeeRepository struct {
	db *sql.DB
}

func NewEmployeeRepository(db *sql.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

// CreateEmployee implements [EmployeeRepository].
func (e *employeeRepository) CreateEmployee(ctx context.Context, employee *Employee) error {
	// Insert new employee to DB
	query := "INSERT INTO employees (name, age, position, department_id, salary) VALUES (?, ?, ?, ?, ?)"
	res, err := e.db.ExecContext(ctx, query, employee.Name, employee.Age, employee.Position, employee.DepartmentId, employee.Salary)
	if err != nil {
		return err
	}

	// Get latest employee id
	id, err := res.LastInsertId()
	if LogError(ErrLastInsertId, err) {
		return err
	}

	employee.Id = int(id)

	return nil
}
