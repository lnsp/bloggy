package routes

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/mooxmirror/blog/config"
	"github.com/mooxmirror/blog/templates"
	"github.com/mooxmirror/blog/posts"
)

var cfg config.Config

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "happy new year from %s!", cfg.HostCountry)
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	post := posts.PublicPosts[0]
	context := templates.GetPostContext(cfg, post)

	renderError := templates.PostTemplate.Execute(w, context)
	if renderError != nil {
		log.Fatal(renderError)
		fmt.Fprintln(w, "Sorry, an error occured. Please try again later.")
		return
	}
}

func Setup(configuration config.Config) (*mux.Router) {
	cfg = configuration
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	r.HandleFunc("/post", PostHandler)
	return r
}
