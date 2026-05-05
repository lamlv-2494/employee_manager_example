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
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateEmployee(r.Context(), &emp)
	if err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
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
		RespondWithError(w, http.StatusBadRequest, ErrPageOrLimitInvalid)
		return
	}

	employees, totalCount, err := h.service.GetEmployees(r.Context(), page, limit)

	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"employees": employees,
		"total":     totalCount,
	}

	ResponseSuccess(w, http.StatusOK, response)
}
