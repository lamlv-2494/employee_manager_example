package handlers

import (
	. "employee_manager_example/internal/models"
	"employee_manager_example/internal/services"
	. "employee_manager_example/internal/utils"
	"employee_manager_example/internal/texts"
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
		ResponseError(w, http.StatusMethodNotAllowed, texts.MethodNotAllowed)
		return
	}

	var emp Employee
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateEmployee(r.Context(), &emp)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusCreated, emp)
}

func (h *EmployeeHandler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		ResponseError(w, http.StatusMethodNotAllowed, texts.MethodNotAllowed)
		return
	}

	query := r.URL.Query()
	page, pageErr := strconv.Atoi(query.Get("page"))
	limit, limitErr := strconv.Atoi(query.Get("limit"))

	if pageErr != nil || limitErr != nil {
		ResponseError(w, http.StatusBadRequest, texts.PageOrLimitInvalid)
		return
	}

	// Search
	keyword := query.Get("keyword")
	position := query.Get("position")

	employees, totalCount, err := h.service.GetEmployees(r.Context(), page, limit, keyword, position)

	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
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
		ResponseError(w, http.StatusMethodNotAllowed, texts.MethodNotAllowed)
		return
	}

	employeeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, texts.InvalidEmployeeId)
		return
	}

	employee, err := h.service.GetEmployeeById(r.Context(), employeeId)

	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusOK, employee)
}

func (h *EmployeeHandler) UpdateEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		ResponseError(w, http.StatusMethodNotAllowed, texts.MethodNotAllowed)
		return
	}

	employeeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, texts.InvalidEmployeeId)
		return
	}

	var req UpdateEmployeeRequest
	defer r.Body.Close()
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		ResponseError(w, http.StatusBadRequest, texts.InvalidRequestBody)
		return
	}

	err = h.service.UpdateEmployee(r.Context(), employeeId, req)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusNoContent, nil)
}

func (h *EmployeeHandler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		ResponseError(w, http.StatusMethodNotAllowed, texts.MethodNotAllowed)
		return
	}
	employeeId, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		ResponseError(w, http.StatusBadRequest, texts.InvalidEmployeeId)
		return
	}
	err = h.service.DeleteEmployee(r.Context(), employeeId)
	if err != nil {
		ResponseError(w, http.StatusInternalServerError, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusNoContent, nil)
}
