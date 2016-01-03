package templates

import (
	"html/template"
	"github.com/mooxmirror/blog/posts"
)

const (
	baseName = "base"
	postName = "post"
	indexName = "index"
	indexEntryName = "entry"
	fileExtension = ".html"
)

var (
	BaseTemplate, PostTemplate, IndexTemplate, IndexEntryTemplate *template.Template
)

func Load(folder string) {
	BaseTemplate = template.Must(template.New(baseName).ParseFiles(folder + "/" + baseName + fileExtension))
	PostTemplate = template.Must(template.New(postName).ParseFiles(folder + "/" + postName + fileExtension))
	IndexTemplate = template.Must(template.New(indexName).parseFiles(folder + "/" + indexName + fileExtension))
	IndexEntryTemplate = template.Must(template.New(indexEntryName).parseFiles(folder + "/" + indexEntryName + fileExtension))
}
