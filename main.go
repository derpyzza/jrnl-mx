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
	Id 			int 		`json:"id"`
	Title 	string		`json:"title"`
	Date 		string		`json:"date"`
	Tags 		[]string	`json:"tags"`
	Post 		[]string		`json:"post"`
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
		file, err := os.ReadFile("newPost.html")
		check(err)
		fmt.Fprintf(w, string(file))
	}

	createPostHandler := func(w http.ResponseWriter, r *http.Request) {

		r.ParseForm()

		form := r.Form
		post := form["post-body"]
		title := form["post-title"]
		date := form["post-date"]
		tags := form["post-tags"]

		var err error
		err = nil

		// error check

		posts = append(posts,
			Post{
				Id: len(posts) + 1,
				Title: title[0],
				Post: post,
				Date: date[0],
				Tags: tags,
			},
		)

		fmt.Println("form: ", form)
		fmt.Printf(
			`title: %s 
date: %s 
tags: %s
post: %s`,
			title, date, tags, post)

		if (err == nil) {
		// no error
			fmt.Fprintf(w,
				"<button disabled class=\"outline\">Post Created Successfully</button>")
		} else {
		// error
			fmt.Fprintf(w, "<button disabled class=\"outline pico-color-red-500\">Error Creating Post</button>")
		}
		
	}

	http.HandleFunc("/", mainHandler)
	http.HandleFunc("/newPost", postHandler)
	http.HandleFunc("/createPost", createPostHandler)

	log.Printf("Starting server at port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
