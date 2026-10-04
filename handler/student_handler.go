package handler

import (
	"errors"
	"net/http"

	"university-events/middleware"
	"university-events/pkg/response"
	"university-events/repository"
)

type StudentHandler struct {
	studentRepo *repository.StudentRepository
}

func NewStudentHandler(studentRepo *repository.StudentRepository) *StudentHandler {
	return &StudentHandler{studentRepo: studentRepo}
}

// Me returns the currently authenticated student's public profile.
func (h *StudentHandler) Me(w http.ResponseWriter, r *http.Request) {
	studentID, _ := r.Context().Value(middleware.ContextUserID).(string)

	student, err := h.studentRepo.GetByID(r.Context(), studentID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "student not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not load student")
		return
	}

	response.JSON(w, http.StatusOK, student.ToPublic())
}
