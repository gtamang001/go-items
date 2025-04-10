package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type Item struct {
	ID   int
	Name string
}

var items []Item
var idCounter int

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	log.Println("Starting server on http://localhost:8080")

	http.HandleFunc("/", loggingMiddleware(indexHandler))
	http.HandleFunc("/create", loggingMiddleware(createHandler))
	http.HandleFunc("/update", loggingMiddleware(updateHandler))
	http.HandleFunc("/delete", loggingMiddleware(deleteHandler))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "index.html", items)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		idCounter++
		item := Item{ID: idCounter, Name: name}
		items = append(items, item)
		log.Printf("Created item: %+v\n", item)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	templates.ExecuteTemplate(w, "create.html", nil)
}

func updateHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	for i := range items {
		if items[i].ID == id {
			if r.Method == http.MethodPost {
				oldName := items[i].Name
				items[i].Name = r.FormValue("name")
				log.Printf("Updated item ID %d: '%s' → '%s'\n", id, oldName, items[i].Name)
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}
			templates.ExecuteTemplate(w, "update.html", items[i])
			return
		}
	}
	http.NotFound(w, r)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, _ := strconv.Atoi(idStr)

	for i, item := range items {
		if item.ID == id {
			log.Printf("Deleted item: %+v\n", item)
			items = append(items[:i], items[i+1:]...)
			break
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// middleware to log requests
func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.URL.Path)
		next(w, r)
	}
}
