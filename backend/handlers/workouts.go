package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jt00721/workout-log-app/backend/database"
)

func WorkoutsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWorkouts(w, r)
	case http.MethodPost:
		createWorkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func WorkoutsByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/workouts/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid workout ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateWorkout(w, r, id)
	case http.MethodDelete:
		deleteWorkout(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getWorkouts(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDb()
	defer db.Close()

	// Optional filter by user_id
	userID := r.URL.Query().Get("user_id")
	var rows *sql.Rows
	var err error

	if userID != "" {
		rows, err = db.Query("SELECT * FROM workouts WHERE user_id = ?", userID)
	} else {
		rows, err = db.Query("SELECT * FROM workouts")
	}

	if err != nil {
		http.Error(w, "Failed to fetch workouts", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workouts := []map[string]interface{}{}
	for rows.Next() {
		var id, userID int
		var date, notes string
		var duration int
		if err := rows.Scan(&id, &userID, &date, &duration, &notes); err != nil {
			http.Error(w, "Failed to parse workouts", http.StatusInternalServerError)
			return
		}
		workouts = append(workouts, map[string]interface{}{
			"id":       id,
			"user_id":  userID,
			"date":     date,
			"duration": duration,
			"notes":    notes,
		})
	}

	respondWithJSON(w, http.StatusOK, workouts)
}

func createWorkout(w http.ResponseWriter, r *http.Request) {
	var workout struct {
		UserId   int    `json:"user_id"`
		Date     string `json:"date"`
		Duration int    `json:"duration"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if workout.UserId == 0 || workout.Date == "" || workout.Duration == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("INSERT INTO workouts (user_id, date, duration, notes) VALUES (?, ?, ?, ?)", workout.UserId, workout.Date, workout.Duration, workout.Notes)
	if err != nil {
		http.Error(w, "Failed to create workout", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Workout created successfully"})
}

func updateWorkout(w http.ResponseWriter, r *http.Request, id int) {
	var workout struct {
		Date     string `json:"date"`
		Duration int    `json:"duration"`
		Notes    string `json:"notes"`
	}

	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if workout.Date == "" || workout.Duration == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("UPDATE workouts SET date = ?, duration = ?, notes = ? WHERE id = ?", workout.Date, workout.Duration, workout.Notes, id)
	if err != nil {
		http.Error(w, "Failed to update workout", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Workout updated successfully"})
}

func deleteWorkout(w http.ResponseWriter, r *http.Request, id int) {
	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("DELETE FROM workouts WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete workout", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Workout deleted successfully"})
}
