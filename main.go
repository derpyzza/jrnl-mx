package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"sort"

	"github.com/charmbracelet/log"
)

var postString = `
{{range .}}
<article>
<header>
	<h6>{{.Title}}</h6>
	<small>22/12/25</small>
</header>

{{range $post := .Post}}
<p>{{$post}}</p>
{{end}}
<footer>
	<div class="group">
		<div>
		{{range $key, $value := .Tags}}
		<small>#{{$value}}</small>
		{{end}}
		</div>
	</div>				
</footer>
</article>
{{end}}`

type Post struct {
	Id 			int 		`json:"id"`
	Title 	string		`json:"title"`
	Date 		string		`json:"date"`
	Tags 		[]string	`json:"tags"`
	Post 		[]string		`json:"post"`
}

var posts []Post

func check(err error){}

func main() {

	port := ":5050"
	dbfile := "data.json"

	log.SetLevel(log.DebugLevel)

	file, err := os.ReadFile(dbfile)
	check(err)

	if err := json.Unmarshal(file, &posts); err != nil {
		panic(err)
	}

	postHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("newPost.html")
		check(err)
		fmt.Fprintf(w, string(file))
	}

	mux := http.NewServeMux()

	// logger middleware
	logger := func(handler string, next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			log.Info("Incoming request", "Handler", handler, "Method", r.Method, "URL", r.URL)
			next.ServeHTTP(w, r)
		}
	}

	mux.HandleFunc("/", logger("/", mainHandler))
	mux.HandleFunc("/newPost/", logger("/newPost/", postHandler))
	mux.HandleFunc("/createPost/", logger("/createPost/", createPostHandler))
	mux.HandleFunc("/oldest/", logger("/oldest/",
		func (w http.ResponseWriter, r *http.Request){sortHandler(false, w, r)}))
	mux.HandleFunc("/recent/", logger("/recent/",
		func (w http.ResponseWriter, r *http.Request){sortHandler(true, w, r)}))
	mux.HandleFunc("/sort", logger("/sort", ssortHandler))

	log.Printf("Starting server at port %s...", port)
	log.Fatal(http.ListenAndServe(port, mux))
}

func ssortHandler( w http.ResponseWriter, r *http.Request ) {
	query := r.URL.Query().Get("recent")
	log.Debug("hiiiit")
	if (query == "true"){
		sorted := make([]Post, len(posts))
		copy(sorted, posts)
		// sort list back to front
		sort.Slice(sorted, func(i, j int) bool {
			res := sorted[i].Id > sorted[j].Id
			return res
		})
		templ := template.Must(template.New("post").Parse(postString))
		templ.Execute(w, sorted)
		return
	} else if (query == "false") {
		
		log.Debug("oldest sorting\n")
		
		templ := template.Must(template.New("post").Parse(postString))
		templ.Execute(w, posts)
		return
	} else {
		fmt.Fprintf(w, "unknown query")
	}
}

func createPostHandler (w http.ResponseWriter, r *http.Request) {
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
	if (err == nil) {
	// no error
		fmt.Fprintf(w,
			"<button disabled class=\"outline\">Post Created Successfully</button>")
	} else {
	// error
		fmt.Fprintf(w, "<button disabled class=\"outline pico-color-red-500\">Error Creating Post</button>")
	}
}	

func mainHandler ( w http.ResponseWriter, r *http.Request ) {

	if r.URL.Path != "/" {
		log.Error("Unknown path", "path", r.URL.Path)
		http.NotFound(w, r)
		return
	}

	temp := template.Must(template.ParseFiles("index.html"))

	temp.Execute(w, posts);
}

func sortHandler (s bool, w http.ResponseWriter, r *http.Request) {
		// s 	== sort recent
		// !s == sort oldest	- It's stupid ik, but it saves space!
		// In other words, s represents whether or not to sort the posts list
		// or return the original list
		log.Info("s: ", s)

		if (!s) {
			// return list as it is
			templ := template.Must(template.New("post").Parse(postString))
			templ.Execute(w, posts)
			// log.Print(posts)
			return
		}
		// else sort list
		sorted := posts

		// sort list back to front
		sort.Slice(sorted, func(i, j int) bool {
			res := sorted[i].Id > sorted[j].Id
			return res
		})
		// log.Print(sorted)
		templ := template.Must(template.New("post").Parse(postString))
		templ.Execute(w, sorted)
		// log.Print(sorted)
		return
	}
