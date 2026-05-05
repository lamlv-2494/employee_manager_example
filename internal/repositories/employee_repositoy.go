package repositories

import (
	"context"
	"database/sql"
	. "employee_manager_example/internal/models"
	. "employee_manager_example/internal/utils"
	"errors"
	"strings"
)

type EmployeeRepository interface {
	CreateEmployee(ctx context.Context, employee *Employee) error

	GetEmployees(ctx context.Context, limit, offet int) ([]Employee, int, error)

	GetEmployeeById(ctx context.Context, id int) (*Employee, error)

	UpdateEmployee(ctx context.Context, id int, req *UpdateEmployeeRequest) error
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

// GetEmployeeById implements [EmployeeRepository].
func (e *employeeRepository) GetEmployeeById(ctx context.Context, id int) (*Employee, error) {
	query := "SELECT id, name, age, position, department_id, salary FROM employees WHERE id = ? AND deleted_at IS NULL"
	row := e.db.QueryRowContext(ctx, query, id)

	var employee Employee
	err := row.Scan(&employee.Id, &employee.Name, &employee.Age, &employee.Position, &employee.DepartmentId, &employee.Salary)
	if err != nil {
		LogError(LogErrScan, err)
		return nil, err
	}

	return &employee, nil
}

// UpdateEmployee implements [EmployeeRepository].
func (e *employeeRepository) UpdateEmployee(ctx context.Context, id int, req *UpdateEmployeeRequest) error {
	if req == nil {
		return errors.New(ErrNoFieldUpdate)
	}

	var setClauses []string
	var args []any

	if req.Name != nil {
		setClauses = append(setClauses, "name = ?")
		args = append(args, *req.Name)
	}

	if req.Age != nil {
		setClauses = append(setClauses, "age = ?")
		args = append(args, *req.Age)
	}

	if req.Position != nil {
		setClauses = append(setClauses, "position = ?")
		args = append(args, *req.Position)
	}

	if req.DepartmentId != nil {
		setClauses = append(setClauses, "department_id = ?")
		args = append(args, *req.DepartmentId)
	}

	if req.Salary != nil {
		setClauses = append(setClauses, "salary = ?")
		args = append(args, *req.Salary)
	}

	setClause := strings.Join(setClauses, ", ") // "name = ?, age = ?"

	if setClause == "" {
		return errors.New(ErrNoFieldUpdate)
	}

	query := "UPDATE employees SET " + setClause + " WHERE id = ? AND deleted_at IS NULL"
	args = append(args, id)

	result, err := e.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	affect, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affect == 0 {
		return sql.ErrNoRows
	}

	return nil
}
