// Package httpx wires the HTTP router and middleware.
package httpx

import (
	"edugov-back-v2/internal/http/handlers"
	appMiddleware "edugov-back-v2/internal/http/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// App holds the configured router.
type App struct {
	Router http.Handler
}

// NewRouter builds the chi router with all routes and middleware.
func NewRouter(h *handlers.Handlers) *App {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	r.Route("/emplpoyee-degree", func(r chi.Router) {
		r.Get("/", h.GetEmployeeDegrees)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeDegree)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeDegree)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeDegree)
	})

	r.Route("/emplpoyee-main-research-area", func(r chi.Router) {
		r.Get("/", h.GetEmployeeMainResearchAreas)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeMainResearchArea)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeMainResearchArea)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeMainResearchArea)
	})

	r.Route("/emplpoyee-participation-in-event", func(r chi.Router) {
		r.Get("/", h.GetEmployeeParticipationInEvents)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeParticipationInEvent)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeParticipationInEvent)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeParticipationInEvent)
	})

	r.Route("/emplpoyee-participation-in-professional-community", func(r chi.Router) {
		r.Get("/", h.GetEmployeeParticipationInProfessionalCommunities)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeParticipationInProfessionalCommunity)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeParticipationInProfessionalCommunity)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeParticipationInProfessionalCommunity)
	})

	r.Route("/emplpoyee-patent", func(r chi.Router) {
		r.Get("/", h.GetEmployeePatents)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeePatent)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeePatent)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeePatent)
	})

	r.Route("/emplpoyee-publication", func(r chi.Router) {
		r.Get("/", h.GetEmployeePublications)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeePublication)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeePublication)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeePublication)
	})

	r.Route("/emplpoyee-refresher-course", func(r chi.Router) {
		r.Get("/", h.GetEmployeeRefresherCourses)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeRefresherCourse)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeRefresherCourse)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeRefresherCourse)
	})

	r.Route("/emplpoyee-research-activity", func(r chi.Router) {
		r.Get("/", h.GetEmployeeResearchActivities)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeResearchActivity)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeResearchActivity)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeResearchActivity)
	})

	r.Route("/emplpoyee-scientific-award", func(r chi.Router) {
		r.Get("/", h.GetEmployeeScientificAwards)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeScientificAward)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeScientificAward)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeScientificAward)
	})

	r.Route("/emplpoyee-social", func(r chi.Router) {
		r.Get("/", h.GetEmployeeSocials)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeSocial)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeSocial)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeSocial)
	})

	r.Route("/emplpoyee-work-experience", func(r chi.Router) {
		r.Get("/", h.GetEmployeeWorkExperiences)
		r.With(appMiddleware.JWTMiddleware()).Post("/", h.CreateEmployeeWorkExperience)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeWorkExperience)
		r.With(appMiddleware.JWTMiddleware()).Delete("/{id}", h.DeleteEmployeeWorkExperience)
	})

	r.Route("/emplpoyee-detail", func(r chi.Router) {
		r.Get("/", h.GetEmployeeDetails)
		r.With(appMiddleware.JWTMiddleware()).Put("/", h.UpdateEmployeeDetails)
	})

	r.Route("/employee", func(r chi.Router) {
		r.Get("/", h.GetEmployee)
		r.Get("/profile-picture/{uid}", h.GetEmployeeProfilePicture)
		r.With(appMiddleware.JWTMiddleware()).Put("/profile-picture/{uid}", h.UpdateEmployeeProfilePicture)
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)
		r.Post("/logout", h.Logout)
	})

	return &App{Router: r}
}
