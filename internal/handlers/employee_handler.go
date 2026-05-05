package handlers

import (
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/services"
	. "employee_manager_example/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type EmployeeHandler struct {
	service services.EmployeeService
}

func NewEmployeeHandler(service services.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{service: service}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var emp Employee
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateEmployee(r.Context(), &emp)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusCreated, emp)
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	page, pageErr := strconv.Atoi(query.Get("page"))
	limit, limitErr := strconv.Atoi(query.Get("limit"))

	if pageErr != nil || limitErr != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrPageOrLimitInvalid)
		return
	}

	employees, totalCount, err := h.service.GetEmployees(r.Context(), page, limit)

	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"employees": employees,
		"total":     totalCount,
	}

	ResponseSuccess(w, http.StatusOK, response)
}

func (h *EmployeeHandler) GetEmployeeById(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponseWithError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	employeeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrInvalidEmployeeId)
		return
	}

	employee, err := h.service.GetEmployeeById(r.Context(), employeeId)

	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ResponseWithError(w, http.StatusMethodNotAllowed, ErrMethodNotAllowed)
		return
	}

	employeeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrInvalidEmployeeId)
		return
	}

	var req UpdateEmployeeRequest
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrInvalidRequestBody)
		return
	}

	err = h.service.UpdateEmployee(r.Context(), employeeId, req)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusNoContent, nil)
}
