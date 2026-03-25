package main

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestTodoHandler checks if the main page loads with a 200 OK status.
func TestTodoHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(todoHandler)

	// Since the global tmpl is needed, we initialize it for the test
	// from the new "static" directory.
	tmpl = template.Must(template.ParseFiles("static/index.html"))

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}
