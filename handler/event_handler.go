package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"university-events/middleware"
	"university-events/models"
	"university-events/pkg/response"
	"university-events/repository"
)

type EventHandler struct {
	eventRepo *repository.EventRepository
}

func NewEventHandler(eventRepo *repository.EventRepository) *EventHandler {
	return &EventHandler{eventRepo: eventRepo}
}

// List returns all events ordered by event_time (soonest first).
// Fully public — no account or token needed to view it.
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	events, err := h.eventRepo.List(r.Context())
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not load events")
		return
	}
	response.JSON(w, http.StatusOK, events)
}

func (h *EventHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	event, err := h.eventRepo.GetByID(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not load event")
		return
	}
	response.JSON(w, http.StatusOK, event)
}

// Create adds a new event. Restricted to the admin (RequireRole middleware).
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var in models.CreateEventInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if in.Title == "" || in.Location == "" || in.Capacity <= 0 || in.EventTime.IsZero() {
		response.Error(w, http.StatusBadRequest, "title, location, event_time and a positive capacity are required")
		return
	}

	adminID, _ := r.Context().Value(middleware.ContextUserID).(string)

	event := &models.Event{
		Title:       in.Title,
		Description: in.Description,
		EventTime:   in.EventTime,
		Location:    in.Location,
		Capacity:    in.Capacity,
		CreatedBy:   adminID,
	}

	if err := h.eventRepo.Create(r.Context(), event); err != nil {
		response.Error(w, http.StatusInternalServerError, "could not create event")
		return
	}

	response.JSON(w, http.StatusCreated, event)
}

// Update edits an existing event. Restricted to the admin.
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var in models.UpdateEventInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if in.Title == "" || in.Location == "" || in.Capacity <= 0 || in.EventTime.IsZero() {
		response.Error(w, http.StatusBadRequest, "title, location, event_time and a positive capacity are required")
		return
	}

	event, err := h.eventRepo.Update(r.Context(), id, in)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}
	if errors.Is(err, repository.ErrCapacityBelowRegistered) {
		response.Error(w, http.StatusBadRequest, "capacity cannot be lower than the number of already registered participants")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not update event")
		return
	}

	response.JSON(w, http.StatusOK, event)
}

// Delete removes an event. Restricted to the admin.
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.eventRepo.Delete(r.Context(), id)
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusNotFound, "event not found")
		return
	}
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not delete event")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "event deleted"})
}
