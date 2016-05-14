package main

import (
	"bufio"
	"errors"
	"io/ioutil"
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

// FileHeader is a structure to store post file header data.
type FileHeader struct {
	Title       string `yaml:"title"`
	Subtitle    string `yaml:"subtitle"`
	PublishDate string `yaml:"date"`
	Slug        string `yaml:"slug"`
}

// Post is a structure to store a posts title, subtitle, publishing date and content.
type Post struct {
	Title       string
	Subtitle    string
	PublishDate time.Time
	MDContent   string
	HTMLContent string
	Slug        string
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
	return b[i].PublishDate.Unix() > b[j].PublishDate.Unix()
}

// BlogPosts is the slice of published posts.
var BlogPosts, BlogPages []Post

// Render generates HTML from a posts markdown content.
func (p *Post) Render() {
	output := blackfriday.MarkdownCommon([]byte(p.MDContent))
	p.HTMLContent = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(output)))
}

// GetURL generates the absolute URL from /.
func (p Post) GetURL() string {
	rgx, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := rgx.ReplaceAllString(p.Slug, "")
	return PostBaseURL + "/" + strings.ToLower(slugged)
}

// NewPost creates a new post with a specified title, subtitle, publishing date and content.
func NewPost(title, subtitle string, date time.Time, content, slug string) Post {
	p := Post{title, subtitle, date, content, "", slug}
	p.Render()
	return p
}

func parsePostFile(folder, file string) (Post, error) {
	// Open the post file
	input, openError := os.Open(path.Join(BlogFolder, folder, file))
	if openError != nil {
		return Post{}, openError
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

	headerData := FileHeader{}

	// Decode JSON header
	yamlError := yaml.Unmarshal([]byte(header), &headerData)
	if yamlError != nil {
		return Post{}, yamlError
	}

	// parse Date from header data
	date, parseError := time.Parse("2006-Jan-02", headerData.PublishDate)
	if parseError != nil {
		return Post{}, parseError
	}

	if len(headerData.Slug) == 0 {
		headerData.Slug = strings.TrimSuffix(file, filepath.Ext(file))
	}

	return NewPost(headerData.Title, headerData.Subtitle, date, body, headerData.Slug), nil
}

// LoadPosts loads all published posts from the blog folder.
func LoadPosts() {
	BlogPosts = make([]Post, 0)
	dirEntries, readDirError := ioutil.ReadDir(path.Join(BlogFolder, PostsFolder))
	if readDirError != nil {
		Warning.Println("Failed to open post folder:", readDirError)
		return
	}
	for _, entry := range dirEntries {
		// Ignore directories
		if entry.IsDir() {
			continue
		}

		post, parseError := parsePostFile(PostsFolder, entry.Name())
		if parseError != nil {
			Warning.Println("Failed to parse file:", parseError)
			continue
		}
		BlogPosts = append(BlogPosts, post)
		Trace.Println("Read post file", entry.Name())
	}

	// Sort all posts by age
	sort.Sort(ByAge(BlogPosts))
	Trace.Println("Serving", len(BlogPosts), "blog posts")
}

// LoadPages loads all pages from the pages/ folder.
func LoadPages() {
	BlogPages = make([]Post, 0)
	dirEntries, readDirError := ioutil.ReadDir(path.Join(BlogFolder, PagesFolder))
	if readDirError != nil {
		Warning.Println("Failed to open page folder:", readDirError)
		return
	}
	for _, entry := range dirEntries {
		// Ignore directories
		if entry.IsDir() {
			continue
		}

		page, parseError := parsePostFile(PagesFolder, entry.Name())
		if parseError != nil {
			Warning.Println("Failed to parse file:", parseError)
			continue
		}
		BlogPages = append(BlogPages, page)
		Trace.Println("Read page file", entry.Name())
	}

	// Sort all posts by age
	Trace.Println("Serving", len(BlogPages), "blog pages")
}

// GetLatestsPosts returns a slice of the latest posts.
func GetLatestsPosts(count int) []Post {
	size := len(BlogPosts)
	if size < count {
		return BlogPosts
	}
	return BlogPosts[:count]
}

// FindPost returns a post with the specified slug.
func FindPost(slug string) (*Post, error) {
	for _, p := range BlogPosts {
		if p.Slug == slug {
			return &p, nil
		}
	}
	return nil, errors.New("Post not found.")
}
