package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"simple-task-api-golang/internal/service"
	"simple-task-api-golang/internal/store"
	"simple-task-api-golang/internal/transport"
	"strings"

	_ "modernc.org/sqlite"
)



func main() {

	var port = "8080"

	// Connect to SQLite
	db, err := sql.Open("sqlite", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close() // Close connection

	// Create table if not exists
	q := `
		CREATE TABLE IF NOT EXISTS tasks(
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			completed INTEGER NOT NULL)`

	if _, err := db.Exec(q); err != nil {
		log.Fatal(err.Error())
	}

	// Dependency inyection
	taskStore := store.New(db)
	taskService := service.New(taskStore)
	taskHandler := transport.New(taskService)

	// Route config
	http.HandleFunc("/tasks",taskHandler.HandleTasks)
	http.HandleFunc("/tasks/",taskHandler.HandleTaskById)

	fmt.Println(strings.TrimSpace("	Server running in http://localhost:"+port))
	fmt.Println("	API endpoints:")
	fmt.Println("	GET    /tasks        - Get all tasks")
	fmt.Println("	POST   /tasks        - Create a new task")
	fmt.Println("	GET    /tasks/{id}   - Get specific task")
	fmt.Println("	PUT    /tasks/{id}   - Update a task")
	fmt.Println("	DELETE /tasks/{id}   - Delete a task")

	port = ":" + port

	// Start and hear server
	log.Fatal(http.ListenAndServe(port,nil))

}
