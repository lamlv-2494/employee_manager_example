package repositories

import (
	"context"
	"database/sql"
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/texts"
	. "employee_manager_example/internal/utils"
)

type DepartmentRepository interface {
	CreateDepartment(ctx context.Context, department *Department) error

	GetDepartments(ctx context.Context, limit, offset int) ([]Department, int, error)

	GetEmployeesByDepartmentId(ctx context.Context, departmentId, limit, offset int) ([]Employee, int, error)
}

type departmentRepository struct {
	db *sql.DB
}

func NewDepartmentRepository(db *sql.DB) DepartmentRepository {
	return &departmentRepository{db: db}
}

// CreateDepartment implements [DepartmentRepository].
func (d *departmentRepository) CreateDepartment(ctx context.Context, department *Department) error {
	query := "INSERT INTO departments (name) VALUES (?)123"

	res, err := d.db.ExecContext(ctx, query, department.Name)
	if LogError(texts.ExecuteContext, err) {
		return err
	}

	id, err := res.LastInsertId()
	if LogError(texts.LastInsertId, err) {
		return err
	}

	department.ID = int(id)
	return nil
}

// GetDepartments implements [DepartmentRepository].
func (d *departmentRepository) GetDepartments(ctx context.Context, limit int, offset int) ([]Department, int, error) {
	countQuery := "SELECT COUNT(id) FROM departments WHERE deleted_at IS NULL"

	var totalCount int
	err := d.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if LogError(texts.QueryRowErr, err) {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []Department{}, 0, nil
	}

	query := "SELECT id, name FROM departments WHERE deleted_at IS NULL LIMIT ? OFFSET ?"
	rows, err := d.db.QueryContext(ctx, query, limit, offset)
	if LogError(texts.QueryErr, err) {
		return nil, 0, err
	}
	defer rows.Close()

	var departments []Department

	for rows.Next() {
		var department Department
		err := rows.Scan(&department.ID, &department.Name)
		if LogError(texts.ScanErr, err) {
			return nil, 0, err
		}
		departments = append(departments, department)
	}

	return departments, totalCount, nil
}

// GetEmployeesByDepartmentId implements [DepartmentRepository].
func (d *departmentRepository) GetEmployeesByDepartmentId(ctx context.Context, departmentId, limit, offset int) ([]Employee, int, error) {
	countQuery := "SELECT COUNT(*) FROM employees WHERE department_id = ? AND deleted_at IS NULL"
	var totalCount int
	err := d.db.QueryRowContext(ctx, countQuery, departmentId).Scan(&totalCount)
	if LogError(texts.QueryRowErr, err) {
		return nil, 0, err
	}

	if totalCount == 0 {
		return []Employee{}, 0, nil
	}

	dataQuery := "SELECT e.id, e.name, e.age, e.position, e.department_id, e.salary, d.name AS department_name " +
		"FROM employees e " +
		"JOIN departments d ON e.department_id = d.id " +
		"WHERE e.department_id = ? AND e.deleted_at IS NULL AND d.deleted_at IS NULL " +
		"LIMIT ? OFFSET ?"

	rows, err := d.db.QueryContext(ctx, dataQuery, departmentId, limit, offset)
	if LogError(texts.QueryErr, err) {
		return nil, 0, err
	}

	var employees []Employee
	defer rows.Close()

	for rows.Next() {
		var emp Employee
		err := rows.Scan(&emp.Id, &emp.Name, &emp.Age, &emp.Position, &emp.DepartmentId, &emp.Salary, &emp.DepartmentName)
		if LogError(texts.ScanErr, err) {
			return nil, 0, err
		}
		employees = append(employees, emp)
	}
	return employees, totalCount, nil
}
