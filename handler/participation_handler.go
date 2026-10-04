package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"university-events/middleware"
	"university-events/models"
	"university-events/pkg/response"
	"university-events/repository"
)

type ParticipationHandler struct {
	participationRepo *repository.ParticipationRepository
	eventRepo         *repository.EventRepository
}

func NewParticipationHandler(participationRepo *repository.ParticipationRepository, eventRepo *repository.EventRepository) *ParticipationHandler {
	return &ParticipationHandler{participationRepo: participationRepo, eventRepo: eventRepo}
}


func (h *ParticipationHandler) Register(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	studentID, _ := r.Context().Value(middleware.ContextUserID).(string)

	var in models.RegisterInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fullName := strings.TrimSpace(in.FullName)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	studentNumber := normalizeDigits(strings.TrimSpace(in.StudentNumber))

	if fullName == "" || email == "" || studentNumber == "" {
		response.Error(w, http.StatusBadRequest, "full_name, email and student_number are required")
		return
	}

	if !validEmail(email) {
		response.Error(w, http.StatusBadRequest, "email is not valid")
		return
	}

	if !validStudentNumber(studentNumber) {
		response.Error(w, http.StatusBadRequest, "student_number must contain only digits (5 to 20 digits)")
		return
	}

	participation, err := h.participationRepo.Register(r.Context(), studentID, fullName, email, studentNumber, eventID)

	switch {
	case errors.Is(err, repository.ErrNotFound):
		response.Error(w, http.StatusNotFound, "event not found")
	case errors.Is(err, models.ErrEventFull):
		response.Error(w, http.StatusConflict, "this event has reached its capacity")
	case errors.Is(err, models.ErrAlreadyRegistered):
		response.Error(w, http.StatusConflict, "you are already registered for this event")
	case err != nil:
		log.Printf("register for event failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not register for event")
	default:
		response.JSON(w, http.StatusCreated, participation)
	}
}


func (h *ParticipationHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")
	studentID, _ := r.Context().Value(middleware.ContextUserID).(string)

	err := h.participationRepo.Cancel(r.Context(), eventID, studentID)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "you have no active registration for this event")
		return
	}
	if err != nil {
		log.Printf("cancel registration failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not cancel registration")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "registration cancelled"})
}

func (h *ParticipationHandler) MyRegistrations(w http.ResponseWriter, r *http.Request) {
	studentID, _ := r.Context().Value(middleware.ContextUserID).(string)

	regs, err := h.participationRepo.ListByStudent(r.Context(), studentID)
	if err != nil {
		log.Printf("list student registrations failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not load registrations")
		return
	}

	response.JSON(w, http.StatusOK, regs)
}


func (h *ParticipationHandler) ListParticipants(w http.ResponseWriter, r *http.Request) {
	eventID := chi.URLParam(r, "id")

	if _, err := h.eventRepo.GetByID(r.Context(), eventID); errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}

	participants, err := h.participationRepo.ListParticipants(r.Context(), eventID)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not load participants")
		return
	}

	response.JSON(w, http.StatusOK, participants)
}
