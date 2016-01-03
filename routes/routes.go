package routes

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/mooxmirror/blog/config"
)

var cfg config.Config

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "happy new year from %s!", cfg.HostCountry)
}

func Setup(configuration config.Config) (*mux.Router) {
	cfg = configuration
	r := mux.NewRouter()
	r.HandleFunc("/", HelloHandler)
	return r
}
