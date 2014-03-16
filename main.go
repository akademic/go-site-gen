package main

import (
    "os"
    "path/filepath"
    "log"
    "fmt"
)

type SiteData struct {
    DomainName string
}

var (
    PostsDir string
    PublicDir string
    TemplatesDir string
    IncludesDir string
    BlogDir string
    SiteDir string
    SiteDataVar SiteData
)

const (
    DateFormat     = "2006-01-02"
    DateTimeFormat = DateFormat + " 15:04"
    RecentPostsCount = 4
)

func _init() {

    // Initialize directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}

    PublicDir = filepath.Join(pwd, "public")
    BlogDir = filepath.Join(PublicDir, "blog")
	PostsDir = filepath.Join(pwd, "_posts")
	TemplatesDir = filepath.Join(pwd, "_layouts")
	SiteDir = filepath.Join(pwd, "_site")
    IncludesDir = filepath.Join(pwd, "_includes")

    //create directories
    os.MkdirAll(PostsDir, 0700)
    os.MkdirAll(TemplatesDir, 0700)
    os.MkdirAll(SiteDir, 0700)
    os.MkdirAll(IncludesDir, 0700)


    SiteDataVar.DomainName = ""
}

func main() {
    _init()
    genPosts()
    createIndexPage()
    genSite()
}

func die(format string, v ...interface{}) {
    os.Stderr.WriteString(fmt.Sprintf(format+"\n", v...))
    os.Exit(1)
}
