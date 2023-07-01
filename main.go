package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

type Todo struct {
	ID   int
	Text string
}

type TodoList struct {
	Todos []Todo
}

var todos TodoList
var decoder = schema.NewDecoder()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	err := tmpl.Execute(w, todos)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var todo Todo
	err = decoder.Decode(&todo, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	todo.ID = len(todos.Todos) + 1
	todos.Todos = append(todos.Todos, todo)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteTodoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	todoID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Println(err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	for i, todo := range todos.Todos {
		if todo.ID == todoID {
			todos.Todos = append(todos.Todos[:i], todos.Todos[i+1:]...)
			break
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func deleteAllHandler(w http.ResponseWriter, r *http.Request) {
	todos.Todos = []Todo{}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/todos", addTodoHandler).Methods("POST")
	r.HandleFunc("/todos/{id}/delete", deleteTodoHandler).Methods("POST")
	r.HandleFunc("/todos/deleteAll", deleteAllHandler).Methods("POST")

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	log.Println("Server started on http://localhost:8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
