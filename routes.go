package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"

	"github.com/gorilla/mux"
)

const (
	// IndexBaseURL for routing index requests.
	IndexBaseURL = "/"
	// PageBaseURL for routing page requests.
	PageBaseURL = "/"
	// PostBaseURL for routing post requests.
	PostBaseURL = "/post/"
	// AssetBaseURL for routing asset requests.
	StaticBaseURL  = "/static/"
	FaviconBaseURL = "/favicon.ico"
	StaticFolder   = "static"
)

// ErrorHandler handles the errors.
func ErrorHandler(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	Error.Println(err)
	err = RenderPage(w, "error", NewErrorContext(err))
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
}

// IndexHandler handles the index page and displays a list of the recent blog posts.
func IndexHandler(w http.ResponseWriter, r *http.Request) {
	posts := LatestPosts(10)
	err := RenderPage(w, "index", NewIndexContext(posts))
	if err != nil {
		ErrorHandler(w, err, 500)
		return
	}
}

// PostHandler handles a post request and displays the post.
func PostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context, err := NewPostContext(vars["slug"])
	if err != nil {
		ErrorHandler(w, err, 404)
		return
	}
	err = RenderPage(w, "post", context)
	if err != nil {
		ErrorHandler(w, err, 500)
		return
	}
}

// PageHandler handles a page request and displays the page.
func PageHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	context, err := NewPageContext(vars["slug"])
	if err != nil {
		ErrorHandler(w, err, 404)
		return
	}
	err = RenderPage(w, "page", context)
	if err != nil {
		ErrorHandler(w, err, 500)
		return
	}
}

// FaviconHandler
func FaviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(BlogFolder, Config.Meta.Favicon))
}

// LoadRoutes configures a new blog router.
func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix(StaticBaseURL).Handler(http.StripPrefix(StaticBaseURL, http.FileServer(http.Dir(path.Join(BlogFolder, StaticFolder)))))
	r.HandleFunc(IndexBaseURL, IndexHandler)
	if Config.Meta.Favicon != "" {
		r.HandleFunc(FaviconBaseURL, FaviconHandler)
	}
	r.HandleFunc(PostBaseURL+"{slug}", PostHandler)
	r.HandleFunc(PageBaseURL+"{slug}", PageHandler)
	Trace.Println("initialized routes")
	return r
}
