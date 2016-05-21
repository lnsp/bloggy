package main

import (
	"bufio"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// PostsFolder is the default path for posts.
const (
	PostsFolder = "posts"
	PagesFolder = "pages"
)

// ParseData is a structure to store post file header data.
type ParseData struct {
	Title       string `yaml:"title"`
	Subtitle    string `yaml:"subtitle"`
	PublishDate string `yaml:"date"`
	Slug        string `yaml:"slug"`
	content     string
}

func (p *ParseData) SetContent(c string) {
	p.content = c
}

func (p *ParseData) Content() string {
	return p.content
}

type Entry interface {
	GetContent() string
	GetURL() string
}

// Post is a structure to store a posts title, subtitle, publishing date and content.
type Post struct {
	Title       string
	Subtitle    string
	PublishDate time.Time
	Slug        string
	Content     string
}

func (p *Post) GetContent() string {
	return p.Content
}

// GetURL generates the absolute URL from /.
func (p *Post) GetURL() string {
	rgx, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := rgx.ReplaceAllString(p.Slug, "")
	return PostBaseURL + strings.ToLower(slugged)
}

func (p *Post) Age() int64 {
	return time.Now().Unix() - p.PublishDate.Unix()
}

// Page is a structure to store a page title etc.
type Page struct {
	Title   string
	Slug    string
	Content string
}

func (p *Page) GetContent() string {
	return p.Content
}

func (p *Page) GetURL() string {
	regex, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := regex.ReplaceAllString(p.Slug, "")
	return PageBaseURL + strings.ToLower(slugged)
}

// ByAge implements a interface to sort a slice of posts by publishing date.
type ByAge []Post

func (b ByAge) Len() int {
	return len(b)
}
func (b ByAge) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b ByAge) Less(i, j int) bool {
	return b[i].Age() < b[j].Age()
}

// BlogPosts is the slice of published posts.
var Posts []Post
var Pages []Page
var PostBySlug map[string]*Post
var PageBySlug map[string]*Page

// Render generates HTML from a posts markdown content.
func Render(e Entry) string {
	output := blackfriday.MarkdownCommon([]byte(e.GetContent()))
	return string(bluemonday.UGCPolicy().SanitizeBytes([]byte(output)))
}

func parseFile(file string) (data *ParseData, err error) {
	// Open the post file
	input, openError := os.Open(file)
	if openError != nil {
		return nil, openError
	}
	defer input.Close()

	// Scan all input lines
	inputScanner := bufio.NewScanner(input)
	inputScanner.Split(bufio.ScanLines)

	header, body, headerIndex := "", "", 0
	for inputScanner.Scan() {
		line := inputScanner.Text()
		switch headerIndex {
		case 0:
			if line == "---" {
				headerIndex = 1
			}
		case 1:
			if line == "---" {
				headerIndex = 2
			} else {
				header += line + "\n"
			}
		case 2:
			body += line + "\n"
		}
	}

	data = new(ParseData)
	data.SetContent(body)

	// Decode JSON header
	yamlError := yaml.Unmarshal([]byte(header), &data)
	if yamlError != nil {
		Error.Println("Error while parsing", file, ":", yamlError)
	}

	// Generate slug from file name if needed
	if len(data.Slug) == 0 {
		data.Slug = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	}
	data.Slug = strings.ToLower(data.Slug)

	Trace.Printf("parseFile: %!v\n", *data)
	return data, nil
}

func loadDirectory(dir string, callback func(*ParseData) error) error {
	glob := path.Join(dir, "*.md")
	dirEntries, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	Trace.Println("loadDirectory: searching in", glob)
	Trace.Println("loadDirectory: found entries:", strings.Join(dirEntries, ","))
	for _, entry := range dirEntries {
		// Parse file entry
		data, err := parseFile(entry)
		if err != nil {
			Warning.Println("Parse error:", err)
			continue
		}

		err = callback(data)
		if err != nil {
			Error.Println(err)
			continue
		}
		Trace.Println("Read file:", entry)
	}

	return nil
}

func addPost(data *ParseData) error {
	p := Post{
		Title:    data.Title,
		Subtitle: data.Subtitle,
		Slug:     data.Slug,
		Content:  data.Content(),
	}
	date, err := time.Parse("2006-Jan-02", data.PublishDate)
	if err != nil {
		return err
	}
	p.PublishDate = date

	Posts = append(Posts, p)
	PostBySlug[p.Slug] = &p
	return nil
}

func addPage(data *ParseData) error {
	p := Page{
		Title:   data.Title,
		Slug:    data.Slug,
		Content: data.Content(),
	}

	Pages = append(Pages, p)
	PageBySlug[p.Slug] = &p
	return nil
}

// LoadPosts loads all published posts from the blog folder.
func LoadPosts() error {
	Posts = make([]Post, 0)
	PostBySlug = make(map[string]*Post)
	err := loadDirectory(path.Join(BlogFolder, PostsFolder), addPost)
	if err != nil {
		return err
	}
	// Sort all posts by age
	sort.Sort(ByAge(Posts))
	return nil
}

// LoadPages loads all pages from the pages/ folder.
func LoadPages() error {
	Pages = make([]Page, 0)
	PageBySlug = make(map[string]*Page)
	err := loadDirectory(path.Join(BlogFolder, PagesFolder), addPage)
	return err
}

// LatestPosts returns a slice of the latest posts.
func LatestPosts(count int) []Post {
	size := len(Posts)
	if size < count {
		return Posts
	}
	return Posts[:count]
}
