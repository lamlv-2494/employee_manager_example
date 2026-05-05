package repositories

import (
	"context"
	"database/sql"
	. "employee_manager_example/internal/models"
	. "employee_manager_example/internal/utils"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *Employee) error

	GetEmployees(ctx context.Context, limit, offet int) ([]Employee, int, error)
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

// GetEmployees implements [EmployeeRepository].
func (e *employeeRepository) GetEmployees(ctx context.Context, limit int, offet int) ([]Employee, int, error) {
	// Count total employees in DB
	var totalCount int

	countQuery := "SELECT COUNT(id) FROM employees WHERE deleted_at IS NULL"
	err := e.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)

	if err != nil {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []Employee{}, 0, nil
	}

	// Get employees from DB
	query := "SELECT id, name, age, position, department_id, salary FROM employees WHERE deleted_at IS NULL LIMIT ? OFFSET ?"
	rows, err := e.db.QueryContext(ctx, query, limit, offet)
	if LogError(LogErrQuery, err) {
		return nil, 0, err
	}
	defer rows.Close()

	var employees []Employee

	for rows.Next() {
		var employee Employee
		err := rows.Scan(&employee.Id, &employee.Name, &employee.Age, &employee.Position, &employee.DepartmentId, &employee.Salary)
		if LogError(LogErrScan, err) {
			return nil, 0, err
		}
		employees = append(employees, employee)
	}

	return employees, totalCount, nil
}
