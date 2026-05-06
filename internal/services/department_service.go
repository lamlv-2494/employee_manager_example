package services

import (
	"context"
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/repositories"
	"employee_manager_example/internal/utils"
	"errors"
	"strings"
)

type DepartmentService interface {
	CreateDepartment(ctx context.Context, department *Department) error

	GetDepartments(ctx context.Context, page, limit int) ([]Department, int, error)

	GetEmployeesByDepartmentId(ctx context.Context, departmentId, page, limit int) ([]Employee, int, error)
}

type departmentService struct {
	repro repositories.DepartmentRepository
}

func NewDepartmentService(repo repositories.DepartmentRepository) DepartmentService {
	return &departmentService{repro: repo}
}

// CreateDepartment implements [DepartmentService].
func (d *departmentService) CreateDepartment(ctx context.Context, department *Department) error {
	department.Name = strings.TrimSpace(department.Name)
	if department.Name == "" {
		return errors.New(utils.ErrDepartmentNameEmpty)
	}
	err := d.repro.CreateDepartment(ctx, department)
	if err != nil {
		return err
	}
	return nil
}

// GetDepartments implements [DepartmentService].
func (d *departmentService) GetDepartments(ctx context.Context, page int, limit int) ([]Department, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	departments, totalCount, err := d.repro.GetDepartments(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return departments, totalCount, nil
}

// GetEmployeesByDepartmentId implements [DepartmentService].
func (d *departmentService) GetEmployeesByDepartmentId(ctx context.Context, departmentId, page, limit int) ([]Employee, int, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	offset := (page - 1) * limit

	employees, totalCount, err := d.repro.GetEmployeesByDepartmentId(ctx, departmentId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	return employees, totalCount, nil
}
