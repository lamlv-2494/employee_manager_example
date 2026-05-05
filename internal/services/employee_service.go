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
	CreateEmployee(ctx context.Context, emp Employee) error
}

type employeeService struct {
	repo repositories.EmployeeRepository
}

func NewEmployeeService(repo repositories.EmployeeRepository) EmployeeService {
	return &employeeService{repo: repo}
}

// CreateEmployee implements [EmployeeService].
func (e *employeeService) CreateEmployee(ctx context.Context, emp Employee) error {
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

	err := e.repo.CreateEmployee(ctx, &emp)

	if err != nil {
		return err
	}
	return nil
}
