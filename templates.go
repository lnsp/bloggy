package main

import (
	"fmt"
	"html/template"
	"path"
	"time"

	humanize "github.com/dustin/go-humanize"
)

// TemplatesFolder is the default folder for templates.
const TemplatesFolder = "templates"

// BaseContext stores basic context information like title, author etc.
type BaseContext struct {
	BlogTitle, BlogSubtitle, BlogAuthor, BlogYear, BlogEmail, BlogURL string
}

// PostContext stores additional information like post title, date etc.
type PostContext struct {
	BaseContext
	PostTitle, PostSubtitle, PostDate string
	PostContent                       template.HTML
}

// IndexContext stores index information like a list of posts.
type IndexContext struct {
	BaseContext
	ListOfContents []Post
}

// GetBaseContext creates a BaseContext from the global blog configuration.
func GetBaseContext() BaseContext {
	return BaseContext{
		GlobalConfig.BlogTitle,
		GlobalConfig.BlogSubtitle,
		GlobalConfig.BlogAuthor,
		fmt.Sprint(time.Now().Year()),
		GlobalConfig.BlogEmail,
		GlobalConfig.BlogURL,
	}
}

// GetPostContext creates a PostContext from the global blog configuration and a specified post.
func GetPostContext(post Post) PostContext {
	return PostContext{
		GetBaseContext(),
		post.Title,
		post.Subtitle,
		humanize.Time(post.PublishDate),
		template.HTML(post.HTMLContent),
	}
}

// GetIndexContext creates a IndexContext from the global blog configuration and a list of posts.
func GetIndexContext(posts []Post) IndexContext {
	return IndexContext{GetBaseContext(), posts}
}

var (
	// BaseTemplate is the base template for all pages.
	BaseTemplate *template.Template
	// PostTemplate is the template for post pages.
	PostTemplate *template.Template
	// IndexTemplate is the template for index pages.
	IndexTemplate *template.Template
)

// LoadTemplates loads the templates from the blog folder.
func LoadTemplates() {
	baseName := path.Join(BlogFolder, TemplatesFolder, "base.html")
	postName := path.Join(BlogFolder, TemplatesFolder, "post.html")
	indexName := path.Join(BlogFolder, TemplatesFolder, "index.html")
	entryName := path.Join(BlogFolder, TemplatesFolder, "entry.html")

	PostTemplate = template.Must(template.New("post").ParseFiles(baseName, postName))
	IndexTemplate = template.Must(template.New("index").ParseFiles(baseName, indexName, entryName))
}
