package content

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"path"
	"path/filepath"
	"strings"
	"time"

	humanize "github.com/dustin/go-humanize"
	"github.com/lnsp/bloggy/pkg/config"
	"github.com/sirupsen/logrus"
)

// TemplatesFolder is the default folder for templates.
const (
	TemplateFolder = "templates"
	DisplayFolder  = "displays"
	IncludeFolder  = "includes"
)

// NavItemContext stores the information of a navigation item.
type NavItemContext struct {
	Title string
	URL   string
}

// BaseContext stores basic context information like title, author etc.
type BaseContext struct {
	BlogTitle    string
	BlogSubtitle string
	BlogAuthor   string
	BlogYear     string
	BlogEmail    string
	BlogURL      string
	BlogNav      []NavItemContext
}

// PostContext stores additional information for posts.
type PostContext struct {
	BaseContext
	PostTitle    string
	PostSubtitle string
	PostDate     string
	PostContent  template.HTML
	PostURL      string
}

// PageContext stores additional information for pages.
type PageContext struct {
	BaseContext
	PageTitle   string
	PageContent template.HTML
	PageURL     string
}

// IndexContext stores a list of the latest posts.
type IndexContext struct {
	BaseContext
	LatestPosts []Post
}

// ErrorContext stores error information.
type ErrorContext struct {
	BaseContext
	Message string
}

type Templater struct {
	Config      *config.Config
	templates   map[string]*template.Template
	navItems    []NavItemContext
	blogContext *BaseContext
	cachedPages map[string]*PageContext
	cachedPosts map[string]*PostContext
	cachedIndex *IndexContext
	index       *Index
}

func (t *Templater) ClearNav() {
	t.navItems = make([]NavItemContext, 0)
}

// AddNavItem adds a item to the navigation.
func (t *Templater) AddNavItem(e Entry) {
	item := NavItemContext{
		Title: e.GetTitle(),
		URL:   e.GetURL(),
	}
	t.navItems = append(t.navItems, item)
	logrus.WithField("title", item.Title).Debug("added item to navigation bar")
}

// NewBaseContext either creates a new BaseContext from the global blog configuration or returns the cached version.
func (t *Templater) NewBaseContext() *BaseContext {
	if t.blogContext != nil {
		return t.blogContext
	}
	t.blogContext = &BaseContext{
		BlogTitle:    t.Config.Meta.Title,
		BlogSubtitle: t.Config.Meta.Subtitle,
		BlogAuthor:   t.Config.Author.Name,
		BlogYear:     fmt.Sprint(time.Now().Year()),
		BlogEmail:    t.Config.Author.Email,
		BlogURL:      "/",
		BlogNav:      t.navItems,
	}
	return t.blogContext
}

// NewPostContext either creates a new post context or returns the cached version.
func (t *Templater) NewPostContext(slug string) (*PostContext, error) {
	context, ok := t.cachedPosts[slug]
	if !ok {
		post, ok := t.index.PostBySlug[slug]
		if !ok {
			return nil, errors.New("post not found")
		}
		context = &PostContext{
			*t.NewBaseContext(),
			post.Title,
			post.Subtitle,
			humanize.Time(post.PublishDate),
			template.HTML(Render(post)),
			post.GetURL(),
		}
		t.cachedPosts[slug] = context
		logrus.WithField("slug", slug).Debug("created cache version of post")
	}
	return context, nil
}

// NewPageContext either creates a new page context or returns the cached version.
func (t *Templater) NewPageContext(slug string) (*PageContext, error) {
	context, ok := t.cachedPages[slug]
	if !ok {
		page, ok := t.index.PageBySlug[slug]
		if !ok {
			return nil, errors.New("page '" + slug + "' not found")
		}
		context = &PageContext{
			*t.NewBaseContext(),
			page.Title,
			template.HTML(Render(page)),
			page.GetURL(),
		}
		t.cachedPages[slug] = context
		logrus.WithField("slug", slug).Debug("created cache version of page")
	}
	return context, nil
}

// NewIndexContext either creates a new index context or returns the cached version.
func (t *Templater) NewIndexContext() *IndexContext {
	if t.cachedIndex != nil {
		return t.cachedIndex
	}
	t.cachedIndex = &IndexContext{*t.NewBaseContext(), t.index.LatestPosts(10)}
	return t.cachedIndex
}

// NewErrorContext creates a new error context.
func (t *Templater) NewErrorContext(err error) *ErrorContext {
	return &ErrorContext{*t.NewBaseContext(), err.Error()}
}

// NewTemplater loads the templates from the blog folder.
func NewTemplater(cfg *config.Config, index *Index) (*Templater, error) {
	tmpl := &Templater{
		Config:      cfg,
		cachedPages: make(map[string]*PageContext),
		cachedPosts: make(map[string]*PostContext),
		templates:   make(map[string]*template.Template),
		index:       index,
	}
	displays, err := filepath.Glob(path.Join(cfg.Base, TemplateFolder, DisplayFolder, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("displays glob: %w", err)
	}
	logrus.WithField("displays", displays).Debug("loading displays")
	includes, err := filepath.Glob(path.Join(cfg.Base, TemplateFolder, IncludeFolder, "*.html"))
	if err != nil {
		return nil, fmt.Errorf("includes glob: %w", err)
	}
	logrus.WithField("includes", includes).Debug("loading includes")
	for _, display := range displays {
		files := append(includes, display)
		name := strings.TrimSuffix(filepath.Base(display), filepath.Ext(display))
		tmpl.templates[name] = template.Must(template.New(name).ParseFiles(files...))
	}
	for _, page := range index.Pages {
		tmpl.AddNavItem(&page)
	}
	for key, value := range cfg.Links {
		link := NavigationLink{key, value}
		tmpl.AddNavItem(&link)
	}
	return tmpl, nil
}

// RenderPage renders a page or throws an error if the template is missing.
func (t *Templater) RenderPage(w io.Writer, name string, context interface{}) error {
	tmpl, ok := t.templates[name]
	if !ok {
		return errors.New("template not found")
	}
	return tmpl.ExecuteTemplate(w, "base", context)
}

// ClearCache clears the context cache.
func (t *Templater) ClearCache() {
	logrus.Info("clearing cache")
	t.blogContext = nil
	t.cachedPages = make(map[string]*PageContext)
	t.cachedPosts = make(map[string]*PostContext)
	t.cachedIndex = nil
}

type NavigationLink struct {
	Title string
	URL   string
}

func (n *NavigationLink) GetContent() string {
	return ""
}

func (n *NavigationLink) GetURL() string {
	return n.URL
}

func (n *NavigationLink) GetTitle() string {
	return n.Title
}
