package templates

import (
	"html/template"
	"time"
	"fmt"
	"github.com/mooxmirror/blog/posts"
	"github.com/mooxmirror/blog/config"
)

const (
	templatesFolder = "templates"
	baseName = "base"
	postName = "post"
	indexName = "index"
	fileExtension = ".html"
)

type BaseContext struct {
	BlogTitle, BlogSubtitle, BlogAuthor, BlogYear, BlogEmail, BlogUrl string
}

type PostContext struct {
	BaseContext
	PostTitle, PostSubtitle, PostDate string
	PostContent template.HTML
}

func GetBaseContext(cfg config.Config) BaseContext {
	return BaseContext{cfg.BlogTitle, cfg.BlogSubtitle, cfg.BlogAuthor, fmt.Sprint(time.Now().Year()), cfg.BlogEmail, cfg.BlogUrl}
}

func GetPostContext(cfg config.Config, post posts.Post) PostContext {
	return PostContext{GetBaseContext(cfg), post.Title, post.Subtitle, post.PublishDate.String(), template.HTML(post.HTMLContent)};
}

var (
	BaseTemplate, PostTemplate, IndexTemplate *template.Template
)

func Load(folder string) {
	folder += "/" + templatesFolder

	BaseTemplate = template.Must(template.ParseFiles(folder + "/" + baseName + fileExtension))
	PostTemplate = template.Must(template.ParseFiles(folder + "/" + postName + fileExtension))
	//IndexTemplate = template.Must(template.New(indexName).ParseFiles(folder + "/" + indexName + fileExtension))
}
