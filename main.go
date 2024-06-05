package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/signal"
	"sort"
	"strings"
	"syscall"
	"time"
	"errors"

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
			<small>
				<a 
					hx-get="/tags?tag={{$value}}"
					hx-target="#posts"
					hx-swap="innerHTML"
					hx-push-url="/tags?tag={{$value}}"
					>
					#{{$value}}</a></small>
		{{end}}
		</div>
	</div>				
</footer>
</article>
{{end}}`


var postTempl = template.Must(template.New("post").Parse(postString))

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
	mux := http.NewServeMux()

	server := http.Server {
		Addr: port,
		Handler: mux,
	}

	log.SetLevel(log.DebugLevel)

	file, err := os.ReadFile(dbfile)
	check(err)

	if err := json.Unmarshal(file, &posts); err != nil {
		panic(err)
	}

	postHandler := func(w http.ResponseWriter, r *http.Request) {
		file, err := os.ReadFile("newPost.html")
		check(err)
		w.Header().Add("Content-Type", "text/html")
		fmt.Fprintf(w, string(file))
	}

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
	mux.HandleFunc("/sort", logger("/sort", sortHandler))
	mux.HandleFunc("/tags", logger("/tags", tagHandler))

	//== MAGIC CODE==
	// magic server shutdown code. idk how this works.
	go func () {
		log.Info("Starting server", "port", port)
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()

	log.Info("Shutting down server...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP Shutdown error: %v", err)
		server.Close()
	}
	//== MAGIC CODE==

	b, err := json.Marshal(posts)
	f, err := os.Create(dbfile)
	defer f.Close()

	fmt.Fprintf(f, string(b))
	log.Debug("this code is reached")
}

func sortHandler( w http.ResponseWriter, r *http.Request ) {
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
		postTempl.Execute(w, sorted)
		return
	} else if (query == "false") {
		
		log.Debug("oldest sorting\n")
		
		postTempl.Execute(w, posts)
		return
	} else {
		fmt.Fprintf(w, "unknown query")
	}
}

func tagHandler (w http.ResponseWriter, r *http.Request ) {
	tag := r.URL.Query().Get("tag")
	var vposts []Post

	for _, post := range posts {
		for _, t := range post.Tags {
			if t == tag {
				vposts = append(vposts, post)
			}
		}
	}

	if len(vposts) >= 1 {
		postTempl.Execute(w, vposts);
		return
	}

	fmt.Fprintf(w, "no tags found :/")
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
	log.Debug(form)
	log.Debug(tags)
	for index, tag := range tags {
		tags[index] = strings.ToLower(tag)	
	}

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
		errortemp := template.Must(template.ParseFiles("error.html"))

		errortemp.Execute(w, r.URL);
		return
	}

	temp := template.Must(template.ParseFiles("index.html"))

	temp.Execute(w, posts);
}
