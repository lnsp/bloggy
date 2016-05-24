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

const (
	// PostsFolder is the default folder for blog posts.
	PostsFolder = "posts"
	// PagesFolder is the default folder for blog pages.
	PagesFolder = "pages"
)

// ParseData stores the parsed data of a file.
type ParseData struct {
	Title       string `yaml:"title"`
	Subtitle    string `yaml:"subtitle"`
	PublishDate string `yaml:"date"`
	Slug        string `yaml:"slug"`
	content     string
}

// SetContent sets the parsed content.
func (p *ParseData) SetContent(c string) {
	p.content = c
}

// Content returns the parsed content.
func (p *ParseData) Content() string {
	return p.content
}

// Entry has content and be located by an URL.
type Entry interface {
	GetContent() string
	GetURL() string
}

// Post stores a title, a page slug and the body content.
type Post struct {
	Title       string
	Subtitle    string
	PublishDate time.Time
	Slug        string
	Content     string
}

// GetContent returns the content body of the post.
func (p *Post) GetContent() string {
	return p.Content
}

// GetURL generates a URL from the post route url and the post slug.
func (p *Post) GetURL() string {
	rgx, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := rgx.ReplaceAllString(p.Slug, "")
	return PostBaseURL + strings.ToLower(slugged)
}

// Age returns the age of the post in seconds.
func (p *Post) Age() int64 {
	return time.Now().Unix() - p.PublishDate.Unix()
}

// Page stores a title, a page slug and the body content.
type Page struct {
	Title   string
	Slug    string
	Content string
}

// GetContent returns the content body of the page.
func (p *Page) GetContent() string {
	return p.Content
}

// GetURL generates a URL from the page route url and the page slug.
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

// Posts stores all blog posts.
var Posts []Post

// Pages stores all blog pages.
var Pages []Page

// PostBySlug matches each slug to its post.
var PostBySlug map[string]*Post

// PageBySlug matches each slug to its page.
var PageBySlug map[string]*Page

// Render generates HTML from an entries markdown content.
func Render(e Entry) string {
	output := blackfriday.MarkdownCommon([]byte(e.GetContent()))
	return string(bluemonday.UGCPolicy().SanitizeBytes([]byte(output)))
}

// parseFile parses a file and returns a pointer to the parsed data or an error.
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
		Error.Println(file, ":", yamlError)
	}

	// Generate slug from file name if needed
	if len(data.Slug) == 0 {
		data.Slug = strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	}
	data.Slug = strings.ToLower(data.Slug)
	return data, nil
}

// loadDirectory searches a directory for markdown files, parses them and calls a function for each of them.
func loadDirectory(dir string, callback func(*ParseData) error) error {
	glob := path.Join(dir, "*.md")
	dirEntries, err := filepath.Glob(glob)
	if err != nil {
		return err
	}

	Trace.Println("searching in", glob)
	Trace.Println("found entries:", strings.Join(dirEntries, ","))
	for _, entry := range dirEntries {
		// Parse file entry
		data, err := parseFile(entry)
		if err != nil {
			Warning.Println("parse error:", err)
			continue
		}

		err = callback(data)
		if err != nil {
			Error.Println(err)
			continue
		}
		Trace.Println("read file:", entry)
	}

	return nil
}

// addPost creates a new post from the parsed data.
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

// addPage creates a new page from the parsed data.
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

// LoadPosts loads all posts from the posts/ folder.
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

// LatestPosts returns a slice of the latest blog posts.
func LatestPosts(count int) []Post {
	size := len(Posts)
	if size < count {
		return Posts
	}
	return Posts[:count]
}
