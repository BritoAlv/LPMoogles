package main

import (
	"net/http"
	"text/template"
)

type Item struct {
	Name    string
	Snippet string
}

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		inputValue := r.FormValue("query")
		items := []Item{
			{inputValue, "BritoAlv"},
			{inputValue, "limaJavier"},
			{inputValue, "AlbaroAlb"},
		}
		tmpl, _ := template.ParseFiles("answer.html")
		tmpl.Execute(w, items)
	} else {
		http.ServeFile(w, r, "yourfile.html")
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "principal.html")
	})
	http.HandleFunc("/send_query", formHandler)
	http.ListenAndServe(":8080", nil)
}
