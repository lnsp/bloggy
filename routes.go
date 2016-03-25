package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

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
	post := BlogPosts[0]
	context := GetPostContext(post)

	renderError := PostTemplate.ExecuteTemplate(w, "base", context)
	if renderError != nil {
		Warning.Println("Failed to render template:", renderError)
		fmt.Fprintln(w, "Sorry, an error occured. Please try again later.")
		return
	}
}

// LoadRoutes configures a new blog router.
func LoadRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/", IndexHandler)
	r.HandleFunc("/post", PostHandler)
	Info.Println("Initialized default router")
	return r
}
