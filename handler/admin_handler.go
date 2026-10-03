package handler

import (
	"errors"
	"net/http"

	"university-events/middleware"
	"university-events/pkg/response"
	"university-events/repository"
)

type AdminHandler struct {
	adminRepo *repository.AdminRepository
}

func NewAdminHandler(adminRepo *repository.AdminRepository) *AdminHandler {
	return &AdminHandler{adminRepo: adminRepo}
}

// Me returns the currently authenticated admin's public profile.
func (h *AdminHandler) Me(w http.ResponseWriter, r *http.Request) {
	adminID, _ := r.Context().Value(middleware.ContextUserID).(string)

	admin, err := h.adminRepo.GetByID(r.Context(), adminID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "admin not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not load admin")
		return
	}

	response.JSON(w, http.StatusOK, admin.ToPublic())
}
