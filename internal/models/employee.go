package models

type Employee struct {
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Age            int     `json:"age"`
	Position       string  `json:"position"`
	Salary         float64 `json:"salary"`
	DepartmentId   int     `json:"department_id"`
	DepartmentName string  `json:"department_name,omitempty"`
}
