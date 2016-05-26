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
	blogContext *BaseContext
	cachedPages map[string]*PageContext
	cachedPosts map[string]*PostContext
	cachedIndex *IndexContext
)

// BaseContext stores basic context information like title, author etc.
type BaseContext struct {
	BlogTitle, BlogSubtitle, BlogAuthor, BlogYear, BlogEmail, BlogURL string
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
	PostURL     string
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

// NewBaseContext either creates a new BaseContext from the global blog configuration or returns the cached version.
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
			return nil, errors.New("page not found")
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
}

// RenderPage renders a page or throws an error if the template is missing.
func RenderPage(w io.Writer, name string, context interface{}) error {
	tmpl, ok := templates[name]
	if !ok {
		return errors.New("template not found")
	}
	return tmpl.ExecuteTemplate(w, "base", context)
}
