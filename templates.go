package main

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
)

// TemplatesFolder is the default folder for templates.
const (
	TemplateFolder = "templates"
	DisplayFolder  = "displays"
	IncludeFolder  = "includes"
)

var (
	templates   map[string]*template.Template
	navItems    []NavItemContext
	blogContext *BaseContext
	cachedPages map[string]*PageContext
	cachedPosts map[string]*PostContext
	cachedIndex *IndexContext
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

// AddNavItem adds a item to the navigation.
func AddNavItem(e Entry) {
	item := NavItemContext{
		Title: e.GetTitle(),
		URL:   e.GetURL(),
	}
	navItems = append(navItems, item)
	Trace.Println("added", item.Title, "to navigation bar")
}

// NewBaseContext either creates a new BaseContext from the global blog configuration or returns the cached version.
func NewBaseContext() *BaseContext {
	if blogContext == nil {
		// Generate slice of page contexts
		blogContext = &BaseContext{
			BlogTitle:    GlobalConfig.Title,
			BlogSubtitle: GlobalConfig.Subtitle,
			BlogAuthor:   GlobalConfig.Author,
			BlogYear:     fmt.Sprint(time.Now().Year()),
			BlogEmail:    GlobalConfig.Email,
			BlogURL:      GlobalConfig.URL,
			BlogNav:      navItems,
		}
	}
	return blogContext
}

// NewPostContext either creates a new post context or returns the cached version.
func NewPostContext(slug string) (*PostContext, error) {
	context, ok := cachedPosts[slug]
	if !ok {
		post, ok := PostBySlug[slug]
		if !ok {
			return nil, errors.New("post not found")
		}
		context = &PostContext{
			*NewBaseContext(),
			post.Title,
			post.Subtitle,
			humanize.Time(post.PublishDate),
			template.HTML(Render(post)),
			post.GetURL(),
		}
		cachedPosts[slug] = context
		Trace.Println("create cache version of post", slug)
	}
	return context, nil
}

// NewPageContext either creates a new page context or returns the cached version.
func NewPageContext(slug string) (*PageContext, error) {
	context, ok := cachedPages[slug]
	if !ok {
		page, ok := PageBySlug[slug]
		if !ok {
			return nil, errors.New("page '" + slug + "' not found")
		}
		context = &PageContext{
			*NewBaseContext(),
			page.Title,
			template.HTML(Render(page)),
			page.GetURL(),
		}
		cachedPages[slug] = context
		Trace.Println("create cache version of page", slug)
	}
	return context, nil
}

// GetIndexContext either creates a new index context or returns the cached version.
func NewIndexContext(posts []Post) *IndexContext {
	if cachedIndex == nil {
		cachedIndex = &IndexContext{*NewBaseContext(), posts}
	}
	return cachedIndex
}

// NewErrorContext creates a new error context.
func NewErrorContext(err error) *ErrorContext {
	return &ErrorContext{*NewBaseContext(), err.Error()}
}

// LoadTemplates loads the templates from the blog folder.
func LoadTemplates() error {
	cachedPages = make(map[string]*PageContext)
	cachedPosts = make(map[string]*PostContext)
	templates = make(map[string]*template.Template)

	displays, err := filepath.Glob(path.Join(BlogFolder, TemplateFolder, DisplayFolder, "*.html"))
	if err != nil {
		Error.Println("error loading display templates:", err)
		return err
	}
	Trace.Println("displays:", strings.Join(displays, ","))

	includes, err := filepath.Glob(path.Join(BlogFolder, TemplateFolder, IncludeFolder, "*.html"))
	if err != nil {
		Error.Println("error loading include templates:", err)
		return err
	}
	Trace.Println("includes:", strings.Join(includes, ","))

	for _, display := range displays {
		files := append(includes, display)
		name := strings.TrimSuffix(filepath.Base(display), filepath.Ext(display))
		templates[name] = template.Must(template.New(name).ParseFiles(files...))
		Trace.Println("load display template:", name)
	}
	return nil
}

// RenderPage renders a page or throws an error if the template is missing.
func RenderPage(w io.Writer, name string, context interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("template not found")
	}
	return tmpl.ExecuteTemplate(w, "base", context)
}

// ClearCache clears the context cache.
func ClearCache() {
	Info.Println("clearing cache")
	blogContext = nil
	cachedPages = make(map[string]*PageContext)
	cachedPosts = make(map[string]*PostContext)
	cachedIndex = nil
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

// AddLinks add links to navigation.
func AddLinks() {
	for key, value := range GlobalConfig.Links {
		link := NavigationLink{key, value}
		AddNavItem(&link)
	}
}
