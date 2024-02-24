package main

import (
	"net/http"
	"sort"
	"text/template"
)

func formHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		inputValue := r.FormValue("query")
		items := startSearchFromQuery(inputValue)
		tmpl, _ := template.ParseFiles("answer.html")
		tmpl.Execute(w, items)
	} else {
		http.ServeFile(w, r, "yourfile.html")
	}
}

func startSearchFromQuery(inputValue string) []ResultToWebDto {
	txtItems := read_txt_files_local()
	items := make([]ResultFromDto, len(txtItems))
	for i, item := range txtItems {
		items[i] = &item
	}

	// passing a sort function here determines which method should be
	// better.
	
	/*
	Return the document with most ocurrences of the query as a single word, only works for queries with one word.
	*/
	/* sort.Slice(items, func(i, j int) bool {
		return Compare(inputValue, items[i], items[j]) > 0
	}) */


	modelTdIdf := *ConstructormodelTfIdf(items)
	queryTfIdf := tfIdfQuery(inputValue, modelTdIdf)
	
	sort.Slice(items, func(i, j int) bool {
		result1 := cos_sim(tf_idf_doc(items[i].Name(), modelTdIdf), queryTfIdf)
		result2 := cos_sim(tf_idf_doc(items[j].Name(), modelTdIdf), queryTfIdf)  
		return result1 > result2
		})
	
	return Map(items, func(r ResultFromDto) ResultToWebDto {
		return r.(ResultToWebDto)
	})
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "principal.html")
	})
	http.HandleFunc("/send_query", formHandler)
	http.ListenAndServe(":8080", nil)
}
