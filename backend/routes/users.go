package routes

import (
	"net/http"

	"github.com/jt00721/workout-log-app/backend/handlers"
)

func UserRoutes(router *http.ServeMux) {
	router.HandleFunc("/api/users", handlers.UsersHandler)
}
