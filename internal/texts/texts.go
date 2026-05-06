package texts

const (
	ServerError      = "Server error"
	MethodNotAllowed = "Method not allowed"
	NoFieldUpdate    = "No fields to update"
	ParseErr         = "Parse error"

	// Invalid error
	PageOrLimitInvalid = "Invalid page or limit parameter"
	InvalidEmployeeId  = "Invalid employee ID"
	InvalidRequestBody = "Invalid request body"

	// Employee
	EmployeeNameEmpty      = "Employee name must not be empty"
	EmployeeAgeInvalid     = "Employee age must be greater than 0"
	EmployeeSalaryInvalid  = "Employee salary must be greater than 0"
	AgeMustBePositive      = "Age must be >= 0"
	SalaryMustBePositive   = "Salary must be >= 0"
	EmployeeNotFound       = "No Employee found"
	FailedToCreateEmployee = "Failed to create employee"

	// Department
	DepartmentNameEmpty = "Department name must not be empty"
	InvalidDepartmentId = "Invalid department ID"

	// DB messages
	CanNotOpenDB    = "Can't open DB"
	CanNotConnectDB = "Can't connect to DB"
	ExecuteContext  = "Execute ctx error"
	ScanErr         = "Scan error"
	QueryErr        = "Query error"
	QueryRowErr     = "Query row error"
	RowsAffected    = "Rows affected error"
	LastInsertId    = "Can't get last insert id"
)
