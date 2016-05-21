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
	displayFolder  = "displays"
	includeFolder  = "includes"
)

var (
	// BlogTemplates stores all required blog templates.
	templates   map[string]*template.Template
	blogContext *BaseContext
	cachedPages map[string]*PageContext
	cachedPosts map[string]*PostContext
	cachedIndex *IndexContext
)

// BaseContext stores basic context information like title, author etc.
type BaseContext struct {
	BlogTitle, BlogSubtitle, BlogAuthor, BlogYear, BlogEmail, BlogURL string
}

// PostContext stores additional information like post title, date etc.
type PostContext struct {
	BaseContext
	PostTitle    string
	PostSubtitle string
	PostDate     string
	PostContent  template.HTML
	PostURL      string
}

type PageContext struct {
	BaseContext
	PageTitle   string
	PageContent template.HTML
	PostURL     string
}

// IndexContext stores index information like a list of posts.
type IndexContext struct {
	BaseContext
	LatestPosts []Post
}

// ErrorContext stores error information.
type ErrorContext struct {
	BaseContext
	Message string
}

// NewBaseContext creates a BaseContext from the global blog configuration.
func NewBaseContext() *BaseContext {
	if blogContext == nil {
		blogContext = &BaseContext{
			GlobalConfig.BlogTitle,
			GlobalConfig.BlogSubtitle,
			GlobalConfig.BlogAuthor,
			fmt.Sprint(time.Now().Year()),
			GlobalConfig.BlogEmail,
			GlobalConfig.BlogURL,
		}
	}
	return blogContext
}

// GetPostContext creates a PostContext from the global blog configuration and a specified post.
func NewPostContext(slug string) (*PostContext, error) {
	context, ok := cachedPosts[slug]
	if !ok {
		post, ok := PostBySlug[slug]
		if !ok {
			return nil, errors.New("Post not found")
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
		Trace.Println("Create cache version of post", slug)
	}
	return context, nil
}

func NewPageContext(slug string) (*PageContext, error) {
	context, ok := cachedPages[slug]
	if !ok {
		page, ok := PageBySlug[slug]
		if !ok {
			return nil, errors.New("Page not found")
		}
		context = &PageContext{
			*NewBaseContext(),
			page.Title,
			template.HTML(Render(page)),
			page.GetURL(),
		}
		cachedPages[slug] = context
		Trace.Println("Create cache version of page", slug)
	}
	return context, nil
}

// GetIndexContext creates a IndexContext from the global blog configuration and a list of posts.
func NewIndexContext(posts []Post) *IndexContext {
	if cachedIndex == nil {
		cachedIndex = &IndexContext{*NewBaseContext(), posts}
	}
	return cachedIndex
}

// GetErrorContext creates a ErrorContext from the error.
func NewErrorContext(err error) *ErrorContext {
	return &ErrorContext{*NewBaseContext(), err.Error()}
}

// LoadTemplates loads the templates from the blog folder.
func LoadTemplates() {
	cachedPages = make(map[string]*PageContext)
	cachedPosts = make(map[string]*PostContext)
	templates = make(map[string]*template.Template)

	displays, err := filepath.Glob(path.Join(BlogFolder, TemplateFolder, displayFolder, "*.html"))
	if err != nil {
		Error.Println("Error loading display templates:", err)
		return
	}
	Trace.Println("Found displays:", strings.Join(displays, ","))

	includes, err := filepath.Glob(path.Join(BlogFolder, TemplateFolder, includeFolder, "*.html"))
	if err != nil {
		Error.Println("Error loading include templates:", err)
		return
	}
	Trace.Println("Found includes:", strings.Join(includes, ","))

	for _, display := range displays {
		files := append(includes, display)
		name := strings.TrimSuffix(filepath.Base(display), filepath.Ext(display))
		templates[name] = template.Must(template.New(name).ParseFiles(files...))
		Trace.Println("Load display template:", name)
	}
}

func RenderPage(w io.Writer, name string, context interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("Template not found")
	}
	return tmpl.ExecuteTemplate(w, "base", context)
}
