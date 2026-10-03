package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"university-events/models"
	"university-events/pkg/jwt"
	"university-events/pkg/response"
	"university-events/repository"
)

type AuthHandler struct {
	adminRepo   *repository.AdminRepository
	studentRepo *repository.StudentRepository
	jwtManager  *jwt.Manager
}

func NewAuthHandler(adminRepo *repository.AdminRepository, studentRepo *repository.StudentRepository, jwtManager *jwt.Manager) *AuthHandler {
	return &AuthHandler{adminRepo: adminRepo, studentRepo: studentRepo, jwtManager: jwtManager}
}


type authResponse struct {
	Token string      `json:"token"`
	Role  string      `json:"role"`
	User  interface{} `json:"user"`
}

type loginRequest struct {
	// Username is the admin's username, or a student's email.
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	minPasswordLen = 8
	maxPasswordLen = 72 
)


func (h *AuthHandler) gitgitLogin(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	username := strings.TrimSpace(req.Username)
	if username == "" || req.Password == "" {
		response.Error(w, http.StatusBadRequest, "username and password are required")
		return
	}

	// 1) Admin?
	admin, err := h.adminRepo.GetByEmail(r.Context(), username)
	if err == nil {
		if bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)) != nil {
			response.Error(w, http.StatusUnauthorized, "invalid username or password")
			return
		}
		h.respondWithToken(w, http.StatusOK, admin.ID, jwt.RoleAdmin, admin.ToPublic())
		return
	}
	if !errors.Is(err, repository.ErrNotFound) {
		log.Printf("login (admin lookup) failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not verify credentials")
		return
	}

	// 2) Student (logs in with email)?
	student, err := h.studentRepo.GetByEmail(r.Context(), strings.ToLower(username))
	if errors.Is(err, repository.ErrNotFound) {
		response.Error(w, http.StatusUnauthorized, "invalid username or password")
		return
	}
	if err != nil {
		log.Printf("login (student lookup) failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not verify credentials")
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(student.PasswordHash), []byte(req.Password)) != nil {
		response.Error(w, http.StatusUnauthorized, "invalid username or password")
		return
	}

	h.respondWithToken(w, http.StatusOK, student.ID, jwt.RoleStudent, student.ToPublic())
}

func (h *AuthHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var in models.SignupInput
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	fullName := strings.TrimSpace(in.FullName)
	email := strings.ToLower(strings.TrimSpace(in.Email))
	studentNumber := normalizeDigits(strings.TrimSpace(in.StudentNumber))

	if fullName == "" || email == "" || studentNumber == "" || in.Password == "" {
		response.Error(w, http.StatusBadRequest, "full_name, email, student_number and password are required")
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
	if len(in.Password) < minPasswordLen || len(in.Password) > maxPasswordLen {
		response.Error(w, http.StatusBadRequest, "password must be between 8 and 72 characters")
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(in.Password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not create account")
		return
	}

	student := &models.Student{
		FullName:      fullName,
		Email:         email,
		StudentNumber: studentNumber,
		PasswordHash:  string(hash),
	}

	err = h.studentRepo.Create(r.Context(), student)
	switch {
	case errors.Is(err, repository.ErrEmailTaken):
		response.Error(w, http.StatusConflict, "an account with this email already exists")
		return
	case errors.Is(err, repository.ErrStudentNumberTaken):
		response.Error(w, http.StatusConflict, "an account with this student number already exists")
		return
	case err != nil:
		log.Printf("signup failed: %v", err)
		response.Error(w, http.StatusInternalServerError, "could not create account")
		return
	}

	h.respondWithToken(w, http.StatusCreated, student.ID, jwt.RoleStudent, student.ToPublic())
}

func (h *AuthHandler) respondWithToken(w http.ResponseWriter, status int, userID, role string, user interface{}) {
	token, err := h.jwtManager.Generate(userID, role)
	if err != nil {
		response.Error(w, http.StatusInternalServerError, "could not generate token")
		return
	}
	response.JSON(w, status, authResponse{Token: token, Role: role, User: user})
}


func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{"message": "logged out"})
}
