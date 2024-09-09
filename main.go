package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

type Router struct {
	Name        string
	Method      string
	Path        string
	HandlerFunc http.HandlerFunc
}

type Routes []Router

type Todo struct {
	Id      int
	Message string
}

var data = map[string][]Todo{}

func todosHandler(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("index.html"))

	// data["Todos"] = append(data["Todos"], Todo{Id: 0, Message: "Pray to God"})

	err := templ.Execute(w, data)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)

		return
	}
}

func addTodoHandler(w http.ResponseWriter, r *http.Request) {
	message := r.PostFormValue("message")
	todo := Todo{Id: len(data["Todos"]) + 1, Message: message}
	data["Todos"] = append(data["Todos"], todo)

	templ := template.Must(template.ParseFiles("index.html"))
	err := templ.ExecuteTemplate(w, "todo-list-element", todo)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

func createRoutes() *Routes {
	routes := &Routes{
		Router{
			Name:   "Index Page",
			Method: "GET",
			Path:   "/",
			HandlerFunc: func(w http.ResponseWriter, r *http.Request) {
				todosHandler(w, r)
			},
		},
		Router{
			Name:        "Add Todo",
			Method:      "POST",
			Path:        "/add-todo",
			HandlerFunc: addTodoHandler,
		},
	}

	return routes
}

func main() {
	routes := createRoutes()
	router := mux.NewRouter().StrictSlash(true)

	for _, route := range *routes {
		router.
			Methods(route.Method).
			Path(route.Path).
			Name(route.Name).
			Handler(route.HandlerFunc)
	}

	fs := http.FileServer(http.Dir("assets/"))

	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", fs))

	port := "8082"

	log.Printf("Listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
