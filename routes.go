package main

import (
	"fmt"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

const (
	// IndexBaseURL for routing index requests.
	IndexBaseURL = "/"
	// PostBaseURL for routing post requests.
	PostBaseURL = "/post"
	// AssetBaseURL for routing asset requests.
	StaticBaseURL = "/static/"
	StaticFolder  = "static"
)

// ErrorHandler handles the errors.
func ErrorHandler(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	renderError := ErrorTemplate.ExecuteTemplate(w, "base", GetErrorContext(err))
	if renderError != nil {
		Error.Println("Failed to render template:", renderError)
	}
}

// IndexHandler handles the index page and displays a list of the recent blog posts.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	posts := GetLatestsPosts(10)
	context := GetIndexContext(posts)

	renderError := IndexTemplate.ExecuteTemplate(w, "base", context)
	if renderError != nil {
		Warning.Println("Failed to render template:", renderError)
		fmt.Fprintln(w, "Sorry, an error occured. Please try again later.")
		return
	}
}

// PostHandler handles a post request and displays the post.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	post, notFoundError := FindPost(vars["slug"])
	if notFoundError != nil {
		ErrorHandler(w, notFoundError, 404)
		return
	}
	renderError := PostTemplate.ExecuteTemplate(w, "base", GetPostContext(*post))
	if renderError != nil {
		ErrorHandler(w, renderError, 500)
		Error.Println("Failed to render template:", renderError)
		return
	}
}

// LoadRoutes configures a new blog router.
func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix(StaticBaseURL).Handler(http.StripPrefix(StaticBaseURL, http.FileServer(http.Dir(path.Join(BlogFolder, StaticFolder)))))
	r.HandleFunc(IndexBaseURL, IndexHandler)
	r.HandleFunc(PostBaseURL+"/{slug}", PostHandler)
	Trace.Println("Initialized default router")
	return r
}
