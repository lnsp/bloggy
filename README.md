# go-blog
A simple blog.

## Set up go-blog
You need to have a working Go environment installed (v1.5+).

```bash
$ go get github.com/mooxmirror/go-blog # download it from GitHub
$ go install github.com/mooxmirror/go-blog # install dependencies, build it
$ $GOPATH/bin/go-blog -init="my-blog-name" # resets your blog folder, starts the server
$ $GOPATH/bin/go-blog "my-blog-name" # starts the server without reset
```
