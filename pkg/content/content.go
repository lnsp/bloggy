package content

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/go-yaml/yaml"
	"github.com/lnsp/bloggy/pkg/config"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/sirupsen/logrus"
)

const (
	// PostsFolder is the default folder for blog posts.
	PostsFolder = "posts"
	// PagesFolder is the default folder for blog pages.
	PagesFolder = "pages"
	// FileDateFormat is the date format required in a post's header.
	FileDateFormat = "2006-Jan-02"
)

type URLResolver interface {
	Page(string) string
	Post(string) string
}

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
	GetTitle() string
}

// Post stores a title, a page slug and the body content.
type Post struct {
	Title       string
	Subtitle    string
	PublishDate time.Time
	Slug        string
	Content     string
	Resolver    URLResolver
}

// GetContent returns the content body of the post.
func (p *Post) GetContent() string {
	return p.Content
}

// GetTitle returns the title of the post.
func (p *Post) GetTitle() string {
	return p.Title
}

// GetURL generates a URL from the post route url and the post slug.
func (p *Post) GetURL() string {
	rgx, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := rgx.ReplaceAllString(p.Slug, "")
	return p.Resolver.Post(strings.ToLower(slugged))
}

// Age returns the age of the post in seconds.
func (p *Post) Age() int64 {
	return time.Now().Unix() - p.PublishDate.Unix()
}

// Page stores a title, a page slug and the body content.
type Page struct {
	Title    string
	Slug     string
	Content  string
	Resolver URLResolver
}

// GetContent returns the content body of the page.
func (p *Page) GetContent() string {
	return p.Content
}

// GetTitle returns the title of the page.
func (p *Page) GetTitle() string {
	return p.Title
}

// GetURL generates a URL from the page route url and the page slug.
func (p *Page) GetURL() string {
	regex, _ := regexp.Compile("[^A-Za-z\\-]")
	slugged := regex.ReplaceAllString(p.Slug, "")
	return p.Resolver.Page(strings.ToLower(slugged))
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

type Index struct {
	// Posts stores all blog posts.
	Posts []Post

	// Pages stores all blog pages.
	Pages []Page

	// PostBySlug matches each slug to its post.
	PostBySlug map[string]*Post

	// PageBySlug matches each slug to its page.
	PageBySlug map[string]*Page

	Resolver URLResolver
}

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
	if err := yaml.Unmarshal([]byte(header), &data); err != nil {
		return nil, err
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

	logrus.WithFields(logrus.Fields{
		"glob":    glob,
		"entries": dirEntries,
	}).Debug("scanning directory")
	for _, entry := range dirEntries {
		// Parse file entry
		data, err := parseFile(entry)
		if err != nil {
			logrus.WithField("file", entry).Warn("failed to parse")
			continue
		}

		err = callback(data)
		if err != nil {
			logrus.WithField("file", entry).Warn("failed to callback")
			continue
		}
	}

	return nil
}

// AddPost creates a new post from the parsed data.
func (c *Index) AddPost(data *ParseData) error {
	p := Post{
		Title:    data.Title,
		Subtitle: data.Subtitle,
		Slug:     data.Slug,
		Content:  data.Content(),
		Resolver: c.Resolver,
	}
	date, err := time.Parse(FileDateFormat, data.PublishDate)
	if err != nil {
		return err
	}
	p.PublishDate = date

	c.Posts = append(c.Posts, p)
	c.PostBySlug[p.Slug] = &p
	return nil
}

// AddPage creates a new page from the parsed data.
func (c *Index) AddPage(data *ParseData) error {
	p := Page{
		Title:    data.Title,
		Slug:     data.Slug,
		Content:  data.Content(),
		Resolver: c.Resolver,
	}

	c.Pages = append(c.Pages, p)
	c.PageBySlug[p.Slug] = &p
	return nil
}

func NewIndex(cfg *config.Config, resolver URLResolver) (*Index, error) {
	index := &Index{
		PostBySlug: make(map[string]*Post),
		PageBySlug: make(map[string]*Page),
		Resolver:   resolver,
	}
	err := loadDirectory(path.Join(cfg.Base, PostsFolder), index.AddPost)
	if err != nil {
		return nil, fmt.Errorf("load posts dir: %w", err)
	}
	// Sort all posts by age
	sort.Sort(ByAge(index.Posts))
	if err := loadDirectory(path.Join(cfg.Base, PagesFolder), index.AddPage); err != nil {
		return nil, fmt.Errorf("load pages dir: %w", err)
	}
	return index, nil
}

// LatestPosts returns a slice of the latest blog posts.
func (c *Index) LatestPosts(count int) []Post {
	size := len(c.Posts)
	if size < 0 || size < count {
		return c.Posts
	}
	return c.Posts[:count]
}
