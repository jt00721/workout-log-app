const API_BASE_URL = "http://localhost:8080/api";
const workoutList = document.getElementById("workout-list");
const exerciseList = document.getElementById("exercise-list");
const addForm = document.getElementById("add-form");
const addSection = document.getElementById("add-section");
const formTitle = document.getElementById("form-title");

document.getElementById("add-workout-button").addEventListener("click", () => showAddForm("workout"));
document.getElementById("add-exercise-button").addEventListener("click", () => showAddForm("exercise"));

// Fetch workouts and display them
async function fetchWorkouts() {
    try {
        const response = await fetch(`${API_BASE_URL}/workouts`);
        const workouts = await response.json();

        workoutList.innerHTML = workouts.map(workout => `
            <div class="item">
                <h3>${workout.date} (${workout.duration} mins)</h3>
                <p>${workout.notes}</p>
                <button onclick="fetchExercises(${workout.id})">View Exercises</button>
            </div>
        `).join('');
    } catch (err) {
        console.error("Failed to fetch workouts", err);
    }
}

// Fetch exercises for a specific workout
async function fetchExercises(workoutId) {
    try {
        const response = await fetch(`${API_BASE_URL}/exercises?workout_id=${workoutId}`);
        const exercises = await response.json();

        exerciseList.innerHTML = exercises.map(exercise => `
            <div class="item">
                <h3>${exercise.name}</h3>
                <p>Sets: ${exercise.sets}, Reps: ${exercise.reps}, Weight: ${exercise.weight} kg</p>
            </div>
        `).join('');
        exerciseList.parentElement.style.display = "block";
    } catch (err) {
        console.error("Failed to fetch exercises", err);
    }
}

// Show the add form
function showAddForm(type) {
    formTitle.textContent = type === "workout" ? "Add New Workout" : "Add New Exercise";
    addForm.innerHTML = type === "workout" ? `
        <label>User ID: <input type="number" name="user_id" required></label>
        <label>Date: <input type="date" name="date" required></label>
        <label>Duration: <input type="number" name="duration" required></label>
        <label>Notes: <textarea name="notes" required></textarea></label>
        <button type="submit">Add Workout</button>
    ` : `
        <label>Workout ID: <input type="number" name="workout_id" required></label>
        <label>Name: <input type="text" name="name" required></label>
        <label>Sets: <input type="number" name="sets" required></label>
        <label>Reps: <input type="number" name="reps" required></label>
        <label>Weight: <input type="number" name="weight" required></label>
        <button type="submit">Add Exercise</button>
    `;

    addSection.style.display = "block";
    addForm.onsubmit = type === "workout" ? handleAddWorkout : handleAddExercise;
}

// Handle adding a workout
async function handleAddWorkout(e) {
    e.preventDefault();
    const formData = new FormData(addForm);
    const data = Object.fromEntries(formData);

    data.user_id = parseInt(data.user_id)
    data.duration = parseInt(data.duration);

    try {
        await fetch(`${API_BASE_URL}/workouts`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        fetchWorkouts();
        addSection.style.display = "none";
    } catch (err) {
        console.error("Failed to add workout", err);
    }
}

// Handle adding an exercise
async function handleAddExercise(e) {
    e.preventDefault();
    const formData = new FormData(addForm);
    const data = Object.fromEntries(formData);

    data.workout_id = parseInt(data.workout_id)
    data.sets = parseInt(data.sets);
    data.reps = parseInt(data.reps);
    data.weight = parseFloat(data.weight);

    console.log(`NEW DATA: ${data}`);

    try {
        await fetch(`${API_BASE_URL}/exercises`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(data),
        });
        addSection.style.display = "none";
    } catch (err) {
        console.error("Failed to add exercise", err);
    }
}

// Initialize the app
fetchWorkouts();
