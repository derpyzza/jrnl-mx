package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Post struct {
	Id 			string 		`json:"id"`
	Title 	string		`json:"title"`
	Date 		string		`json:"date"`
	Tags 		[]string	`json:"tags"`
	Post 		string		`json:"post"`

	Saved 	bool			`json:"saved"`
}

func check(err error){}

func main() {

	var posts []Post

	port := ":5050"
	dbfile := "data.json"

	file, err := os.ReadFile(dbfile)
	check(err)

	if err := json.Unmarshal(file, &posts); err != nil {
		panic(err)
	}

	mainHandler := func ( w http.ResponseWriter, r *http.Request ) {
		temp := template.Must(template.ParseFiles("index.html"))

		temp.Execute(w, posts);

		log.Printf("Request from %s", r.URL)
	}

	postHandler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `<p>Hello</p>
<button
	hx-post="/home"
	hx-target="#container"
	hx-swap="innerHTML"
	>
	Go Back
</button>
		`)
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/home", mainHandler)
	http.HandleFunc("/newPost", postHandler)

	log.Printf("Starting server at port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
