# bloggy
**Remember: This is heavily Work-in-Progress and far from being a stable software.**

This is a blog system I am currently developing. It is focused on simplicity and conventions over configuration. I recommend to proxy the service using NGINX or comparable software.

## Set up blog
You need to have a working Go environment installed (v1.5+).

```bash
$ go get github.com/mooxmirror/go-blog # download it from GitHub
$ go install github.com/mooxmirror/go-blog # install dependencies, build it
$ $GOPATH/bin/blog -reset -blog="example-blog" # resets your blog folder, starts the server
$ $GOPATH/bin/blog -blog="example-blog" # starts the server without reset
```

## Folder structure
```
/
  /blog (runnable)
  /example-blog (folder)
    /config.json
    /posts
      2015-12-31.md
      2015-12-13.md
      2015-11-09.md
      2015-10-31.md
    /templates
      base.html
      index.html
      post.html
      header.html
      footer.html
    /files
      header-image.png
      favicon.ico
      2015-12-31.png
```

The **config.json** file stores basic configuration options like the blog's name, host address etc.
The blog posts are stored in the **posts** folder. Every post file has to begin with a date representing the publishing date of the post. Every post file has to contain a header marked by `---`. This header has to be in JSON.