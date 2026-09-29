package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"

	"university-events/database"
	"university-events/handler"
	"university-events/internal/config"
	appmiddleware "university-events/middleware"
	"university-events/pkg/jwt"
	"university-events/repository"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(database.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("connected to database")

	// Repositories
	adminRepo := repository.NewAdminRepository(db)
	eventRepo := repository.NewEventRepository(db)
	studentRepo := repository.NewStudentRepository(db)
	participationRepo := repository.NewParticipationRepository(db)

	// Auth
	jwtManager := jwt.NewManager(cfg.JWTSecret, cfg.JWTExpiryHours)

	// Handlers
	authHandler := handler.NewAuthHandler(adminRepo, studentRepo, jwtManager)
	adminHandler := handler.NewAdminHandler(adminRepo)
	studentHandler := handler.NewStudentHandler(studentRepo)
	eventHandler := handler.NewEventHandler(eventRepo)
	participationHandler := handler.NewParticipationHandler(participationRepo, eventRepo)

	r := chi.NewRouter()
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)
	r.Use(appmiddleware.CORS())

	r.Route("/api/v1", func(r chi.Router) {
		// Public — no token needed. Anyone can browse events, log in, or
		// sign up as a student.
		r.Post("/auth/login", authHandler.Login)
		r.Post("/auth/signup", authHandler.Signup)

		r.Get("/events", eventHandler.List)
		r.Get("/events/{id}", eventHandler.Get)

		// Any logged-in account (admin or student).
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(jwtManager))

			r.Post("/auth/logout", authHandler.Logout)
		})

		// Admin only — manages events (and can see who registered).
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(jwtManager))
			r.Use(appmiddleware.RequireRole(jwt.RoleAdmin))

			r.Get("/admin/me", adminHandler.Me)

			r.Post("/events", eventHandler.Create)
			r.Put("/events/{id}", eventHandler.Update)
			r.Delete("/events/{id}", eventHandler.Delete)
			r.Get("/events/{id}/participants", participationHandler.ListParticipants)
		})

		// Student only — registers for / cancels events.
		r.Group(func(r chi.Router) {
			r.Use(appmiddleware.Auth(jwtManager))
			r.Use(appmiddleware.RequireRole(jwt.RoleStudent))

			r.Get("/student/me", studentHandler.Me)
			r.Get("/student/registrations", participationHandler.MyRegistrations)

			r.Post("/events/{id}/register", participationHandler.Register)
			r.Post("/events/{id}/cancel", participationHandler.Cancel)
		})
	})

	// Serve the static frontend directly from the Go binary for local dev.
	fileServer := http.FileServer(http.Dir("./web"))
	r.Handle("/*", fileServer)

	addr := ":" + cfg.ServerPort
	log.Printf("server listening on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
