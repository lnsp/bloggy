package posts

import (
	"time"
	"bufio"
	"os"
	"io/ioutil"
	"log"
	"sort"
	"encoding/json"
	"github.com/russross/blackfriday"
	"github.com/microcosm-cc/bluemonday"
)
type Post struct {
	Title, Subtitle string
	PublishDate time.Time
	MDContent string
	HTMLContent string
}
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

var PublicPosts []Post

func (p *Post) Render() {
	output := blackfriday.MarkdownCommon([]byte(p.MDContent))
	p.HTMLContent = string(bluemonday.UGCPolicy().SanitizeBytes([]byte(output)))
}

func NewPost(title, subtitle string, date time.Time, content string) Post {
	p := Post{title, subtitle, date, content, ""}
	p.Render()
	return p
}

func Load(folder string) {
	PublicPosts = make([]Post, 0)
	dirEntries, readError := ioutil.ReadDir(folder + "/posts")
	if readError != nil {
		log.Fatal(readError)
		os.Exit(1)
	}
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		input, err := os.Open(folder + "/posts/" + entry.Name())
		if err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		defer input.Close()

		inputScanner := bufio.NewScanner(input)
		inputScanner.Split(bufio.ScanLines)

		header, body, headerActive := "{\n", "", true
		for inputScanner.Scan() {
			line := inputScanner.Text()
			log.Println("Read line", line)
			if headerActive {
				if line == "---" {
					headerActive = false
				} else {
					header += line + "\n"
				}
			} else {
				body += line + "\n";
			}
		}
		header += "}"
		log.Println(header)
		log.Println(body)

		var headerData struct {
			Title string `json:"title"`
			Subtitle string `json:"subtitle"`
			PublishDate string `json:"date"`
		}

		// load json header data
		jsonError := json.Unmarshal([]byte(header), &headerData)
		if jsonError != nil {
			log.Fatal(jsonError)
			os.Exit(1)
		}

		// parse Date from header data
		date, parseError := time.Parse("2015-12-31", headerData.PublishDate)
		if parseError != nil {
			log.Fatal(parseError)
			os.Exit(1)
		}
		PublicPosts = append(PublicPosts, NewPost(headerData.Title, headerData.Subtitle, date, body))
		log.Println("Read post file", entry.Name())
	}
	sort.Sort(ByAge(PublicPosts))
}

func GetNewest(count int) []Post {
	return PublicPosts[:count]
}
