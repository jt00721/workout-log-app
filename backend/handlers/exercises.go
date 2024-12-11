package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/jt00721/workout-log-app/backend/database"
)

func ExercisesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getExercises(w, r)
	case http.MethodPost:
		addExercise(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func ExercisesByIdHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/exercises/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid exercise ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateExercise(w, r, id)
	case http.MethodDelete:
		deleteExercise(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getExercises(w http.ResponseWriter, r *http.Request) {
	db := database.ConnectDb()
	defer db.Close()

	// Optional filter by workout_id
	workoutID := r.URL.Query().Get("workout_id")
	if workoutID == "" {
		http.Error(w, "workout_id is required", http.StatusInternalServerError)
		return
	}

	rows, err := db.Query("SELECT id, name, sets, reps, weight FROM exercises WHERE workout_id = ?", workoutID)
	if err != nil {
		http.Error(w, "Failed to fetch exercises", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	exercises := []map[string]interface{}{}
	for rows.Next() {
		var id, sets, reps int
		var weight float64
		var name string

		if err := rows.Scan(&id, &name, &sets, &reps, &weight); err != nil {
			http.Error(w, "Failed to parse exercises", http.StatusInternalServerError)
			return
		}

		exercises = append(exercises, map[string]interface{}{
			"id":     id,
			"name":   name,
			"sets":   sets,
			"reps":   reps,
			"weight": weight,
		})
	}

	respondWithJSON(w, http.StatusOK, exercises)
}

func addExercise(w http.ResponseWriter, r *http.Request) {
	var exercise struct {
		WorkoutId int     `json:"workout_id"`
		Name      string  `json:"name"`
		Sets      int     `json:"sets"`
		Reps      int     `json:"reps"`
		Weight    float64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if exercise.WorkoutId == 0 || exercise.Name == "" || exercise.Sets == 0 || exercise.Reps == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDb()
	defer db.Close()

	var workoutExists bool
	err := db.QueryRow("SELECT EXISTS (SELECT 1 FROM workouts WHERE id = ?)", exercise.WorkoutId).Scan(&workoutExists)
	if err != nil || !workoutExists {
		http.Error(w, "Workout not found", http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO exercises (workout_id, name, sets, reps, weight) VALUES (?, ?, ?, ?, ?)", exercise.WorkoutId, exercise.Name, exercise.Sets, exercise.Reps, exercise.Weight)
	if err != nil {
		http.Error(w, "Failed to create exercise", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Exercise created successfully"})
}

func updateExercise(w http.ResponseWriter, r *http.Request, id int) {
	var exercise struct {
		Name   string  `json:"name"`
		Sets   int     `json:"sets"`
		Reps   int     `json:"reps"`
		Weight float64 `json:"weight"`
	}

	if err := json.NewDecoder(r.Body).Decode(&exercise); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if exercise.Name == "" || exercise.Sets == 0 || exercise.Reps == 0 {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		return
	}

	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("UPDATE exercises SET name = ?, sets = ?, reps = ?, weight = ? WHERE id = ?", exercise.Name, exercise.Sets, exercise.Reps, exercise.Weight, id)
	if err != nil {
		http.Error(w, "Failed to update exercise", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusCreated, map[string]string{"message": "Exercise updated successfully"})
}

func deleteExercise(w http.ResponseWriter, r *http.Request, id int) {
	db := database.ConnectDb()
	defer db.Close()

	_, err := db.Exec("DELETE FROM exercises WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete exercise", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Exercise deleted successfully"})
}
