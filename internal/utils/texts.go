package utils

const (
	ErrCanNotOpenDB    = "Can't open DB"
	ErrCanNotConnectDB = "Can't connect to DB"

	ErrServerError = "Server error"

	ErrMethodNotAllowed = "Method not allowed"

	ErrFailedToCreateEmployee = "Failed to create employee"
	ErrLastInsertId           = "Can't get last insert id"

	// Employee validation
	ErrEmployeeNameEmpty     = "Employee name must not be empty"
	ErrEmployeeAgeInvalid    = "Employee age must be greater than 0"
	ErrEmployeeSalaryInvalid = "Employee salary must be greater than 0"
	ErrAgeMustBePositive     = "Age must be >= 0"
	ErrSalaryMustBePositive  = "Salary must be >= 0"
)
