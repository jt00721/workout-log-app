package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/jt00721/workout-log-app/backend/database"
	"github.com/jt00721/workout-log-app/backend/routes"
)

func main() {
	db := database.ConnectDb()
	defer db.Close()

	database.InitializeSchema(db)
	// database.SeedDatabase(db)

	router := http.NewServeMux()

	routes.UserRoutes(router)
	routes.WorkoutRoutes(router)
	routes.ExerciseRoutes(router)

	fs := http.FileServer(http.Dir(filepath.Join("..", "frontend")))
	router.Handle("/", fs)

	log.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
