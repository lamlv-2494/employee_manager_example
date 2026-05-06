package handlers

type AppHandler struct {
	EmployeeHandler   *EmployeeHandler
	DepartmentHandler *DepartmentHandler
}
