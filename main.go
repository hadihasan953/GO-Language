package main

import (
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseFiles("templates/index.html", "templates/result.html"))

func main() {
	mux := http.NewServeMux()

	// Static file server for CSS and assets
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/submit", submitHandler)

	addr := ":8080"
	log.Printf("Starting server at http://localhost%s/", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	if err := templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type FormData struct {
	Username string
	Email    string
}

func submitHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "invalid form", http.StatusBadRequest)
		return
	}
	data := FormData{
		Username: r.PostFormValue("username"),
		Email:    r.PostFormValue("email"),
	}

	if err := templates.ExecuteTemplate(w, "result.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
