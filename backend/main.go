package main

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/jt00721/workout-log-app/backend/database"
)

func main() {
	db := database.ConnectDb()
	defer db.Close()

	database.InitializeSchema(db)
	// database.SeedDatabase(db)

	fs := http.FileServer(http.Dir(filepath.Join("..", "frontend")))
	http.Handle("/", fs)

	log.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
