package services

import (
	"context"
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/repositories"
	"employee_manager_example/internal/texts"
	"employee_manager_example/internal/utils"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

type EmployeeService interface {
	CreateEmployee(ctx context.Context, emp *Employee) error

	GetEmployees(ctx context.Context, page, limit int, keyword, position string) ([]Employee, int, error)

	GetEmployeeById(ctx context.Context, id int) (*Employee, error)

	UpdateEmployee(ctx context.Context, id int, req UpdateEmployeeRequest) error

	DeleteEmployee(ctx context.Context, id int) error

	ExportData(ctx context.Context) error
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
		return errors.New(texts.EmployeeNameEmpty)
	}
	if emp.Age <= 0 {
		return errors.New(texts.EmployeeAgeInvalid)
	}
	if emp.Salary <= 0 {
		return errors.New(texts.EmployeeSalaryInvalid)
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
		return nil, errors.New(texts.EmployeeNotFound)
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

// ExecuteConcurrentExport implements [EmployeeService].
func (e *employeeService) ExportData(ctx context.Context) error {
	fmt.Println("⏳ [1/3] Getting data from MySQL...")
	startTime := time.Now()

	employees, totalCount, err := e.repo.GetEmployees(ctx, 0, 0, "", "")
	if err != nil {
		log.Fatalf("Error when get data from SQL: %v", err)
	}

	fmt.Printf("✅ Fetched %d employees from DB. Time taken: %v\n", totalCount, time.Since(startTime))

	if totalCount == 0 {
		fmt.Println("No data to export!")
		return errors.New("No data to export!")
	}

	var wg sync.WaitGroup
	var mutex sync.Mutex

	totalFilesCreated := 0
	totalRecordsExported := 0

	fmt.Println("⏳ [2/3] start writing JSON and CSV files parallel...")
	exportStartTime := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := utils.ExportToJson(employees)
		if err != nil {
			log.Printf("Error while export data to JSON: %v", err)
			return
		}

		mutex.Lock()
		totalFilesCreated++
		totalRecordsExported += len(employees)
		mutex.Unlock()
		fmt.Println("   -> [JSON] Done file employees.json")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := utils.ExportToCSV(employees)
		if err != nil {
			log.Printf("Error while export data to CSV: %v", err)
			return
		}

		mutex.Lock()
		totalFilesCreated++
		totalRecordsExported += len(employees)
		mutex.Unlock()
		fmt.Println("   -> [CSV] Done file employees.csv")
	}()

	wg.Wait()

	fmt.Println("✅ [3/3] Export Done!")
	fmt.Println("=====================================")
	fmt.Printf("[FILE] File created successfully  : %d files\n", totalFilesCreated)
	fmt.Printf("[DATA] Total records exported     : %d records\n", totalRecordsExported)
	fmt.Printf("[TIME] Writing file times         : %v\n", time.Since(exportStartTime))
	fmt.Printf("[DONE] Total times                : %v\n", time.Since(startTime))
	fmt.Println("=====================================")
	return nil
}
