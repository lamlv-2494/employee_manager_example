package utils

import (
	"employee_manager_example/internal/models"
	"employee_manager_example/internal/texts"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
)

func ExportToJson(data []models.Employee) error {
	filename := "exports/employees.json"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(data); err != nil {
		return errors.New(texts.ErrEncodeJSON + err.Error())
	}
	return nil
}

func ExportToCSV(data []models.Employee) error {
	filename := "exports/employees.csv"
	file, err := os.Create(filename)
	if err != nil {
		return err
	}

	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	if err := writer.Write([]string{"ID", "Name", "Age", "Position", "DepartmentID", "Salary"}); err != nil {
		return errors.New(texts.ErrWriteCSVHeader + err.Error())
	}

	// Write employee data
	for _, emp := range data {
		record := []string{
			strconv.Itoa(emp.Id),
			emp.Name,
			strconv.Itoa(emp.Age),
			emp.Position,
			strconv.Itoa(emp.DepartmentId),
			fmt.Sprintf("%.2f", emp.Salary),
		}
		if err := writer.Write(record); err != nil {
			return errors.New(texts.ErrWriteCSVRecord + err.Error())
		}
	}
	return nil
}
