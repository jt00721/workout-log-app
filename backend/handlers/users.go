package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jt00721/workout-log-app/backend/database"
)

func UsersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUsers(w, r)
	case http.MethodPost:
		createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDb()
	defer db.Close()

	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []map[string]interface{}{}
	for rows.Next() {
		var id int
		var name, email string
		if err := rows.Scan(&id, &name, &email); err != nil {
			http.Error(w, "Failed to parse users", http.StatusInternalServerError)
			return
		}

		users = append(users, map[string]interface{}{
			"id":    id,
			"name":  name,
			"email": email,
		})
	}

	respondWithJSON(w, http.StatusOK, users)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Name == "" || user.Email == "" {
		http.Error(w, "Name and Email are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", user.Name, user.Email)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "User created successfully"})
}

func respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}
