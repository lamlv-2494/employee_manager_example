package handlers

import (
	"employee_manager_example/internal/models"
	"employee_manager_example/internal/services"
	. "employee_manager_example/internal/utils"
	"encoding/json"
	"net/http"
	"strconv"
)

type DepartmentHandler struct {
	service services.DepartmentService
}

func NewDepartmentHandler(service services.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{service: service}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	var department models.Department
	defer r.Body.Close()

	err := json.NewDecoder(r.Body).Decode(&department)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = h.service.CreateDepartment(r.Context(), &department)
	if err != nil {
		ResponseWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	ResponseSuccess(w, http.StatusCreated, department)
}

func (h *DepartmentHandler) GetDepartments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	page, pageErr := strconv.Atoi(query.Get("page"))
	limit, limitErr := strconv.Atoi(query.Get("limit"))
	if pageErr != nil || limitErr != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrParseErr)

	}

	departments, totalCount, err := h.service.GetDepartments(r.Context(), page, limit)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"departments": departments,
		"total":       totalCount,
	}

	ResponseSuccess(w, http.StatusOK, response)
}

func (h *DepartmentHandler) GetEmployeesByDepartmentId(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, ErrMethodNotAllowed, http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query()
	departmentId, departmentErr := strconv.Atoi(r.PathValue("id"))
	page, pageErr := strconv.Atoi(query.Get("page"))
	limit, limitErr := strconv.Atoi(query.Get("limit"))
	if departmentErr != nil || pageErr != nil || limitErr != nil {
		ResponseWithError(w, http.StatusBadRequest, ErrParseErr)
		return
	}

	employees, totalCount, err := h.service.GetEmployeesByDepartmentId(r.Context(), departmentId, page, limit)
	if err != nil {
		ResponseWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]any{
		"employees":   employees,
		"total_count": totalCount,
	}

	ResponseSuccess(w, http.StatusOK, response)

}
