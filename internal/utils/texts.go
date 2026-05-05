package utils

const (
	ErrCanNotOpenDB    = "Can't open DB"
	ErrCanNotConnectDB = "Can't connect to DB"

	ErrServerError = "Server error"

	ErrMethodNotAllowed = "Method not allowed"

	ErrFailedToCreateEmployee = "Failed to create employee"
	ErrLastInsertId           = "Can't get last insert id"

	// Invalid error
	ErrPageOrLimitInvalid = "Invalid page or limit parameter"
	ErrInvalidEmployeeId  = "Invalid employee ID"
	ErrInvalidRequestBody = "Invalid request body"

	// Employee validation

	ErrEmployeeNameEmpty     = "Employee name must not be empty"
	ErrEmployeeAgeInvalid    = "Employee age must be greater than 0"
	ErrEmployeeSalaryInvalid = "Employee salary must be greater than 0"
	ErrAgeMustBePositive     = "Age must be >= 0"
	ErrSalaryMustBePositive  = "Salary must be >= 0"

	// Log messages
	LogErrPrepareStatement = "Prepare statement error: %v\n"
	LogErrExecuteStatement = "Execute statement error: %v\n"
	LogErrScan             = "Scan error:"
	LogErrQuery            = "Query error"

	// Log Empty
	ErrEmpty         = "No Employee found"
	ErrNoFieldUpdate = "No fields to update"
)
