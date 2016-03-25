package main

import (
	"testing"
	"time"
)

// TestGetURL testes the URL creation method.
func TestGetURL(t *testing.T) {
	testCases := []struct {
		post Post
		url  string
	}{
		{
			Post{"", "", time.Now(), "", "", "my-SLUG"},
			PostBaseURL + "/my-slug",
		},
		{
			Post{"", "", time.Now(), "", "", "example slug"},
			PostBaseURL + "/exampleslug",
		},
	}

	for _, c := range testCases {
		url := c.post.GetURL()
		if url != c.url {
			t.Error(
				"For", c.post.Slug,
				"expected", c.url,
				"got", url,
			)
		}
	}
}

// TestRender testes rendering and sanitizing.
func TestRender(t *testing.T) {
	testCases := []struct {
		post    Post
		content string
	}{
		{
			Post{"", "", time.Now(), "# Hello", "", ""},
			"<h1>Hello</h1>\n",
		},
		{
			Post{"", "", time.Now(), "<script></script>", "", ""},
			"\n",
		},
	}

	for _, c := range testCases {
		c.post.Render()
		if c.post.HTMLContent != c.content {
			t.Error(
				"For", c.post.MDContent,
				"expected", c.content,
				"got", c.post.HTMLContent,
			)
		}
	}
}
