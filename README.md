# Bloggy
**Remember: This is heavily Work-in-Progress and far from being a stable software.**

This is a blog system I am currently developing. It is focused on simplicity and conventions over configuration. I recommend to proxy the service using NGINX or comparable software.

## Screenshot
![Screenshot](https://cloud.githubusercontent.com/assets/3391295/22408821/0a8ee4e2-e680-11e6-9c8b-6d123fb21e97.png)
## Getting started
You need to have Git and a working Go environment installed (v1.5+).

```bash
$ # download and install it from GitHub
$ go get github.com/lnsp/bloggy
$ # reset your blog folder, starts the server
$ $GOPATH/bin/bloggy -reset -blog="example-blog"
$ # start the server without resetting the blog
$ $GOPATH/bin/bloggy -blog="example-blog"
```

## Folder structure
```
my-blog
├── config.json
├── pages
│   └── about.md
├── posts
│   ├── first-post.md
│   ├── second-post.md
├── static
│   └── screenshot.png
└── templates
    ├── displays
    │   ├── error.html
    │   ├── index.html
    │   ├── page.html
    │   └── post.html
    └── includes
        └── base.html
```

The **config.json** file stores basic configuration options like the blog's name, host address etc.
The blog posts are stored in the **posts** folder. Every post file has to begin with a date representing the publishing date of the post. Every post file has to contain a header marked by `---`. This header has to be in YAML.

## Post example
```markdown
---
title: Happy New Year!
subtitle: Ideas for 2016
date: 2015-Dec-31
slug: open-source-land
---
## Hello World from Open Source Land!

It is wonderful in here!
```

## Page example
```markdown
---
title: About me
slug: about
---

This is an about-me page.
```
