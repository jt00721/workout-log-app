package routes

import (
	"net/http"

	"github.com/jt00721/workout-log-app/backend/handlers"
)

func ExerciseRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/exercises", handlers.ExercisesHandler)
	router.HandleFunc("/api/exercises/", handlers.ExercisesByIdHandler) // Handles PUT, DELETE
}
