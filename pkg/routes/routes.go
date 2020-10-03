package routes

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
	"github.com/lnsp/bloggy/pkg/config"
	"github.com/lnsp/bloggy/pkg/content"
	"github.com/sirupsen/logrus"
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

type Router struct {
	mux       *mux.Router
	templater *content.Templater
	config    *config.Config
}

// ErrorHandler handles the errors.
func (router *Router) error(w http.ResponseWriter, err error, status int) {
	w.WriteHeader(status)

	logrus.WithError(err).Error("failed to render page")
	err = router.templater.RenderPage(w, "error", router.templater.NewErrorContext(err))
	if err != nil {
		fmt.Fprintln(w, err.Error())
		return
	}
}

// IndexHandler handles the index page and displays a list of the recent blog posts.
func (router *Router) indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := router.templater.RenderPage(w, "index", router.templater.NewIndexContext())
		if err != nil {
			router.error(w, err, 500)
			return
		}
	})
}

// PostHandler handles a post request and displays the post.
func (router *Router) postHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		context, err := router.templater.NewPostContext(vars["slug"])
		if err != nil {
			router.error(w, err, 404)
			return
		}
		err = router.templater.RenderPage(w, "post", context)
		if err != nil {
			router.error(w, err, 500)
			return
		}
	})
}

// PageHandler handles a page request and displays the page.
func (router *Router) pageHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		context, err := router.templater.NewPageContext(vars["slug"])
		if err != nil {
			router.error(w, err, 404)
			return
		}
		if err := router.templater.RenderPage(w, "page", context); err != nil {
			router.error(w, err, 500)
			return
		}
	})
}

// FaviconHandler initializes a new favicon handler.
func (router *Router) faviconHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(router.config.Base, router.config.Meta.Favicon))
	})
}

// Save renders all routes statically to a local directory.
func (router *Router) Save(dir string) error {
	return nil
}

// Serve waits for incoming connections on the configured port.
func (router *Router) Serve() error {
	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", router.config.Server.Port),
		ReadHeaderTimeout: time.Minute,
		ReadTimeout:       time.Minute,
		WriteTimeout:      time.Minute,
		Handler:           router.mux,
	}
	return server.ListenAndServe()
}

// NewRouter configures a new blog router.
func NewRouter(cfg *config.Config, templater *content.Templater) *Router {
	rtr := &Router{
		mux:       mux.NewRouter(),
		templater: templater,
		config:    cfg,
	}
	rtr.mux.PathPrefix(StaticBaseURL).Handler(
		http.StripPrefix(StaticBaseURL, http.FileServer(http.Dir(path.Join(cfg.Base, StaticFolder)))))
	rtr.mux.Handle(IndexBaseURL, rtr.indexHandler())
	if cfg.Meta.Favicon != "" {
		rtr.mux.Handle(FaviconBaseURL, rtr.faviconHandler())
	}
	rtr.mux.Handle(PostBaseURL+"{slug}", rtr.postHandler())
	rtr.mux.Handle(PageBaseURL+"{slug}", rtr.pageHandler())
	return rtr
}

type simpleResolver struct {
	cfg *config.Config
}

func (r *simpleResolver) Page(slug string) string {
	return fmt.Sprintf("%s%s", PageBaseURL, slug)
}

func (r *simpleResolver) Post(slug string) string {
	return fmt.Sprintf("%s%s", PostBaseURL, slug)
}

// NewResolver creates a new simple URL resolver.
func NewResolver(cfg *config.Config) content.URLResolver {
	return &simpleResolver{cfg}
}
