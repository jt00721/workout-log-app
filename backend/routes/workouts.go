package routes

import (
	"net/http"

	"github.com/jt00721/workout-log-app/backend/handlers"
)

func WorkoutRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/workouts", handlers.WorkoutsHandler)      // Handles GET, POST
	router.HandleFunc("/api/workouts/", handlers.WorkoutsByIdHandler) // Handles PUT, DELETE
}
