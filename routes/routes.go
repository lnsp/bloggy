package routes

import (
	"fmt"
	"log"
	"net/http"
	"html/template"
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
	/*context := struct {
		BlogTitle, BlogSubtitle, BlogAuthor, BlogYear, BlogEmail string
		PostTitle, PostSubtitle string
		PostContent template.HTML
	}*/
	postContext := struct {
		PostTitle, PostSubtitle, PostDate string
		PostContent template.HTML
	}{post.Title, post.Subtitle, post.PublishDate.String(), template.HTML(post.HTMLContent)}

	renderError := templates.PostTemplate.Execute(w, postContext)
	if renderError != nil {
		log.Fatal(renderError)
		fmt.Fprintln(w, renderError)
		return
	}
}

func Setup(configuration config.Config) (*mux.Router) {
	cfg = configuration
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	return r
}
