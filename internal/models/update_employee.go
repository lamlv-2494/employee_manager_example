package models

type UpdateEmployeeRequest struct {
	Name         *string  `json:"name,omitempty"`
	Age          *int     `json:"age,omitempty"`
	Position     *string  `json:"position,omitempty"`
	DepartmentId *int     `json:"department_id,omitempty"`
	Salary       *float64 `json:"salary,omitempty"`
}
