package models

type Employee struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Age          int     `json:"age"`
	Position     string  `json:"position"`
	DepartmentId int     `json:"department_id"`
	Salary       float64 `json:"salary"`
}
