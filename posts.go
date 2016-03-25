package main

import (
	"bufio"
	"io/ioutil"
	"os"
	"sort"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
)

// FileHeader is a structure to store post file header data.
type FileHeader struct {
	Title       string `yaml:"title"`
	Subtitle    string `yaml:"subtitle"`
	PublishDate string `yaml:"date"`
}

// Post is a structure to store a posts title, subtitle, publishing date and content.
type Post struct {
	Title, Subtitle string
	PublishDate     time.Time
	MDContent       string
	HTMLContent     string
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
	return b[i].PublishDate.Unix() < b[j].PublishDate.Unix()
}

// BlogPosts is the slice of published posts.
var BlogPosts []Post

// Render generates HTML from a posts markdown content.
func (p *Post) Render() {
	output := blackfriday.MarkdownCommon([]byte(p.MDContent))
	p.HTMLContent = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(output)))
	Trace.Println("Rendering post", p.Title)
}

// NewPost creates a new post with a specified title, subtitle, publishing date and content.
func NewPost(title, subtitle string, date time.Time, content string) Post {
	p := Post{title, subtitle, date, content, ""}
	p.Render()
	return p
}

// LoadPosts loads all published posts from the blog folder.
func LoadPosts(folder string) {
	BlogPosts = make([]Post, 0)
	dirEntries, readDirError := ioutil.ReadDir(folder + "/posts")
	if readDirError != nil {
		Warning.Println("Failed to open post folder:", readDirError)
		return
	}
	for _, entry := range dirEntries {
		// Ignore directories
		if entry.IsDir() {
			continue
		}

		// Open the post file
		input, openError := os.Open(folder + "/posts/" + entry.Name())
		if openError != nil {
			Warning.Println("Failed to read post")
			continue
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
			Warning.Println("Failed parsing YAML:", yamlError)
			continue
		}

		// parse Date from header data
		date, parseError := time.Parse("2006-Jan-02", headerData.PublishDate)
		if parseError != nil {
			Warning.Println("Failed to parse date:", parseError)
			continue
		}
		BlogPosts = append(BlogPosts, NewPost(headerData.Title, headerData.Subtitle, date, body))
		Info.Println("Read post file", entry.Name())
	}

	// Sort all posts by age
	Info.Println("Sorting posts by age")
	sort.Sort(ByAge(BlogPosts))
	Info.Println("Serving", len(BlogPosts), "blog posts")
}

// GetLatestsPosts returns a slice of the latest posts.
func GetLatestsPosts(count int) []Post {
	size := len(BlogPosts)
	if size < count {
		return BlogPosts
	}
	return BlogPosts[:count]
}
