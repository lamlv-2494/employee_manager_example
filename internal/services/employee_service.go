package services

import (
	"context"
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/repositories"
	. "employee_manager_example/internal/utils"
	"errors"
	"strings"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, emp *Employee) error

	GetEmployees(ctx context.Context, page, limit int, keyword, position string) ([]Employee, int, error)

	GetEmployeeById(ctx context.Context, id int) (*Employee, error)

	UpdateEmployee(ctx context.Context, id int, req UpdateEmployeeRequest) error

	DeleteEmployee(ctx context.Context, id int) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// CreateEmployee implements [EmployeeService].
func (e *employeeService) CreateEmployee(ctx context.Context, emp *Employee) error {
	emp.Name = strings.TrimSpace(emp.Name)
	if emp.Name == "" {
		return errors.New(ErrEmployeeNameEmpty)
	}
	if emp.Age <= 0 {
		return errors.New(ErrEmployeeAgeInvalid)
	}
	if emp.Salary <= 0 {
		return errors.New(ErrEmployeeSalaryInvalid)
	}

	err := e.repo.CreateEmployee(ctx, emp)

	if err != nil {
		return err
	}
	return nil
}

// GetEmployees implements [EmployeeService].
func (e *employeeService) GetEmployees(ctx context.Context, page, limit int, keyword, position string) ([]Employee, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return e.repo.GetEmployees(ctx, limit, offset, keyword, position)
}

// GetEmployeeById implements [EmployeeService].
func (e *employeeService) GetEmployeeById(ctx context.Context, id int) (*Employee, error) {
	employee, err := e.repo.GetEmployeeById(ctx, id)

	if err != nil {
		return nil, err
	}

	if employee == nil {
		return nil, errors.New(ErrEmpty)
	}
	return employee, nil
}

// UpdateEmployee implements [EmployeeService].
func (e *employeeService) UpdateEmployee(ctx context.Context, id int, req UpdateEmployeeRequest) error {
	err := e.repo.UpdateEmployee(ctx, id, &req)
	if err != nil {
		return err
	}
	return nil
}

// DeleteEmployee implements [EmployeeService].
func (e *employeeService) DeleteEmployee(ctx context.Context, id int) error {
	err := e.repo.DeleteEmployee(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
