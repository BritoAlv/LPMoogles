package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

type ResultToWebDto interface {
	Name() string
	Snippet() string
}

type ResultFromDto interface {
	ResultToWebDto
	Text() string
}

type ResultFromTxtDto struct {
	path string
}

func (r *ResultFromTxtDto) Name() string {
	return filepath.Base(r.path)
}

func (r *ResultFromTxtDto) Snippet() string {
	maxLen := 200
	// read the first max_len characters from the file.
	// if the file is shorter than max_len, return the whole file.
	// if the file is empty, return an empty string.
	tt := r.Text()
	if len(tt) < maxLen {
		return tt
	} else {
		return tt[:maxLen]
	}
}

func (r *ResultFromTxtDto) Text() string {
	file, err := os.Open(r.path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	buf := make([]byte, 100000)
	n, err := file.Read(buf)
	if err != nil {
		log.Fatal(err)
	}
	return string(buf[:n])
}

func CountOccurrences(word string, text string) int {
	// count number of occurrences of word in text.
	words := strings.Split(text, " ")
	count := 0
	for _, w := range words {
		if w == word {
			count++
		}
	}
	return count
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func Compare(word string, a, b ResultFromDto) int {
	count1 := CountOccurrences(word, a.Text())
	count2 := CountOccurrences(word, b.Text())
	if count1 < count2 {
		return -1
	} else if count1 > count2 {
		return 1
	} else {
		return BoolToInt(a.Name() < b.Name())
	}
}

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

func Map[T any, U any](vs []T, f func(T) U) []U {
	vsm := make([]U, len(vs))
	for i, v := range vs {
		vsm[i] = f(v)
	}
	return vsm
}

func startSearchFromQuery(inputValue string) []ResultToWebDto {
	txtItems := read_txt_files_local()
	items := make([]ResultFromDto, len(txtItems))
	for i, item := range txtItems {
		items[i] = &item
	}
	sort.Slice(items, func(i, j int) bool {
		return Compare(inputValue, items[i], items[j]) > 0
	})

	for _, b := range items {
		fmt.Println(b.Name())
	}

	return Map(items, func(r ResultFromDto) ResultToWebDto {
		return r.(ResultToWebDto)
	})
}

func read_txt_files_local() []ResultFromTxtDto {
	files, err := os.ReadDir("./database/")
	if err != nil {
		log.Fatal(err)
	}
	var items []ResultFromTxtDto
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".txt" {
			items = append(items, ResultFromTxtDto{"./database/" + file.Name()})
		}
	}
	return items
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "principal.html")
	})
	http.HandleFunc("/send_query", formHandler)
	http.ListenAndServe(":8080", nil)
}
