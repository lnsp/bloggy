package templates

import (
	"html/template"
)

const (
	templatesFolder = "templates"
	baseName = "base"
	postName = "post"
	indexName = "index"
	indexEntryName = "entry"
	fileExtension = ".html"
)

var (
	BaseTemplate, PostTemplate, IndexTemplate, IndexEntryTemplate *template.Template
)

// load all required templates
func Load(folder string) {
	folder += "/" + templatesFolder

	BaseTemplate = template.Must(template.ParseFiles(folder + "/" + baseName + fileExtension))
	PostTemplate = template.Must(template.ParseFiles(folder + "/" + postName + fileExtension))
	//IndexTemplate = template.Must(template.New(indexName).ParseFiles(folder + "/" + indexName + fileExtension))
	//IndexEntryTemplate = template.Must(template.New(indexEntryName).ParseFiles(folder + "/" + indexEntryName + fileExtension))
}
