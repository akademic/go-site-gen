package main

import (
    "os"
    "path/filepath"
    "log"
    "fmt"
    "github.com/pelletier/go-toml"
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

    //from config file
    RecentPostsCount int64
    PublicDirName string
    BlogDirName string
    RssFileName string
)

const (
    DateFormat     = "2006-01-02"
    DateTimeFormat = DateFormat + " 15:04"
    ConfigFile = "config.toml"
)

func _init() {

    // Initialize directories
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("FATAL ", err)
	}
    
    loadConf()

    PublicDir = filepath.Join(pwd, PublicDirName)
    BlogDir = filepath.Join(PublicDir, BlogDirName)
	PostsDir = filepath.Join(pwd, "_posts")
	TemplatesDir = filepath.Join(pwd, "_layouts")
	SiteDir = filepath.Join(pwd, "_site")
    IncludesDir = filepath.Join(pwd, "_includes")

    //create directories
    os.MkdirAll(PostsDir, 0700)
    os.MkdirAll(TemplatesDir, 0700)
    os.MkdirAll(SiteDir, 0700)
    os.MkdirAll(IncludesDir, 0700)
}

func main() {
    _init()
    genPosts()
    createIndexPage()
    createPostRss()
    genSite()
}

func loadConf() {
    cfg, err := toml.LoadFile(ConfigFile)
    if err != nil {
        panic(err)
    }

    SiteDataVar.DomainName = cfg.Get("domain").(string)
    RecentPostsCount = cfg.Get("recent_posts_count").(int64)
    PublicDirName = cfg.Get("path.public_dir_name").(string)
    BlogDirName = cfg.Get("path.blog_dir_name").(string)
    RssFileName = cfg.Get("path.rss_file_name").(string)
}

func die(format string, v ...interface{}) {
    os.Stderr.WriteString(fmt.Sprintf(format+"\n", v...))
    os.Exit(1)
}
