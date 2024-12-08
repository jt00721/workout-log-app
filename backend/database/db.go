package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectDb() *sql.DB {
	db, err := sql.Open("sqlite3", "database/workoutlog.db")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		log.Fatal("Failed to enable foreign keys:", err)
	}

	db.SetMaxOpenConns(1)
	return db
}

func InitializeSchema(db *sql.DB) {
	_, err := db.Exec("create table if not exists users(id integer primary key autoincrement, name text not null, email text not null unique)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table if not exists workouts(id integer primary key autoincrement, user_id integer not null, date text not null, duration integer not null, notes text, foreign key (user_id) references users (id) on delete cascade)")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("create table if not exists exercises(id integer primary key autoincrement, workout_id integer not null, name text not null, sets integer not null, reps integer not null, weight real not null, foreign key (workout_id) references workouts (id) on delete cascade)")
	if err != nil {
		log.Fatal(err)
	}
}

func SeedDatabase(db *sql.DB) {
	_, err := db.Exec("INSERT INTO users (name, email) VALUES ('Julian', 'julian@example.com')")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO workouts (user_id, date, duration, notes) VALUES (1, '2024-11-29', 60, 'Leg Day');")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("INSERT INTO exercises (workout_id, name, sets, reps, weight) VALUES (1, 'Squats', 3, 12, 70.5);")
	if err != nil {
		log.Fatal(err)
	}
}
