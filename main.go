package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

// Todo struct represents a single to-do item.
type Todo struct {
	Task      string
	Completed bool
}

// A thread-safe in-memory store for our to-do items.
var (
	todos []Todo
	mu    sync.RWMutex
)

// Global template object.
var tmpl *template.Template

// main is the entry point of our application.
func main() {
	// Parse the index.html template from the "static" directory.
	tmpl = template.Must(template.ParseFiles("static/index.html"))

	// Add some initial data.
	todos = []Todo{
		{Task: "Learn Go", Completed: true},
		{Task: "Build a web app", Completed: false},
	}

	// This creates a file server to serve files from the "static" directory.
	// For example, if you add an image at "static/images/logo.png",
	// it will be accessible at "http://localhost:8080/static/images/logo.png".
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers for our application logic.
	http.HandleFunc("/", todoHandler)
	http.HandleFunc("/add-todo", addTodoHandler)
	http.HandleFunc("/delete-todo", deleteTodoHandler)

	// Start the server.
	log.Println("Server starting on :8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// todoHandler renders the main page with the list of to-do items.
func todoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mu.RLock()
	defer mu.RUnlock()

	err := tmpl.Execute(w, todos)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

// addTodoHandler handles adding a new to-do item.
func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	task := r.FormValue("task")
	if task == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	newTodo := Todo{Task: task, Completed: false}
	todos = append(todos, newTodo)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// deleteTodoHandler handles deleting a to-do item.
func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    taskToDelete := r.FormValue("task")
    if taskToDelete == "" {
        http.Error(w, "Bad Request: Missing task", http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    for i, todo := range todos {
        if todo.Task == taskToDelete {
            todos = append(todos[:i], todos[i+1:]...)
            break
        }
    }

    http.Redirect(w, r, "/", http.StatusSeeOther)
}